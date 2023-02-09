package files_s3

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/infrastructure/files_s3"

	"github.com/google/uuid"
)

// s estructura de conexión s3
type s3 struct {
	TxID string
}

func newDocumentFileS3Repository(txID string) *s3 {
	return &s3{
		TxID: txID,
	}
}

func (s *s3) upload(id string, file *File) (*File, error) {
	c := env.NewConfiguration()
	var fullPath strings.Builder
	if file.Encoding == "" || file.OriginalFile == "" {
		return file, fmt.Errorf("couldn't create encoded file does not exist")
	}
	fl, err := base64.StdEncoding.DecodeString(file.Encoding)
	if err != nil {
		return file, err
	}
	file.Encoding = ""
	if fl == nil {
		return file, fmt.Errorf("couldn't create encoded file is null")
	}
	r := bytes.NewReader(fl)
	file.Path, file.FileName = s.getFullPath(file.OriginalFile, id)
	fullPath.WriteString(file.Path)
	fullPath.WriteString(file.FileName)
	file.Hash = s.getHashFromFile(fl)
	file.FileSize = int(r.Size())
	//TODO getNumberPage
	file.NumberPage = 1
	file.Bucket = c.Files.S3.Bucket
	err = files_s3.UploadFile(r, fullPath.String(), file.Bucket)
	if err != nil {
		return file, err
	}
	return file, nil
}

func (s *s3) getHashFromFile(file []byte) string {
	h := sha256.Sum256(file)
	return fmt.Sprintf("%x", h)
}

func (s *s3) getFullPath(originalFile string, id string) (string, string) {
	fPath := fmt.Sprintf("/%s/%d/%d/%d/%d/", id, time.Now().Year(), time.Now().YearDay(), time.Now().Hour(), time.Now().Minute())
	fileName := fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(originalFile))
	return fPath, fileName
}

func (s *s3) getFile(bucket, path, fileName string) (string, error) {
	return files_s3.GetObjectS3(bucket, path, fileName)
}
