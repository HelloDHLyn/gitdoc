package gitdoc

import (
	"errors"
	"io"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Repository struct {
	gitRepo      *git.Repository
	gitRepoPath  string
	gitSignature *object.Signature
}

type InitOptions struct {
	Path string

	// Author** fields used as signature for commit.
	AuthorName  string // deafult value: "gitdoc"
	AuthorEmail string // default value: "-"
}

var (
	ErrInitNotEmpty = errors.New("path to initialized not empty")
	ErrInitFailed   = errors.New("failed to initialize repository")
)

func Init(opt *InitOptions) (*Repository, error) {
	f, _ := os.Open(opt.Path)
	if _, err := f.Readdirnames(1); err != io.EOF {
		return nil, ErrInitNotEmpty
	}

	// Make new directories.
	os.MkdirAll(opt.Path, 0744)
	os.Mkdir(opt.Path+"/docs", 0744)

	r, err := git.PlainInit(opt.Path, false)
	if err != nil {
		return nil, ErrInitFailed
	}

	return &Repository{
		gitRepo:     r,
		gitRepoPath: opt.Path,
		gitSignature: &object.Signature{
			Name:  stringOrDefault(opt.AuthorName, "gitdoc"),
			Email: stringOrDefault(opt.AuthorEmail, "-"),
		},
	}, nil
}
