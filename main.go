package main

import (
	"fmt"

	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/go-playground/webhooks/v6/github"
)

const (
	path = "/github"
	// ghSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	ghSecret = "MyGitHubSuperSecretSecrect...?"
	reviewer = "some-team"
)

func handlePullRequest(payload github.PullRequestPayload) error {
	log.Info(fmt.Sprintf("Pull Request of type %s received", payload.Action))

	switch payload.Action {

	case "opened":
		log.Info(fmt.Sprintf("Assign PR to team %s", reviewer))

	default:
		log.Info(fmt.Sprintf("Pull Request hook of type not handled. Type = %s", payload.Action))
	}


}

func main() {
	hook, _ := github.New(github.Options.Secret(ghSecret))
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Info("Unhandled webhook received")
			} else {
				log.Info("Non-valid webhook received")
			}
		}
		switch payload.(type) {

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			err := handlePullRequest(pullRequest)
			if err != nil {
				log.Info("Failed to handle")
			}
		}
	})

	log.Info("Listening...")
	http.ListenAndServe(":3000", nil)
}
