package asana

type BaseRequest struct {
	Data interface{} `json:"data"`
}

type TaskObjectResponse struct {
	Data Task
}

type CreateTaskRequest struct {
	Name     string   `json:"name"`
	Projects []string `json:"projects,omitempty"`
	Assignee string `json:"assignee,omitempty"`
}

type UpdateTaskRequest struct {
	Completed string `json:"completed,omitempty"` // true / false
	Assignee string `json:"assignee,omitempty"` // true / false
}

type Task struct {
	Gid       string `json:"gid"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
