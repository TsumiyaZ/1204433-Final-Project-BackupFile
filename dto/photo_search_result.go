package dto

type PhotoSearchResult struct {
	ID       uint    `json:"id"`
	Filename string  `json:"filename"`
	Score    float64 `json:"score"`
}