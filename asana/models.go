package asana

type BaseRequest struct {
	Data interface{} `json:"data"`
}

type CreateTaskResponse struct {
	Data Task
}

type CreateTaskRequest struct {
	Name     string   `json:"name"`
	Projects []string `json:"projects"`
}

type Task struct {
	Gid  string `json:"gid"`
	Name string `json:"name"`
}
