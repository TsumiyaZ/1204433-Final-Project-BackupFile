package dto

type ScannedPhoto struct {
	Filename     string `json:"filename"`
	Path         string `json:"path"`
	RelativePath string `json:"relativePath"`
	Extension    string `json:"extension"`
	Size         int64  `json:"size"`
}
