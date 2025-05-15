package rest

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/gin-gonic/gin"
)

const (
	FILE_UPLOAD_FORM_KEY = "files"
)

var (
	ErrFailedToParseFormData = errors.New("failed to parse form data")
	ErrNoFilesProvided       = errors.New("no files provided")
	ErrFileSizeExceeded      = errors.New("file size exceeded")
	ErrMaxFilesExceeded      = errors.New("maximum files exceeded")
	ErrFileTypeNotAllowed    = errors.New("file type not allowed")
)

type HandleMultipartFormFilesInput struct {
	MaxFileSize  int64
	AllowedTypes string
	MaxFiles     int
}

type InvalidFile struct {
	Filename string
	Reason   error
}

type HandleMultipartFormFilesOutput struct {
	ValidFiles   []*multipart.FileHeader
	InvalidFiles []InvalidFile
}

func HandleMultipartFormFiles(
	c *gin.Context,
	logger *instrumentation.Logger,
	in HandleMultipartFormFilesInput,
) (*HandleMultipartFormFilesOutput, error) {
	if err := c.Request.ParseMultipartForm(in.MaxFileSize); err != nil {
		return nil, ErrFailedToParseFormData
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, ErrFailedToParseFormData
	}

	mpFiles := form.File[FILE_UPLOAD_FORM_KEY]
	if len(mpFiles) == 0 {
		return nil, ErrNoFilesProvided
	}

	if len(mpFiles) > in.MaxFiles {
		return nil, ErrMaxFilesExceeded
	}

	var validFiles []*multipart.FileHeader
	var invalidFiles []InvalidFile

	for _, file := range mpFiles {
		if file.Size > in.MaxFileSize {
			invalidFiles = append(invalidFiles, InvalidFile{
				Filename: file.Filename,
				Reason:   ErrFileSizeExceeded,
			})
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !strings.Contains(in.AllowedTypes, ext) {
			invalidFiles = append(invalidFiles, InvalidFile{
				Filename: file.Filename,
				Reason:   ErrFileTypeNotAllowed,
			})
			continue
		}

		validFiles = append(validFiles, file)
	}

	return &HandleMultipartFormFilesOutput{
		ValidFiles:   validFiles,
		InvalidFiles: invalidFiles,
	}, nil
}
