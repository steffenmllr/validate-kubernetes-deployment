package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/flant/kubedog/pkg/kube"
	"github.com/flant/kubedog/pkg/tracker"
	"github.com/flant/kubedog/pkg/trackers/rollout"
	"github.com/steffenmllr/validate-kubernetes-deployment/slack"
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

	// Check deployments
	for _, deployment := range cfg.Deployments {
		errRollout := rollout.TrackDeploymentTillReady(deployment, cfg.Namespace, kube.Kubernetes, tracker.Options{Timeout: cfg.Timout})
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
