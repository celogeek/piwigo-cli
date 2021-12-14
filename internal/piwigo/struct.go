package piwigo

type Piwigo struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

type PiwigoResult struct {
	Stat       string      `json:"stat"`
	Err        int         `json:"err"`
	ErrMessage string      `json:"message"`
	Result     interface{} `json:"result"`
}
