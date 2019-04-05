package filesystem

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type (
	Filesystem interface {
		Create(User, File) error
		Remove(User, File) error
		Rename(User, File) error
		Lookup(User, File) (File, error)
	}

	File interface {
		Body() io.Reader
		Name() string
		Extension() string
	}

	User interface {
		ID() string
	}
)

func NewFilesystem(root string) Filesystem {
	return &defaultFilesystem{
		rootpath: root,
	}
}

type defaultFilesystem struct {
	rootpath string
}

func (fs *defaultFilesystem) Create(user User, file File) error {
	path, err := fs.userfolder(user)
	if err != nil {
		return err
	}

	return create(path, file)
}

func (fs *defaultFilesystem) Rename(user User, file File) error {
	return nil
}

func (fs *defaultFilesystem) Remove(user User, file File) error {
	return nil
}

func (fs *defaultFilesystem) Lookup(user User, file File) (File, error) {
	path, err := fs.userfolder(user)
	if err != nil {
		return nil, err
	}

	return lookup(path, file)
}

func (fs defaultFilesystem) userfolder(user User) (string, error) {
	userfold := folderNameForUser(user)
	pathto := filepath.Join(fs.rootpath, userfold)
	err := os.Mkdir(pathto, os.ModeDir)
	log.Println("pathto", pathto)

	return pathto, err
}

func create(pathto string, f File) error {
	pathto = filepath.Join(pathto, f.Name()+"."+f.Extension())
	_, err := os.Create(pathto)

	return err
}

func lookup(pathto string, f File) (File, error) {
	pathto = filepath.Join(pathto, f.Name())
	file, err := os.Open(pathto)
	if err != nil {
		return nil, err
	}

	return NewFile(file.Name(), strings.TrimRight(file.Name(), "."), file), nil
}