package service

import (
	"context"
	"log"

	"github.com/razorpay/asana-github-action/github"
)

/*
This will be the Orchestration Layer,
basically the whole business logic, and will treat Asana and Github as external communication modules
*/

func ProcessPR() {
	ctx := context.Background()
	pr, err := github.GetCore().GetPR(ctx, getConfig().PRLink)

	if err != nil {
		log.Printf("failed with error: %s\n", err.Error())
		return
	}

	if pr.Draft == true {
		log.Println("PR is in Draft state, exiting.")
		return
	}

	if len(pr.RequestedReviewers) == 0 {
		log.Println("no reviewers added, nothing to do, exiting.")
		return
	}

	// now that this is done, check the body
	// for now let this be empty, lets take that up in a bit
	asanaTasksInBody, err := parsePRBodyForAsanaTags(pr.Body)

}

func parsePRBodyForAsanaTags(body string) ([]AsanaGithubDescriptionLineItem, error) {
	var tasks []AsanaGithubDescriptionLineItem
	if body == "" {
		return tasks, nil
	}

	return tasks, nil
}
