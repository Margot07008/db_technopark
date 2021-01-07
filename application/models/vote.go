package models

type Vote struct {
	Nickname   string `json:"nickname"`
	Voice      int32  `json:"voice"`
	Thread     int64  `json:"thread,omitempty"`
	ThreadSlug string `json:"thread_slug,omitempty"`
}
