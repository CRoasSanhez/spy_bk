package models

import (
	"errors"
	"spyc_backend/app"

	"github.com/Reti500/mgomap"

	"gopkg.in/mgo.v2/bson"
)

// ErrDocument returns error when interface not is a mgomap.DocumentBase
var ErrDocument = errors.New("Invalid document name")

// ErrID returns error when id is not a Mongo ObjectID
var ErrID = errors.New("Invalid document")

// GetDocument returns a document searching by ID
func GetDocument(id string, out interface{}) error {
	if !bson.IsObjectIdHex(id) {
		return ErrID
	}

	if Document, ok := app.Mapper.GetModel(out); ok {
		return Document.Find(id).Exec(out)
	}

	return ErrDocument
}

// GetDocumentBy returns a document searching by query
func GetDocumentBy(query, out interface{}) error {
	if Document, ok := app.Mapper.GetModel(out); ok {
		return Document.Query(query).Exec(out)
	}

	return ErrDocument
}

// Query returns prepare query for exec
func Query(model string, query interface{}) mgomap.Query {
	Document := app.Mapper.InitModel(model)

	return Document.Query(query)
}

// CreateDocument save new document on DB
func CreateDocument(model interface{}) error {
	if Document, ok := app.Mapper.GetModel(model); ok {
		return Document.Create(model)
	}

	return ErrDocument
}

// UpdateDocument save full document on DB
func UpdateDocument(model interface{}) error {
	if Document, ok := app.Mapper.GetModel(model); ok {
		return Document.Update(model)
	}

	return ErrDocument
}

// UpdateByQuery is for custom update
func UpdateByQuery(model, selector, query interface{}, multi bool) error {
	if Document, ok := app.Mapper.GetModel(model); ok {
		return Document.UpdateQuery(selector, query, multi)
	}

	return ErrDocument
}

// SetDocument create a update query with $set paramater
// fields are a bson.M
func SetDocument(model, id string, fields interface{}) error {
	Document := app.Mapper.InitModel(model)

	var selector = bson.M{"_id": bson.ObjectIdHex(id)}
	var query = bson.M{"$set": fields}
	return Document.UpdateQuery(selector, query, false)
}

// NumberOfDocuments returns number of exists documents
func NumberOfDocuments(model, query interface{}) (int, error) {
	if Document, ok := app.Mapper.GetModel(model); ok {
		return Document.Query(query).Count()
	}

	return -1, ErrDocument
}

// AggregateQuery make and aggregate
func AggregateQuery(model, pipe, out interface{}) error {
	if Document, ok := app.Mapper.GetModel(model); ok {
		return Document.Pipe(pipe, out)
	}

	return ErrDocument
}

func AsDocumentBase(document interface{}) mgomap.DocumentInterface {
	return document.(mgomap.DocumentInterface)
}
