package filesystem

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type (
	Filesystem interface {
		Create(User, File) error
		Remove(User, File) error
		Rename(User, File, File) error
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
	path, _ := fs.userfolder(user)

	return create(path, file)
}

func (fs *defaultFilesystem) Rename(user User, file, newfile File) error {
	path, _ := fs.userfolder(user)

	return rename(path, file, newfile)
}

func (fs *defaultFilesystem) Remove(user User, file File) error {
	path, _ := fs.userfolder(user)

	return remove(path, file)
}

func (fs *defaultFilesystem) Lookup(user User, file File) (File, error) {
	path, _ := fs.userfolder(user)

	return lookup(path, file)
}

func (fs defaultFilesystem) userfolder(user User) (string, error) {
	userfold := folderNameForUser(user)
	pathto := filepath.Join(fs.rootpath, userfold)
	err := os.Mkdir(pathto, os.ModePerm)
	log.Println("pathto", pathto)
	return pathto, err
}

func rename(pt string, f, f1 File) error {
	oldpath := pathto(pt, f.Name(), f.Extension())
	newpath := pathto(pt, f1.Name(), f1.Extension())

	return os.Rename(oldpath, newpath)
}

func remove(pt string, f File) error {
	pt = pathto(pt, f.Name(), f.Extension())

	return os.Remove(pt)
}

func create(pt string, f File) error {
	pt = pathto(pt, f.Name(), f.Extension())
	fmt.Println("pathto", pt)

	file, err := os.Create(pt)
	if err != nil {
		return err
	}
	if f.Body() == nil {
		return nil
	}

	defer file.Close()

	reader := io.TeeReader(f.Body(), file)
	_, err = ioutil.ReadAll(reader)
	return err
}

func lookup(pt string, f File) (File, error) {
	pt = pathto(pt, f.Name(), f.Extension())
	file, err := os.Open(pt)
	if err != nil {
		return nil, err
	}

	return NewFile(f.Name(), f.Extension(), file), nil
}

func pathto(p, name, extension string) string {
	if extension != "" {
		name += "." + extension
	}

	return filepath.Join(p, name)
}
