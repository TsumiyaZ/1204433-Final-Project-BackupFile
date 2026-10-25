package dto

type GalleryPhoto struct {
	ID             uint     `json:"id"`
	Filename       string   `json:"filename"`
	Path           string   `json:"path"`
	Size           int64    `json:"size"`
	ModifiedAt     string   `json:"modifiedAt"`
	Description    string   `json:"description"`
	Status         string   `json:"status"`
	AnalysisStatus string   `json:"analysisStatus"`
	AnalysisError  string   `json:"analysisError"`
	Tags           []string `json:"tags"`
}
