package filesystem

import "io"

func NewFile(name, extension string, body io.Reader) File {
	return &file{body, name, extension}
}

type file struct {
	body      io.Reader
	name      string
	extension string
}

func (f *file) Body() io.Reader {
	return f.body
}
func (f *file) Name() string {
	return f.name
}
func (f *file) Extension() string {
	return f.extension
}
