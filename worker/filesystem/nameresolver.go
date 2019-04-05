package filesystem

import (
	"crypto/sha1"
	"encoding/base64"
	"path/filepath"
)

func folderNameForUser(user User) string {
	hasher := sha1.New()
	hasher.Write([]byte(user.ID()))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return filepath.Clean(sha + "_" + user.ID())
}
