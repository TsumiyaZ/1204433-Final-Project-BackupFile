package dto

const (
	MoveStatusMoved   = "moved"
	MoveStatusSkipped = "skipped"
	MoveStatusFailed  = "failed"
)

type MovePhotosRequest struct {
	Paths []string `json:"paths"`
}

type MoveItemResult struct {
	SourcePath      string `json:"sourcePath"`
	Filename        string `json:"filename"`
	DestinationPath string `json:"destinationPath"`
	Status          string `json:"status"`
	Message         string `json:"message"`
}

type MoveProgress struct {
	Current int            `json:"current"`
	Total   int            `json:"total"`
	Percent int            `json:"percent"`
	Item    MoveItemResult `json:"item"`
}

type MovePhotosResult struct {
	Requested  int              `json:"requested"`
	Moved      int              `json:"moved"`
	Skipped    int              `json:"skipped"`
	Failed     int              `json:"failed"`
	DurationMs int64            `json:"durationMs"`
	Items      []MoveItemResult `json:"items"`
}

type BackupError struct {
	Message string `json:"message"`
}
