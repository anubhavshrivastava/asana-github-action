package main

import (
	"log"
	"os"

	"github.com/razorpay/asana-github-action/asana"
	"github.com/razorpay/asana-github-action/github"
	"github.com/razorpay/asana-github-action/service"
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

func setupGithub(githubToken string) {
	github.SetConfig(github.Config{
		AccessToken: githubToken,
	})
}

func bootUp(asanaProjectGid string, prLink string, githubToken string, asanaToken string, userMapping string) {
	/* Setup Service */
	service.SetConfig(&service.Config{
		AsanaProjectGid: asanaProjectGid,
		PRLink:          prLink,
		UserMapping:     userMapping,
	})

	/* Setup Asana */
	asana.SetConfig(asana.Config{
		AccessToken: asanaToken,
	})
	asana.SetCore()

	/* Setup Github */
	setupGithub(githubToken)
	github.SetCore()
}

func main() {
	log.Printf("all my agrs brother: %v", os.Args)
	return
	args := os.Args

	log.Println(args)
	if len(args) <= 3 {
		panic("insufficient number of arguments, 3 required")
	}

	log.Println("assuming the args being sent in the order of 'Asana Project GID', 'PR Link', 'GITHUB_TOKEN', 'ASANA_TOKEN'")
	asanaProjectGid := args[1]
	prLink := args[2]
	githubToken := args[3]
	asanaToken := args[4]
	userMapping := args[5]

	bootUp(asanaProjectGid, prLink, githubToken, asanaToken, userMapping)

	// initiate processing
	service.ProcessPR()
}
