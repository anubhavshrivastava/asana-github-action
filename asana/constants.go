package asana

const (
	BasePath        = "https://app.asana.com/api/1.0"
	TaskPath        = "tasks"
	JsonContentType = "application/json"
)

var ReviewerType = struct {
	User string
}{
	User: "User",
}
