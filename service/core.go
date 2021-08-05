package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/razorpay/asana-github-action/asana"
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

	if len(asanaTasksInBody) != 0 {
		// here you get the tasks from asana and load them up for more information
		// we may need to update their assignees or delete them, but that is later
	}

	// not you get the emails of these assignees and create the task and the sub-task on them
	// then we have to update the PR description too, not sure how we will manage the update without screwing the
	// rest of the body
	// var asanaTasks []AsanaGithubDescriptionLineItem
	var asanaLineItems []AsanaGithubDescriptionLineItem
	taskName := fmt.Sprintf("Review of PR: %s", pr.Title)
	for _, reviewer := range pr.RequestedReviewers {
		if reviewer.Type != asana.ReviewerType.User {
			continue
		}

		asanaUserEmail, ok := githubUserNameToAsanaEmail[reviewer.Login]
		if !ok {
			log.Printf("no email mapped to user %s, unable to create PR Review Task on Asana\n", reviewer.Login)
			continue
		}

		log.Printf("creating task %s, with assignee %s\n", taskName, asanaUserEmail)
		task, err := asana.GetCore().CreateTask(ctx, getConfig().AsanaProjectGid, taskName, asanaUserEmail)

		if err != nil {
			log.Printf("Task Creation failed with error %s\n", err.Error())
			continue
		}
		asanaLineItems = append(asanaLineItems, AsanaGithubDescriptionLineItem{
			TaskId:    task.Gid,
			TaskUrl:   task.PermaLinkUrl,
			Assignee:  asanaUserEmail,
			IsSubTask: "false",
		})
	}

	if len(asanaLineItems) == 0 {
		log.Printf("no asana tasks could be created. exiting")
		return
	}

	// Add Subtask on the reviewer, so that a Due-Date is added
	// Will do later, if this is really required

	// if they were created, then we add them to the PR description and update github
	asanaTaskPRDescription := generatePRDescriptionWithAsanaTasks(asanaLineItems)

	currentBody := pr.Body
	currentBody = currentBody + "\n" + asanaTaskPRDescription
	// update github PR Description
	pr, err = github.GetCore().UpdatePR(ctx, getConfig().PRLink, github.UpdatePR{
		Body: currentBody,
	})

	if err != nil {
		log.Printf("unable to update PR description due to %s", err.Error())
		return
	}
	log.Printf("updated PR with asana task description %s", asanaTaskPRDescription)
}

func generatePRDescriptionWithAsanaTasks(asanaTasks []AsanaGithubDescriptionLineItem) string {
	descriptionTemplate := "<asana>%s</asana>"
	tasksStr, _ := json.Marshal(asanaTasks)
	desc := fmt.Sprintf(descriptionTemplate, tasksStr)
	return desc
}

func parsePRBodyForAsanaTags(body string) ([]AsanaGithubDescriptionLineItem, error) {
	var tasks []AsanaGithubDescriptionLineItem
	if body == "" {
		return tasks, nil
	}

	return tasks, nil
}
