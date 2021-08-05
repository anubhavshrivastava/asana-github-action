package asana

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/razorpay/asana-github-action/httpclient"
)

// all the app logic here
type ICore interface {
	CreateTask(ctx context.Context, projectGId string, taskName string, assigneeEmail string) (*Task, error)
	DeleteTask(ctx context.Context, taskGid string) error
	CompleteTask(ctx context.Context, taskGid string) error
	AddAssignee(ctx context.Context, taskGid string, email string) error
	CreateSubtask(ctx context.Context, parentTaskGid string, subtaskName string, subtaskAssigneeEmail string) error
}

type impl struct{}

var core ICore

func SetCore() {
	core = &impl{}
}

func GetCore() ICore {
	return core
}

func (i *impl) AddAssignee(ctx context.Context, taskGid string, email string) error {
	log.Printf("adding %s as assignee to task %s\n", email, taskGid)

	if email == "" {
		log.Printf("empty email, cannot proceed")
	}

	url := fmt.Sprintf("%s/%s/%s", BasePath, TaskPath, taskGid)

	responseHolder := TaskObjectResponse{Data: Task{}}

	request := BaseRequest{
		Data: UpdateTaskRequest{
			Assignee: email,
		},
	}
	err := httpclient.Call(url, request, &responseHolder, headers(), http.MethodPut)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (i *impl) CreateSubtask(ctx context.Context, parentTaskGid string, subtaskName string, subtaskAssigneeEmail string) error {
	log.Printf("creating subtask %s with assignee %s for the task %s\n", subtaskName, subtaskAssigneeEmail, parentTaskGid)

	url := fmt.Sprintf("%s/%s/%s/subtasks", BasePath, TaskPath, parentTaskGid)

	responseHolder := TaskObjectResponse{Data: Task{}}

	request := BaseRequest{
		Data: CreateTaskRequest{
			Name:     subtaskName,
			Assignee: subtaskAssigneeEmail,
		},
	}
	err := httpclient.Call(url, request, &responseHolder, headers(), http.MethodPost)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (i *impl) CompleteTask(ctx context.Context, taskGid string) error {
	log.Printf("marking task %s complete", taskGid)

	url := fmt.Sprintf("%s/%s/%s", BasePath, TaskPath, taskGid)

	responseHolder := TaskObjectResponse{Data: Task{}}

	request := BaseRequest{
		Data: UpdateTaskRequest{
			Completed: "true",
		},
	}
	err := httpclient.Call(url, request, &responseHolder, headers(), http.MethodPut)
	if err != nil {
		log.Println(err)
	}

	if responseHolder.Data.Completed == true {
		return nil
	}
	return fmt.Errorf("couldn't mark task completed")
}

func (i *impl) CreateTask(ctx context.Context, projectGId string, taskName string, assigneeEmail string) (*Task, error) {
	log.Printf("creating task %s in project %s and assigning to %s\n", taskName, projectGId, assigneeEmail)

	url := fmt.Sprintf("%s/%s", BasePath, TaskPath)

	request := BaseRequest{
		Data: CreateTaskRequest{
			Name:     taskName,
			Projects: []string{projectGId},
			Assignee: assigneeEmail,
		},
	}

	responseHolder := TaskObjectResponse{Data: Task{}}
	err := httpclient.Call(url, request, &responseHolder, headers(), http.MethodPost)
	if err != nil {
		log.Println(err)
	}

	return &responseHolder.Data, nil
}

func (i *impl) DeleteTask(ctx context.Context, taskGid string) error {
	log.Printf("deleting task %s", taskGid)

	url := fmt.Sprintf("%s/%s/%s", BasePath, TaskPath, taskGid)

	responseHolder := TaskObjectResponse{Data: Task{}}
	err := httpclient.Call(url, nil, &responseHolder, headers(), http.MethodDelete)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func headers() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", GetConfig().AccessToken),
		"Content-Type":  JsonContentType,
	}
}
