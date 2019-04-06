package commands

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"

	"github.com/openbox/worker/communication"
	"github.com/openbox/worker/filesystem"
)

type command func(filesystem.Filesystem, communication.Message) ([]byte, error)

var ErrNotFoundSuchCommand = errors.New("not found this command")

func NewCommand(fs filesystem.Filesystem, mssg communication.Message) (command, error) {
	m := map[string]command{
		"lookup": lookupfileCommand,
		"create": createfileCommand,
	}

	if foo, ok := m[mssg.Type]; ok {
		return foo, nil
	}

	return nil, ErrNotFoundSuchCommand
}

func lookupfileCommand(fs filesystem.Filesystem, m communication.Message) ([]byte, error) {
	log.Println("look up")
	f, err := fs.Lookup(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
	if err != nil {
		log.Println("DONT Found")
		return nil, err
	}

	log.Println("FILE FOUND")

	return ioutil.ReadAll(f.Body())
}

func createfileCommand(fs filesystem.Filesystem, m communication.Message) ([]byte, error) {
	log.Println("Create")
	return nil, fs.Create(filesystem.NewUser(m.UserID), filesystem.NewFile(m.Name, m.Extension, bytes.NewReader(m.Body)))
}
