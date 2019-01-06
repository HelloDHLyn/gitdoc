package gitdoc

import (
	"fmt"
	"testing"
)

var (
	testDocID    = "test_id"
	testDocBody1 = "This is a content for test document."
	testDocBody2 = "This is a content for test document.\nAnd it's new!"

	testDocIDs = []string{"ID!@#$", "한글ID", "日本語ID", "ID With Whitespace"}
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

	// create new document with non-ascii id
	for _, id := range testDocIDs {
		_, err = testRepo.CreateDocument(id, testDocBody1)
		if err != nil {
			fmt.Printf("CreateDocument failed: %v\n", err)
			t.Fail()
			return
		}
	}

	// create document with existing id (should fail)
	_, err = testRepo.CreateDocument(testDocID, testDocBody1)
	if err == nil || err != ErrDocumentExists {
		fmt.Println("CreateDocument should return error")
		t.Fail()
		return
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
	if err == nil || err != ErrDocumentNotExists {
		fmt.Println("UpdateDocument should return error")
		t.Fail()
		return
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
		newDoc.Revision.Hash != getDoc.Revision.Hash ||
		newDoc.Body != getDoc.Body ||
		newDoc.Revision.UpdatedAt.Equal(*getDoc.Revision.UpdatedAt) {
		fmt.Printf("GetDocument failed: document unmathced - %v", newDoc)
		t.Fail()
		return
	}

	// get document that not exists (should fail)
	_, err = testRepo.GetDocument(testDocID + "_invalid")
	if err == nil || err != ErrGetDocument {
		fmt.Println("GetDocument should return error")
		t.Fail()
		return
	}

	/// ===== Test GetDocumentRevisions =====
	// get revisions
	revs, err := testRepo.GetDocumentRevisions(testDocID)
	if err != nil {
		fmt.Printf("GetDocuemtRevisions failed: %v\n", err)
		t.Fail()
		return
	}

	if len(revs) != 2 || revs[0].Hash != getDoc.Revision.Hash {
		fmt.Printf("GetDocumentRevisions return invalid data: %v\n", revs)
		t.Fail()
		return
	}

	// get revisions that not exists (should return empty slice)
	failRevs, _ := testRepo.GetDocumentRevisions(testDocID + "_invalid")
	if len(failRevs) > 0 {
		fmt.Println("GetDocumentRevisions should return empty slice")
		t.Fail()
		return
	}

	/// ===== Test GetDocumentAtRevision =====
	// get document at the last revision
	revDoc, err := testRepo.GetDocumentAtRevision(testDocID, revs[0].Hash)
	if err != nil {
		fmt.Printf("GetDocumentAtRevision failed: %v\n", err)
		t.Fail()
		return
	}
	if revDoc.ID != testDocID || revDoc.Revision.Hash != revs[0].Hash {
		fmt.Printf("GetDocumentRevisions return invalid data: %v\n", revDoc)
		t.Fail()
		return
	}

	// get document at the invalid revision (should fail)
	_, err = testRepo.GetDocumentAtRevision(testDocID, "invalid_hash")
	if err == nil || err != ErrGetDocument {
		fmt.Println("ErrGetDocument should return error")
		t.Fail()
		return
	}

	// get document with invalid id (should fail)
	_, err = testRepo.GetDocumentAtRevision(testDocID+"_invalid", revs[0].Hash)
	if err == nil || err != ErrGetDocument {
		fmt.Println("ErrGetDocument should return error")
		t.Fail()
		return
	}

	/// ===== Test GetDocumentIDs =====
	// get document ids
	ids, err := testRepo.GetDocumentIDs()
	if err != nil {
		fmt.Printf("GetDocumentIDs failed: %v\n", err)
		t.Fail()
		return
	}
	if len(ids) != 5 {
		fmt.Printf("GetDocumentIDs return invalid data: %v\n", ids)
		t.Fail()
		return
	}

	/// ===== Test CompareDocumentRevisions =====
	// compare revisions (with HTML output)
	diffs, err := testRepo.CompareDocumentRevisions(testDocID, revs[1].Hash, revs[0].Hash, CompareOutputHTML)
	expectedDiffs := "<span>This is a content for test document.</span><ins style=\"background:#e6ffe6;\">&para;<br>And it&#39;s new!</ins>"
	if diffs != expectedDiffs {
		fmt.Println("CompareDocumentRevisions returns invalid data")
		t.Fail()
		return
	}

	// compare revisions (with text output)
	diffs, err = testRepo.CompareDocumentRevisions(testDocID, revs[1].Hash, revs[0].Hash, CompareOutputText)
	expectedDiffs = "This is a content for test document.\x1b[32m\nAnd it's new!\x1b[0m"
	if diffs != expectedDiffs {
		fmt.Println("CompareDocumentRevisions returns invalid data")
		t.Fail()
		return
	}

	// compare revisions for invalid revision (should fail)
	_, err = testRepo.CompareDocumentRevisions(testDocID, revs[1].Hash+"_invalid", revs[0].Hash, CompareOutputHTML)
	if err == nil || err != ErrGetDocument {
		fmt.Println("CompareDocumentRevisions should return error")
		t.Fail()
		return
	}
	_, err = testRepo.CompareDocumentRevisions(testDocID, revs[1].Hash, revs[0].Hash+"_invalid", CompareOutputHTML)
	if err == nil || err != ErrGetDocument {
		fmt.Println("CompareDocumentRevisions should return error")
		t.Fail()
		return
	}
}
