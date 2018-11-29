package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
)

type tarGzip struct {
	file *os.File
	gz   *gzip.Writer
	arch *tar.Writer
}

// NewTarGzip returns a new tar.gz archive.
func NewTarGzip(filename string) (Archive, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}

	gz := gzip.NewWriter(f)
	arch := tar.NewWriter(gz)

	return &tarGzip{
		file: f,
		gz:   gz,
		arch: arch,
	}, nil
}

func (tgz tarGzip) Add(file string, content string) error {
	if tgz.arch == nil {
		return errors.New("archive has already been closed")
	}

	hdr := &tar.Header{
		Name: file,
		Mode: 0600,
		Size: int64(len(content)),
	}

	if err := tgz.arch.WriteHeader(hdr); err != nil {
		return fmt.Errorf("failed to write tar file header: %v", err)
	}

	if _, err := tgz.arch.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write tar file data: %v", err)
	}

	return nil
}

func (tgz tarGzip) Close() {
	if tgz.arch != nil {
		tgz.arch.Close()
		tgz.arch = nil
	}

	if tgz.gz != nil {
		tgz.gz.Close()
		tgz.gz = nil
	}

	if tgz.file != nil {
		tgz.file.Close()
		tgz.file = nil
	}
}