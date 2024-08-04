package apprise

type Type string

const (
	Info    Type = "info"
	Success Type = "success"
	Warning Type = "warning"
	Failure Type = "failure"
)

type Notification struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body"`
	Type  Type   `json:"type,omitempty"`
	Tag   string `json:"tag,omitempty"`
}
