package piwigo

type ImagesDetails struct {
	Id            int         `json:"id"`
	Md5           string      `json:"md5sum"`
	Name          string      `json:"name"`
	DateAvailable TimeResult  `json:"date_available"`
	DateCreation  TimeResult  `json:"date_creation"`
	LastModified  TimeResult  `json:"lastmodified"`
	Width         int         `json:"width"`
	Height        int         `json:"height"`
	Url           string      `json:"page_url"`
	ImageUrl      string      `json:"element_url"`
	Filename      string      `json:"file"`
	Filesize      int64       `json:"filesize"`
	Categories    Categories  `json:"categories"`
	Tags          Tags        `json:"tags"`
	Derivatives   Derivatives `json:"derivatives"`
}
