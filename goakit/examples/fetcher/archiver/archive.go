package archiver

import (
	"sync"
)

type (
	// Archive is the archived documents in-memory "DB"
	Archive struct {
		*sync.RWMutex
		docs []*Document
	}

	// Document represents a single archive document.
	Document struct {
		// Unique ID
		ID int
		// Status is the archived response HTTP status
		Status int
		// Body is the archive response HTTP body
		Body string
	}
)

// Store adds an archived document to the archive. Store takes care of
// initializing the document ID.
func (a *Archive) Store(doc *Document) {
	a.Lock()
	defer a.Unlock()
	id := len(a.docs) + 1
	doc.ID = id
	a.docs = append(a.docs, doc)
}

// Read retrieves an archived document by ID. It returns nil if there isn't one.
func (a *Archive) Read(id int) *Document {
	a.RLock()
	defer a.RUnlock()
	if id > len(a.docs) {
		return nil
	}
	return a.docs[id-1]
}
