package piwigo

import "sync"

type Piwigo struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`

	mu sync.Mutex
}
