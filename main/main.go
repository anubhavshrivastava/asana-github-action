package main

import (
	"fmt"
	"os"
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

func main() {
	name := os.Args[1]
	time := os.Args[2]
	fmt.Printf("::set-output name=time::%s\n", time)
	fmt.Printf("::debug Hey Man you said :%s\n", name)
}
