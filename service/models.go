package service

// this is the struct that will represent the Task in the Github PR Description
type AsanaGithubDescriptionLineItem struct {
	TaskId    string `json:"task_id"`
	TaskUrl   string `json:"task_url"`
	Assignee  string `json:"assignee"`
	IsSubTask string `json:"is_sub_task"`
}
