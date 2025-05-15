package fs

import "io"

type File struct {
	File      io.ReadCloser
	Filename  string
	Extension string
}
