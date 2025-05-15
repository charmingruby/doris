package rest

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/charmingruby/doris/lib/fs"
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
	ErrFailedToOpenFile      = errors.New("failed to open file")
)

type HandleMultipartFormFilesInput struct {
	MaxFileSize  int64
	AllowedTypes string
	MaxFiles     int
}

type InvalidFile struct {
	File   fs.File
	Reason error
}

type HandleMultipartFormFilesOutput struct {
	ValidFiles   []fs.File
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

	var validFiles []fs.File
	var invalidFiles []InvalidFile

	for _, mpFile := range mpFiles {
		file, err := mpFile.Open()
		if err != nil {
			return nil, ErrFailedToOpenFile
		}
		defer file.Close()

		if mpFile.Size > in.MaxFileSize {
			invalidFiles = append(invalidFiles, InvalidFile{
				File:   fs.File{File: file, Filename: mpFile.Filename, Extension: filepath.Ext(mpFile.Filename)},
				Reason: ErrFileSizeExceeded,
			})
			continue
		}

		ext := strings.ToLower(filepath.Ext(mpFile.Filename))
		if !strings.Contains(in.AllowedTypes, ext) {
			invalidFiles = append(invalidFiles, InvalidFile{
				File:   fs.File{File: file, Filename: mpFile.Filename, Extension: ext},
				Reason: ErrFileTypeNotAllowed,
			})
			continue
		}

		validFiles = append(validFiles, fs.File{
			File:      file,
			Filename:  mpFile.Filename,
			Extension: ext,
		})
	}

	return &HandleMultipartFormFilesOutput{
		ValidFiles:   validFiles,
		InvalidFiles: invalidFiles,
	}, nil
}
