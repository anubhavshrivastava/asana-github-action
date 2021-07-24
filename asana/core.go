package asana

import (
	"context"
	"fmt"
	"log"
)

// all the app logic here
type ICore interface {
	// create a task
	// delete a task by id ?
	// mark task as done
	// send reminder to the assignee
	// add assignee to a task
	CreateTask(ctx context.Context, projectGId string, taskName string, assigneeEmail string) (*Task, error)
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

	url := fmt.Sprintf("%s/%s", BasePath, CreateTaskPath)

	request := BaseRequest{
		Data: CreateTaskRequest{
			Name:     taskName,
			Projects: []string{projectGId},
		},
	}

	responseHolder := CreateTaskResponse{Data: Task{}}
	err := Post(url, request, &responseHolder, headers())
	if err != nil {
		log.Println(err)
	}

	return &responseHolder.Data, nil
}

func headers() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", GetConfig().AccessToken),
		"Content-Type":  JsonContentType,
	}
}
