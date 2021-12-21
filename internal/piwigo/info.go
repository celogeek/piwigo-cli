package piwigo

type Infos []Info

type Info struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
