package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

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
	assigneeToTaskMap, err := parsePRBodyForAsanaTags(pr.Body)

	if err != nil {
		log.Printf("some error in parsing body, exiting")
		return
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

		asanaUserEmail, ok := getUserMapping(getConfig().UserMapping)[reviewer.Login]
		if !ok {
			log.Printf("no email mapped to user %s, unable to create PR Review Task on Asana\n", reviewer.Login)
			continue
		}

		// check if the task for this is already there, in which case we skip this user
		if task, ok := assigneeToTaskMap[asanaUserEmail]; ok {
			log.Printf("asana task for %s is already present, skipping", asanaUserEmail)
			asanaLineItems = append(asanaLineItems, task)
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
	err = updateGithubPR(ctx, asanaTaskPRDescription, pr)

	if err != nil {
		log.Printf("unable to update PR description due to %s", err.Error())
		return
	}
	log.Printf("updated PR with asana task description %s", asanaTaskPRDescription)
}

func getUserMapping(something string) map[string]string {
	return githubUserNameToAsanaEmail
}

func updateGithubPR(ctx context.Context, asanaDesc string, pr *github.PR) error {
	currentBody := pr.Body
	re := regexp.MustCompile("<asana>.*</asana>")
	descAlreadyPresent := re.MatchString(currentBody)
	if descAlreadyPresent {
		currentBody = re.ReplaceAllString(currentBody, asanaDesc)
	} else {
		currentBody = currentBody + "\n" + asanaDesc
	}

	// update github PR Description
	_, err := github.GetCore().UpdatePR(ctx, getConfig().PRLink, github.UpdatePR{
		Body: currentBody,
	})

	return err
}

func generatePRDescriptionWithAsanaTasks(asanaTasks []AsanaGithubDescriptionLineItem) string {
	descriptionTemplate := "<asana>%s</asana>"
	var tasksStr []byte
	tasksStr, _ = json.Marshal(asanaTasks)
	desc := fmt.Sprintf(descriptionTemplate, string(tasksStr))
	return desc
}

func parsePRBodyForAsanaTags(body string) (map[string]AsanaGithubDescriptionLineItem, error) {
	log.Printf("parsing %s to extract asana information\n", body)
	re := regexp.MustCompile("<asana>(?P<json_body>.*)</asana>")
	match := re.FindStringSubmatch(body)
	if len(match) != 2 {
		// basically no data found
		return map[string]AsanaGithubDescriptionLineItem{}, nil
	}

	var response []AsanaGithubDescriptionLineItem
	err := json.Unmarshal([]byte(match[1]), &response)
	if err != nil {
		log.Printf("unable to understand the asana data: %v", match[1])
		return nil, err
	}

	// now we have that data, return that, or pre-process it ?
	var x = make(map[string]AsanaGithubDescriptionLineItem)
	for _, line := range response {
		x[line.Assignee] = line
	}
	return x, nil
}
