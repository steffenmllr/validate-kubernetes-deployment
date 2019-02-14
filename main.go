package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/flant/kubedog/pkg/kube"
	"github.com/flant/kubedog/pkg/tracker"
	"github.com/flant/kubedog/pkg/trackers/rollout"
	"github.com/steffenmllr/validate-kubernetes-deployment/slack"
)

type config struct {
	Name        string        `env:"NAME"`
	Link        string        `env:"LINK"`
	Namespace   string        `env:"NAMESPACE,required"`
	Deployments []string      `env:"DEPLOYMENTS,required" envSeparator:","`
	Timout      time.Duration `env:"TIMEOUT" envDefault:"1s"`
	SlackUrl    string        `env:"SLACK_HOOK_URL"`
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

	// Check deployments
	for _, deployment := range cfg.Deployments {
		errRollout := rollout.TrackDeploymentTillReady(deployment, cfg.Namespace, kube.Kubernetes, tracker.Options{Timeout: cfg.Timout * time.Second})

		if errRollout != nil {
			deploymentSuccess = false
			attachments = append(attachments, &slack.Attachment{
				Color: "danger",
				Title: fmt.Sprintf("❌ %s / %s", cfg.Namespace, deployment),
				Text:  fmt.Sprintf("%s", errRollout),
			})
		} else {
			attachments = append(attachments, &slack.Attachment{
				Color: "good",
				Title: fmt.Sprintf("✔️%s / %s", cfg.Namespace, deployment),
			})
		}
	}

	var Color string
	if deploymentSuccess {
		Color = "good"
	} else {
		Color = "danger"
	}

	slackMsg := &slack.Message{
		Text:        fmt.Sprintf("*Deployment Status*: %s \n %s", cfg.Name, cfg.Link),
		LinkNames:   "LinkNames",
		Color:       Color,
		Attachments: attachments,
	}

	if cfg.SlackUrl != "" {
		_ = slack.Send(cfg.SlackUrl, slackMsg)
	}

	if !deploymentSuccess {
		os.Exit(1)
	}
}
