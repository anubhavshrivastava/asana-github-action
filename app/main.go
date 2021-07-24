package main

import (
	"context"
	"fmt"
	"log"

	"github.com/razorpay/asana-github-action/asana"
)

/*
* you have to set the asana credentials via the github action itself
* those will be passed to this, and eventually via command line args would be sent here
* see if you can use named command line args for golang as well as github-action / entrypoint.sh
* you will need the
	* PR endpoint
	* github action
	* Reviewers
    * PR Author
    * PR Description
    * PR State (Draft or not)
    * All these when the PR is
		* Created
        * changes state to draft
        * closed --> not sure if we can delete the task or not, lets see if we can, depends on asana
*/

func setupAsana() {
	asana.SetConfig(asana.Config{
		ProjectId:   "1200653787952913",
		AccessToken: "",
	})

	asana.SetCore()
}

func bootUp() {
	setupAsana()
}

func main() {
	bootUp()

	task, err := asana.GetCore().CreateTask(context.Background(), asana.GetConfig().ProjectId, "anubhav-this-works", "")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(task)
}
