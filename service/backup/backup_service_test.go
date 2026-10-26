package backup

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"FinalProject/dto"
	"FinalProject/model"
)

var pngHeader = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

type fakeBackupRepo struct {
	mu     sync.Mutex
	photos []model.Photo
}

func (r *fakeBackupRepo) ChecksumExists(_ context.Context, dest string, checksum string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, photo := range r.photos {
		if photo.Dest == dest && photo.Checksum == checksum {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeBackupRepo) FilenameExists(_ context.Context, dest string, filename string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, photo := range r.photos {
		if photo.Dest == dest && photo.Filename == filename {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeBackupRepo) PhotosByChecksum(_ context.Context, dest string, checksum string) ([]*model.Photo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	photos := make([]*model.Photo, 0)
	for i := range r.photos {
		if r.photos[i].Dest == dest && r.photos[i].Checksum == checksum {
			photo := r.photos[i]
			photos = append(photos, &photo)
		}
	}
	return photos, nil
}

func (r *fakeBackupRepo) RestorePhoto(_ context.Context, id uint, dest string, needsAnalysis bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.photos {
		if r.photos[i].ID == id && r.photos[i].Dest == dest {
			r.photos[i].Status = "active"
			if needsAnalysis {
				r.photos[i].AnalysisStatus = model.AnalysisStatusPending
				r.photos[i].AnalysisError = ""
			}
			return nil
		}
	}
	return nil
}

func (r *fakeBackupRepo) CreatePhoto(_ context.Context, photo *model.Photo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if photo.ID == 0 {
		photo.ID = uint(len(r.photos) + 1)
	}
	r.photos = append(r.photos, *photo)
	return nil
}

type fakeSettingRepo struct {
	setting model.Setting
}

func (r *fakeSettingRepo) GetSetting(context.Context) (*model.Setting, error) {
	setting := r.setting
	return &setting, nil
}

func (r *fakeSettingRepo) SaveSetting(_ context.Context, setting *model.Setting) error {
	r.setting = *setting
	return nil
}

func writeFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func pngWith(body string) []byte {
	return append(append([]byte{}, pngHeader...), body...)
}

func setup(t *testing.T) (string, string, *fakeBackupRepo, BackupService) {
	t.Helper()
	source := t.TempDir()
	dest := t.TempDir()
	repo := &fakeBackupRepo{}
	settings := &fakeSettingRepo{setting: model.Setting{Source: source, Dest: dest}}
	return source, dest, repo, NewBackupService(repo, settings)
}

func TestMovePhotos(t *testing.T) {
	source, dest, repo, service := setup(t)

	good := filepath.Join(source, "a.png")
	nested := filepath.Join(source, "sub", "b.png")
	sameContent := filepath.Join(source, "sub", "a-copy.png")
	fake := filepath.Join(source, "fake.png")
	outside := filepath.Join(dest, "outside.png")

	writeFile(t, good, pngWith("one"))
	writeFile(t, nested, pngWith("two"))
	writeFile(t, sameContent, pngWith("one"))
	writeFile(t, fake, []byte("not an image"))
	writeFile(t, outside, pngWith("three"))

	var progressCount int
	result, err := service.MovePhotos(context.Background(), dto.MovePhotosRequest{
		Paths: []string{good, nested, sameContent, fake, outside, good},
	}, func(dto.MoveProgress) { progressCount++ })
	if err != nil {
		t.Fatal(err)
	}

	if result.Requested != 5 || result.Moved != 2 || result.Skipped != 1 || result.Failed != 2 {
		t.Fatalf("unexpected counts: %+v", result)
	}
	if progressCount != 5 {
		t.Fatalf("expected 5 progress events, got %d", progressCount)
	}
	if len(repo.photos) != 2 {
		t.Fatalf("expected 2 saved photos, got %d", len(repo.photos))
	}

	if _, err := os.Stat(nested); !os.IsNotExist(err) {
		t.Fatalf("source %s should be removed", nested)
	}

	// good and sameContent share content and run in parallel, so exactly
	// one of them is moved and the other stays in Source as a duplicate.
	_, goodErr := os.Stat(good)
	_, copyErr := os.Stat(sameContent)
	if os.IsNotExist(goodErr) == os.IsNotExist(copyErr) {
		t.Fatalf("expected exactly one of the duplicates to be moved")
	}

	for _, path := range []string{fake, outside} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("source %s should be kept: %v", path, err)
		}
	}
}

func TestMovePhotosRenamesOnNameConflict(t *testing.T) {
	source, dest, _, service := setup(t)

	photo := filepath.Join(source, "a.png")
	writeFile(t, photo, pngWith("new"))
	writeFile(t, filepath.Join(dest, "a.png"), pngWith("old"))

	result, err := service.MovePhotos(context.Background(), dto.MovePhotosRequest{
		Paths: []string{photo},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result.Moved != 1 || result.Items[0].Filename != "a (1).png" {
		t.Fatalf("expected rename to a (1).png, got %+v", result.Items[0])
	}
}

func TestMovePhotosSkipsDuplicateInDestination(t *testing.T) {
	source, dest, _, service := setup(t)

	photo := filepath.Join(source, "a.png")
	writeFile(t, photo, pngWith("same"))
	writeFile(t, filepath.Join(dest, "a.png"), pngWith("same"))

	result, err := service.MovePhotos(context.Background(), dto.MovePhotosRequest{
		Paths: []string{photo},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result.Skipped != 1 {
		t.Fatalf("expected duplicate to be skipped, got %+v", result.Items[0])
	}
	if _, err := os.Stat(photo); err != nil {
		t.Fatal("duplicate source should be kept")
	}
}

func TestCopyVerifyDelete(t *testing.T) {
	sourceDir := t.TempDir()
	destDir := t.TempDir()

	sourcePath := filepath.Join(sourceDir, "a.png")
	destPath := filepath.Join(destDir, "a.png")
	content := pngWith("cross drive")
	writeFile(t, sourcePath, content)

	checksum, err := fileChecksum(sourcePath)
	if err != nil {
		t.Fatal(err)
	}

	if err := copyVerifyDelete(sourcePath, destPath, int64(len(content)), "wrong"); err == nil {
		t.Fatal("expected checksum mismatch")
	}
	if _, err := os.Stat(sourcePath); err != nil {
		t.Fatal("source must be kept when verify fails")
	}
	if entries, _ := os.ReadDir(destDir); len(entries) != 0 {
		t.Fatal("temporary file must be removed when verify fails")
	}

	if err := copyVerifyDelete(sourcePath, destPath, int64(len(content)), checksum); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
		t.Fatal("source should be removed after verify")
	}
	if got, _ := fileChecksum(destPath); got != checksum {
		t.Fatal("destination checksum does not match")
	}
}

func TestDetectImageType(t *testing.T) {
	cases := map[string][]byte{
		"jpeg": {0xFF, 0xD8, 0xFF, 0xE0},
		"png":  pngHeader,
		"gif":  []byte("GIF89a"),
		"webp": []byte("RIFF\x00\x00\x00\x00WEBP"),
		"":     []byte("hello"),
	}
	for want, header := range cases {
		if got := detectImageType(header); got != want {
			t.Errorf("detectImageType(%q) = %q, want %q", header, got, want)
		}
	}
}
