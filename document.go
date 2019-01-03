package gitdoc

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Document struct {
	ID           string
	RevisionHash string
	Body         string
	UpdatedAt    *time.Time
}

var (
	ErrDocumentExists    = errors.New("document already exists")
	ErrDocumentNotExists = errors.New("document not exists")

	ErrCreateDocument = errors.New("failed to create document")
	ErrGetDocument    = errors.New("failed to open document")
)

func (r *Repository) getDocumentPath(id string) string {
	return r.gitRepoPath + "/docs/" + id
}

func (r *Repository) getDocumentPathRel(id string) string {
	return "docs/" + id
}

// CreateDocument makes new file in the repository. It returns ErrDocumentExists
// if given ID already exists.
func (r *Repository) CreateDocument(id string, body string) (*Document, error) {
	if fileExists(r.getDocumentPath(id)) {
		return nil, ErrDocumentExists
	}

	// Write new document into the file.
	err := ioutil.WriteFile(r.getDocumentPath(id), []byte(body), 0644)
	if err != nil {
		fmt.Println(err)
		return nil, ErrCreateDocument
	}

	// Add and commit to the repository.
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return nil, ErrCreateDocument
	}

	now := time.Now()
	_, err = w.Add(r.getDocumentPathRel(id))
	if err != nil {
		return nil, ErrCreateDocument
	}
	hash, err := w.Commit(now.Format("20060102150405"), &git.CommitOptions{Author: r.gitSignature})
	if err != nil {
		return nil, ErrCreateDocument
	}

	return &Document{
		ID:           id,
		RevisionHash: hash.String(),
		Body:         body,
		UpdatedAt:    &now,
	}, nil
}

// GetDocument returns latest revision of the document.
func (r *Repository) GetDocument(id string) (*Document, error) {
	head, _ := r.gitRepo.Head()
	return r.GetDocumentAtRevision(id, head.Hash().String())
}

// GetDocumentAtRevision returns specific revision of the document.
func (r *Repository) GetDocumentAtRevision(id string, revisionHash string) (*Document, error) {
	hash, _ := r.gitRepo.ResolveRevision(plumbing.Revision(revisionHash))
	commit, _ := r.gitRepo.Object(plumbing.CommitObject, *hash)

	tree, _ := commit.(*object.Commit).Tree()
	file, err := tree.File(r.getDocumentPathRel(id))
	if err != nil {
		return nil, ErrGetDocument
	}

	buf := new(bytes.Buffer)
	reader, _ := file.Blob.Reader()
	buf.ReadFrom(reader)

	return &Document{
		ID:           id,
		RevisionHash: hash.String(),
		Body:         buf.String(),
		UpdatedAt:    &commit.(*object.Commit).Author.When,
	}, nil
}

// UpdateDocument update the body of existing document and create new revision.
// It returns ErrDocumentNotExists if given ID doesn't exist.
func (r *Repository) UpdateDocument(id string, newBody string) (*Document, error) {
	if !fileExists(r.getDocumentPath(id)) {
		return nil, ErrDocumentNotExists
	}

	// Write new document into the file.
	err := ioutil.WriteFile(r.getDocumentPath(id), []byte(newBody), 0644)
	if err != nil {
		fmt.Println(err)
		return nil, ErrCreateDocument
	}

	// Add and commit to the repository.
	w, err := r.gitRepo.Worktree()
	if err != nil {
		return nil, ErrCreateDocument
	}

	now := time.Now()
	_, err = w.Add(r.getDocumentPathRel(id))
	if err != nil {
		return nil, ErrCreateDocument
	}
	hash, err := w.Commit(now.Format("20060102150405"), &git.CommitOptions{Author: r.gitSignature})
	if err != nil {
		return nil, ErrCreateDocument
	}

	return &Document{
		ID:           id,
		RevisionHash: hash.String(),
		Body:         newBody,
		UpdatedAt:    &now,
	}, nil
}
