package github

type PR struct {
	Body               string     `json:"body"`
	RequestedReviewers []Reviewer `json:"requested_reviewers"`
	Draft              bool       `json:"draft"`
}

type UpdatePR struct {
	Body string `json:"body,omitempty"`
}

type Reviewer struct {
	Login string
	Type  string
}
