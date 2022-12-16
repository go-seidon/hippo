package restapp

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type FileInfo struct {
	Name      string
	Size      int64
	Extension string
	Mimetype  string
}

/*
	@description: parse a given file and return file info
	@referrence: https://gist.github.com/rayrutjes/db9b9ea8e02255d62ce2?permalink_comment_id=3418419#gistcomment-3418419
*/
func ParseMultipartFile(file io.ReadSeeker, fh *multipart.FileHeader) (*FileInfo, error) {
	info := &FileInfo{Size: fh.Size}
	info.Name = ParseFileName(fh)
	info.Extension = ParseFileExtension(fh)

	buff := make([]byte, 512)
	n, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return nil, err
	}
	buff = buff[:n]

	info.Mimetype = http.DetectContentType(buff)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func ParseFileName(fh *multipart.FileHeader) string {
	names := strings.Split(fh.Filename, ".")
	return names[0]
}

func ParseFileExtension(fh *multipart.FileHeader) string {
	names := strings.Split(fh.Filename, ".")
	if len(names) == 1 {
		return ""
	}
	return names[len(names)-1]
}
