package gitdoc

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type Document struct {
	ID       string
	Body     string
	Revision DocumentRevision
}

type DocumentRevision struct {
	Hash      string
	UpdatedAt *time.Time
}

var (
	ErrDocumentExists    = errors.New("document already exists")
	ErrDocumentNotExists = errors.New("document not exists")

	ErrCreateDocument = errors.New("failed to create document")
	ErrGetDocument    = errors.New("failed to open document")

	ErrInvalidOption = errors.New("invalid option")
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
		ID:   id,
		Body: body,
		Revision: DocumentRevision{
			Hash:      hash.String(),
			UpdatedAt: &now,
		},
	}, nil
}

// GetDocument returns latest revision of the document.
func (r *Repository) GetDocument(id string) (*Document, error) {
	head, _ := r.gitRepo.Head()
	return r.GetDocumentAtRevision(id, head.Hash().String())
}

// GetDocumentIDs returns list of id for all documents.
// TODO - Implement pagination
func (r *Repository) GetDocumentIDs() ([]string, error) {
	files, err := ioutil.ReadDir(r.getDocumentPath(""))
	if err != nil {
		return nil, ErrGetDocument
	}

	var ids []string
	for _, f := range files {
		ids = append(ids, f.Name())
	}
	return ids, nil
}

// GetDocumentAtRevision returns specific revision of the document.
func (r *Repository) GetDocumentAtRevision(id string, revisionHash string) (*Document, error) {
	hash, _ := r.gitRepo.ResolveRevision(plumbing.Revision(revisionHash))
	commit, err := r.gitRepo.Object(plumbing.CommitObject, *hash)
	if err != nil {
		return nil, ErrGetDocument
	}

	tree, _ := commit.(*object.Commit).Tree()
	file, err := tree.File(r.getDocumentPathRel(id))
	if err != nil {
		return nil, ErrGetDocument
	}

	buf := new(bytes.Buffer)
	reader, _ := file.Blob.Reader()
	buf.ReadFrom(reader)

	return &Document{
		ID:   id,
		Body: buf.String(),
		Revision: DocumentRevision{
			Hash:      hash.String(),
			UpdatedAt: &commit.(*object.Commit).Author.When,
		},
	}, nil
}

// GetDocumentRevisions returns all revisions of the document.
// TODO - Implement pagination
func (r *Repository) GetDocumentRevisions(id string) ([]DocumentRevision, error) {
	fileName := r.getDocumentPathRel(id)
	iter, err := r.gitRepo.Log(&git.LogOptions{
		Order:    git.LogOrderCommitterTime,
		FileName: &fileName,
	})
	if err != nil {
		return nil, ErrGetDocument
	}

	var revs []DocumentRevision
	iter.ForEach(func(commit *object.Commit) error {
		revs = append(revs, DocumentRevision{
			Hash:      commit.Hash.String(),
			UpdatedAt: &commit.Author.When,
		})
		return nil
	})

	return revs, nil
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
		ID:   id,
		Body: newBody,
		Revision: DocumentRevision{
			Hash:      hash.String(),
			UpdatedAt: &now,
		},
	}, nil
}

// CompareOutputOption is a type for options to select the style of diffs.
type CompareOutputOption int8

const (
	// CompareOutputHTML is a option to present diffs as HTML string.
	CompareOutputHTML CompareOutputOption = iota

	// CompareOutputText is a option to present diffs as a colored string.
	// Recommended for terminal uses.
	CompareOutputText
)

// CompareDocumentRevisions show diffs of two revisions of the document.
func (r *Repository) CompareDocumentRevisions(id, revHashFrom, revHashTo string, typ CompareOutputOption) (string, error) {
	docFrom, err := r.GetDocumentAtRevision(id, revHashFrom)
	if err != nil {
		return "", err
	}
	docTo, err := r.GetDocumentAtRevision(id, revHashTo)
	if err != nil {
		return "", err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(docFrom.Body, docTo.Body, true)

	switch typ {
	case CompareOutputHTML:
		return dmp.DiffPrettyHtml(diffs), nil
	case CompareOutputText:
		return dmp.DiffPrettyText(diffs), nil
	default:
		return "", ErrInvalidOption
	}
}
