package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/flant/kubedog/pkg/kube"
	"github.com/flant/kubedog/pkg/tracker"
	"github.com/flant/kubedog/pkg/trackers/rollout"
	"github.com/steffenmllr/validate-kubernetes-deployment/slack"
	"github.com/tidwall/gjson"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type config struct {
	Name            string        `env:"NAME"`
	Namespace       string        `env:"NAMESPACE,required"`
	Deployments     []string      `env:"DEPLOYMENTS,required" envSeparator:","`
	Timout          time.Duration `env:"TIMEOUT" envDefault:"60s"`
	SlackUrl        string        `env:"SLACK_HOOK_URL"`
	GithubEventPath string        `env:"GITHUB_EVENT_PATH"`
	GithubEventName string        `env:"GITHUB_EVENT_NAME"`
}

func main() {
	var err error
	var deploymentSuccess = true
	var attachments []*slack.Attachment

	cfg := config{}
	err = env.Parse(&cfg)

	// Parse Config
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	// Init kube
	err = kube.Init(kube.InitOptions{})
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	// Slack Message
	slackMsg := &slack.Message{
		Text: fmt.Sprintf("*Deployment Status* %s", cfg.Name),
	}

	// Check if we run with a github action
	if cfg.GithubEventName != "" && cfg.GithubEventPath != "" {
		if cfg.GithubEventName != "push" {
			fmt.Printf("Sorry only 'push' events are supported currently")
			os.Exit(1)
		}
		eventbuffer, err := ioutil.ReadFile(cfg.GithubEventPath)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		githubEvent := string(eventbuffer)

		username := gjson.Get(githubEvent, "head_commit.committer.username")
		authorIcon := fmt.Sprintf("https://github.com/%s.png", username)
		authorName := gjson.Get(githubEvent, "head_commit.committer.name")
		commitUrl := gjson.Get(githubEvent, "head_commit.url").String()
		commitMessage := gjson.Get(githubEvent, "head_commit.message").String()
		repositoryName := gjson.Get(githubEvent, "repository.name")
		repositoryIcon := gjson.Get(githubEvent, "repository.owner.avatar_url")
		repositoryUrl := gjson.Get(githubEvent, "repository.html_url")

		tagOrBranch := gjson.Get(githubEvent, "ref").String()
		tagOrBranch = strings.Replace(tagOrBranch, "refs/tags/", "", -1)
		tagOrBranch = strings.Replace(tagOrBranch, "refs/heads/", "", -1)

		if cfg.Name == "" {
			slackMsg.Text = fmt.Sprintf("*Deployment Status*")
		}

		// Let's add some metadata to make
		// this a little bit more usefull
		var fields []slack.AttachmentField
		if username.Exists() {
			fields = append(fields, slack.AttachmentField{
				Title: "User",
				Value: username.String(),
				Short: false,
			})
		}

		fields = append(fields, slack.AttachmentField{
			Title: "Branch/Tag",
			Value: tagOrBranch,
			Short: false,
		})

		fields = append(fields, slack.AttachmentField{
			Title: "Commit",
			Value: commitUrl,
			Short: false,
		})

		attachments = append(attachments, &slack.Attachment{
			Text:       fmt.Sprintf("%s\n", commitMessage),
			AuthorIcon: authorIcon,
			Fields:     fields,
			Footer:     fmt.Sprintf("%s - %s", repositoryUrl, repositoryName),
			FooterIcon: repositoryIcon.String(),
			AuthorName: authorName.String(),
		})

	}

	// Check deployments
	for _, deployment := range cfg.Deployments {
		errRollout := rollout.TrackDeploymentTillReady(deployment, cfg.Namespace, kube.Kubernetes, tracker.Options{Timeout: cfg.Timout * time.Second})
		apiDeployment, errGetName := kube.Kubernetes.ExtensionsV1beta1().Deployments(cfg.Namespace).Get(deployment, v1.GetOptions{})
		var images []string
		footer := ""

		if errGetName == nil {

			for _, container := range apiDeployment.Spec.Template.Spec.Containers {
				images = append(images, fmt.Sprintf("%s: %s", container.Name, container.Image))
			}
			revision := apiDeployment.GetAnnotations()["deployment.kubernetes.io/revision"]
			footer = fmt.Sprintf("Replica Count: %d - Revision: %s", apiDeployment.Status.ReadyReplicas, revision)
		}

		if errRollout != nil {
			deploymentSuccess = false
			attachments = append(attachments, &slack.Attachment{
				Color:  "danger",
				Footer: footer,
				Title:  fmt.Sprintf("❌ %s / %s", cfg.Namespace, deployment),
				Text:   fmt.Sprintf("%s", errRollout),
			})
		} else {
			attachments = append(attachments, &slack.Attachment{
				Color:  "good",
				Footer: footer,
				Title:  fmt.Sprintf("✔️️️ %s / %s", cfg.Namespace, deployment),
				Text:   fmt.Sprint(strings.Join(images[:], "\n")),
			})
		}
	}

	// Damn it golang is verbose
	var Color string
	if deploymentSuccess {
		Color = "good"
	} else {
		Color = "danger"
	}

	// Set Color and attachment
	slackMsg.Color = Color
	slackMsg.Attachments = attachments

	if cfg.SlackUrl != "" {
		err = slack.Send(cfg.SlackUrl, slackMsg)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

	}

	if !deploymentSuccess {
		os.Exit(1)
	}

	os.Exit(0)
}
