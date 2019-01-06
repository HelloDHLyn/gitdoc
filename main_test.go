package gitdoc

import (
	"os"
	"os/user"
	"testing"
)

var testRepo *Repository

func TestMain(m *testing.M) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	path := user.HomeDir + "/.gitdoc"
	initOpt := InitOptions{
		Path: path,
	}

	r, err := Init(&initOpt)
	if err != nil {
		panic(err)
	}

	testRepo = r
	os.Exit(m.Run())
}
