package asana

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// all the app logic here
type ICore interface {
	// create a task
	// delete a task by id ?
	// mark task as done
	// send reminder to the assignee
	// add assignee to a task
	CreateTask(ctx context.Context, projectGId string, taskName string, assigneeEmail string) (*Task, error)
	DeleteTask(ctx context.Context, taskGid string) error
}

type impl struct{}

var core ICore

func SetCore() {
	core = &impl{}
}

func GetCore() ICore {
	return core
}

func (i *impl) CreateTask(ctx context.Context, projectGId string, taskName string, assigneeEmail string) (*Task, error) {
	log.Printf("creating task %s in project %s and assigning to %s\n", taskName, projectGId, assigneeEmail)

	url := fmt.Sprintf("%s/%s", BasePath, TaskPath)

	request := BaseRequest{
		Data: CreateTaskRequest{
			Name:     taskName,
			Projects: []string{projectGId},
		},
	}

	responseHolder := CreateTaskResponse{Data: Task{}}
	err := Call(url, request, &responseHolder, headers(), http.MethodPost)
	if err != nil {
		log.Println(err)
	}

	return &responseHolder.Data, nil
}

func (i *impl) DeleteTask(ctx context.Context, taskGid string) error {
	log.Printf("deleting task %s", taskGid)

	url := fmt.Sprintf("%s/%s/%s", BasePath, TaskPath, taskGid)

	responseHolder := CreateTaskResponse{Data: Task{}}
	err := Call(url, nil, &responseHolder, headers(), http.MethodDelete)
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
