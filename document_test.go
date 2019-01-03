package gitdoc

import (
	"fmt"
	"testing"
)

var (
	testDocID    = "test_id"
	testDocBody1 = "This is a content for test document."
	testDocBody2 = "This is a content for test document.\nAnd it's new!"
)

func TestDocument(t *testing.T) {
	/// ===== Test CreateDocument =====
	// create new document
	doc, err := testRepo.CreateDocument(testDocID, testDocBody1)
	if err != nil {
		fmt.Printf("CreateDocument failed: %v\n", err)
		t.Fail()
		return
	}

	if doc == nil ||
		doc.ID != testDocID ||
		doc.Body != testDocBody1 {
		fmt.Printf("CreateDocument failed: invalid data - %v\n", doc)
		t.Fail()
		return
	}

	// create document with existing id (should fail)
	_, err = testRepo.CreateDocument(testDocID, testDocBody1)
	if err == nil ||
		err != ErrDocumentExists {
		fmt.Println("CreateDocument failed: should return error")
	}

	/// ===== Test UpdateDocument =====
	// update document
	newDoc, err := testRepo.UpdateDocument(testDocID, testDocBody2)
	if err != nil {
		fmt.Printf("UpdateDocument failed: %v\n", err)
		t.Fail()
		return
	}

	if newDoc.ID != testDocID ||
		newDoc.Body != testDocBody2 {
		fmt.Printf("CreateDocument failed: invalid data - %v\n", newDoc)
		t.Fail()
		return
	}

	// update document that not exists (should fail)
	_, err = testRepo.UpdateDocument(testDocID+"_invalid", testDocBody1)
	if err == nil ||
		err != ErrDocumentNotExists {
		fmt.Println("UpdateDocument failed: should return error")
	}

	/// ===== Test GetDocument =====
	// get document
	getDoc, err := testRepo.GetDocument(testDocID)
	if err != nil {
		fmt.Printf("GetDocument failed: %v\n", err)
		t.Fail()
		return
	}

	if newDoc.ID != getDoc.ID ||
		newDoc.RevisionHash != getDoc.RevisionHash ||
		newDoc.Body != getDoc.Body ||
		newDoc.UpdatedAt.Equal(*getDoc.UpdatedAt) {
		fmt.Printf("GetDocument failed: document unmathced - %v", newDoc)
		t.Fail()
		return
	}

	// get document that not exists (should fail)
	_, err = testRepo.GetDocument(testDocID + "_invalid")
	if err == nil ||
		err != ErrGetDocument {
		fmt.Println("GetDocument failed: should return error")
	}
}
