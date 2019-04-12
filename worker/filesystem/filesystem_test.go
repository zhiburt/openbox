package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		user     User
		file     File
		expected error
	}{
		{user: NewUser("1"), file: NewFile("1", "go", strings.NewReader("")), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "go", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "", nil), expected: nil},
	}

	rootpath := "."
	fs := NewFilesystem(rootpath)
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			defer os.RemoveAll(filepath.Join(rootpath, folderNameForUser(c.user)))
			if err := fs.Create(c.user, c.file); err != c.expected {
				t.Errorf("FIX FIX FIX\nexpected %v\nbut was %v", c.expected, err)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	var patherror = fmt.Errorf("path erorr")
	cases := []struct {
		user     User
		file     File
		search   File
		expected error
	}{
		{user: NewUser("1"), file: NewFile("1", "go", strings.NewReader("")), search: NewFile("1", "go", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "go", nil), search: NewFile("1", "go", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "", nil), search: NewFile("1", "", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "", nil), search: NewFile("2", "", nil), expected: patherror},
	}

	rootpath := "."
	fs := NewFilesystem(rootpath)
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			defer os.RemoveAll(filepath.Join(rootpath, folderNameForUser(c.user)))
			if err := fs.Create(c.user, c.file); err != nil {
				t.Errorf("FIX FIX FIX\n was happend unexpected error %v", err)
			}
			if f, err := fs.Lookup(c.user, c.search); err != c.expected {
				if _, ok := err.(*os.PathError); ok && c.expected == patherror {

				} else {
					t.Errorf("FIX FIX FIX\nexpected error %v\nbut was %v", c.expected, err)
				}
			} else if f.Name() != c.file.Name() || f.Extension() != c.file.Extension() {
				t.Errorf("FIX FIX FIX\nexpected file%v\nbut was %v", c.file, f)
			}
		})
	}
}

func TestRename(t *testing.T) {
	cases := []struct {
		user     User
		file     File
		new      File
		expected error
	}{
		{user: NewUser("1"), file: NewFile("1", "go", strings.NewReader("")), new: NewFile("1", "go", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "go", nil), new: NewFile("2", "", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "", nil), new: NewFile("2", "go", nil), expected: nil},
	}

	rootpath := "."
	fs := NewFilesystem(rootpath)
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			defer os.RemoveAll(filepath.Join(rootpath, folderNameForUser(c.user)))
			if err := fs.Create(c.user, c.file); err != nil {
				t.Errorf("FIX FIX FIX\n was happend unexpected error %v", err)
			}
			if err := fs.Rename(c.user, c.file, c.new); err != c.expected {
				t.Errorf("FIX FIX FIX\nexpected %v\nbut was %v", c.expected, err)
			}
			if f, err := fs.Lookup(c.user, c.new); err != nil || f.Name() != c.new.Name() || f.Extension() != c.new.Extension() {
				t.Errorf("FIX FIX FIX\nfile wasn't changed %#v\nbut was %#v", c.new, f)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		user     User
		file     File
		expected error
	}{
		{user: NewUser("1"), file: NewFile("1", "go", strings.NewReader("")), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "go", nil), expected: nil},
		{user: NewUser("1"), file: NewFile("1", "", nil), expected: nil},
	}

	rootpath := "."
	fs := NewFilesystem(rootpath)
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			defer os.RemoveAll(filepath.Join(rootpath, folderNameForUser(c.user)))
			if err := fs.Create(c.user, c.file); err != nil {
				t.Errorf("FIX FIX FIX\n was happend unexpected error %v", err)
			}
			if err := fs.Remove(c.user, c.file); err != c.expected {
				t.Errorf("FIX FIX FIX\nexpected %v\nbut was %v", c.expected, err)
			}
			if _, err := fs.Lookup(c.user, c.file); err == nil {
				t.Errorf("FIX FIX FIX\nfile wasn't removed")
			}
		})
	}
}
