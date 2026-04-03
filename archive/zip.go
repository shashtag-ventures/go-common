package archive

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/mholt/archives"
)

// InMemFile implements fs.FileInfo for in-memory content.
type InMemFile struct {
	name    string
	size    int64
	modTime time.Time
}

func (f InMemFile) Name() string       { return filepath.Base(f.name) }
func (f InMemFile) Size() int64        { return f.size }
func (f InMemFile) Mode() fs.FileMode  { return 0644 }
func (f InMemFile) ModTime() time.Time { return f.modTime }
func (f InMemFile) IsDir() bool        { return false }
func (f InMemFile) Sys() interface{}   { return nil }

// BufferFile implements fs.File for in-memory content.
type BufferFile struct {
	*bytes.Reader
	info InMemFile
}

func (b BufferFile) Stat() (fs.FileInfo, error) { return b.info, nil }
func (b BufferFile) Close() error               { return nil }

// Archive creates a ZIP archive from disk paths and in-memory data.
// diskPaths maps a source directory/file on disk to its destination path in the archive.
// inMemFiles maps a destination path in the archive to its byte content.
func Archive(ctx context.Context, w io.Writer, diskPaths map[string]string, inMemFiles map[string][]byte, filter func(pathInArchive string) bool) error {
	var files []archives.FileInfo

	// Add files from disk
	if len(diskPaths) > 0 {
		diskFiles, err := archives.FilesFromDisk(ctx, nil, diskPaths)
		if err != nil {
			return err
		}
		for _, f := range diskFiles {
			if filter == nil || filter(f.NameInArchive) {
				files = append(files, f)
			}
		}
	}

	// Add in-memory files
	now := time.Now()
	for name, content := range inMemFiles {
		if filter == nil || filter(name) {
			contentCopy := content // capture for closure
			info := InMemFile{name: name, size: int64(len(content)), modTime: now}
			files = append(files, archives.FileInfo{
				NameInArchive: name,
				Open: func() (fs.File, error) {
					return &BufferFile{Reader: bytes.NewReader(contentCopy), info: info}, nil
				},
			})
		}
	}

	return archives.Zip{}.Archive(ctx, w, files)
}
