package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}
type updateDescription struct {
	UpdatedFields bson.M   `bson:"updatedFields"`
	RemovedFields []string `bson:"removedFields"`
}
type ChangeEvent struct {
	ID                *changeID           `bson:"_id"`
	OperationType     string              `bson:"operationType"`
	ClusterTime       primitive.Timestamp `bson:"clusterTime"`
	UpdateDescription *updateDescription  `bson:"updateDescription"`
	DocumentKey       *documentKey        `bson:"documentKey"`
	Ns                *namespace          `bson:"ns"`
	FullDocument      *User               `bson:"fullDocument"`
}

// marshall event to an array of bytes
func (e ChangeEvent) Marshal() ([]byte, error) {
	return bson.MarshalExtJSON(e, true, true)
}

// return the document id of the event
func (e ChangeEvent) DocumentID() (string, error) {
	id := e.DocumentKey.ID
	if id.IsZero() {
		return "", fmt.Errorf("documentKey should not be empty")
	}
	return id.Hex(), nil
}
