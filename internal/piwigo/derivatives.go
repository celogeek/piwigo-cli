package piwigo

type Derivatives map[string]Derivative

type Derivative struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}
