package dto

type GalleryPhoto struct {
	ID             uint   `json:"id"`
	Filename       string `json:"filename"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	AnalysisStatus string `json:"analysisStatus"`
}
