package neo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResultDataContent string

var (
	ContentJustRow     = []ResultDataContent{"row"}
	ContentJustGraph   = []ResultDataContent{"graph"}
	ContentRowAndGraph = []ResultDataContent{"row", "graph"}
)

type Endpoints struct {
	// Extensions map[string]interface{} // not sure what this actually looks like.
	NodeBase          string `json:"node"`
	Relationship      string `json:"relationship"`
	NodeIndex         string `json:"node_index"`
	RelationshipTypes string `json:"relationship_types"` // return []string
	ExtensionsInfo    string `json:"extensions_info"`
	Batch             string
	RelationshipIndex string `json:"relationship_index"`
	Indexes           string // returns array of {property_keys: string[], label: string}
	Constraints       string
	NodeLabels        string `json:"node_labels"`
	Version           string `json:"neo4j_version"`
	Transaction       string `json:"transaction"`
}

// /node/:id response
// {
//   "extensions" : { },
//   "metadata" : {
//     "id" : 18003,
//     "labels" : [ "Group" ]
//   },
//   "paged_traverse" : "http://localhost:7474/db/data/node/18003/paged/traverse/{returnType}{?pageSize,leaseTime}",
//   "outgoing_relationships" : "http://localhost:7474/db/data/node/18003/relationships/out",
//   "outgoing_typed_relationships" : "http://localhost:7474/db/data/node/18003/relationships/out/{-list|&|types}",
//   "create_relationship" : "http://localhost:7474/db/data/node/18003/relationships",
//   "labels" : "http://localhost:7474/db/data/node/18003/labels",
//   "traverse" : "http://localhost:7474/db/data/node/18003/traverse/{returnType}",
//   "all_relationships" : "http://localhost:7474/db/data/node/18003/relationships/all",
//   "all_typed_relationships" : "http://localhost:7474/db/data/node/18003/relationships/all/{-list|&|types}",
//   "property" : "http://localhost:7474/db/data/node/18003/properties/{key}",
//   "self" : "http://localhost:7474/db/data/node/18003",
//   "incoming_relationships" : "http://localhost:7474/db/data/node/18003/relationships/in",
//   "properties" : "http://localhost:7474/db/data/node/18003/properties",
//   "incoming_typed_relationships" : "http://localhost:7474/db/data/node/18003/relationships/in/{-list|&|types}",
//   "data" : {
//     "name" : "INFORMATIONTECHNOLOGY1@EXTERNAL.LOCAL"
//   }
// }

func (e *Endpoints) ImmediateCommit() string {
	return e.Transaction + "/commit"
}

func (e *Endpoints) Commit(id string) string {
	return e.Transaction + "/" + id + "/commit"
}

func (e *Endpoints) Node(id string) string {
	return e.NodeBase + "/" + id
}

func (e *Endpoints) NodeProperties(id string) string {
	return e.NodeBase + "/" + id + "/properties"
}

type DB struct {
	addr      string
	creds     Creds
	c         *http.Client
	endpoints *Endpoints
}

type Creds struct {
	Username string
	Password string
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error fulfills the error interface
func (e Error) Error() string {
	return fmt.Sprintf("Neo4j error %s: %s", e.Code, e.Message)
}

type TransactionResponse struct {
	Results []TransactionResult
	Errors  []Error
}

func (resp TransactionResponse) Err() error {
	if len(resp.Errors) > 0 {
		return resp.Errors[0]
	}
	return nil
}

type TransactionResult struct {
	Columns []string
	Data    []Datum
}

type Datum struct {
	Row   []*json.RawMessage
	Meta  []interface{} // Not sure why this is always in inner array.
	Graph Graph
}

type Graph struct {
	Nodes         []Node
	Relationships []Relationship
}

type Node struct {
	Id         string                 `json:"id"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

type Relationship struct {
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	StartNode  string                 `json:"startNode"`
	EndNode    string                 `json:"endNode"`
	Properties map[string]interface{} `json:"properties"`
}

type Meta struct {
	Id      int64  `json:"id"`
	Type    string `json:"type"` // relationship or node
	Deleted bool   `json:"deleted"`
}

type Statement struct {
	Statement string `json:"statement"`
	// statement          string              `json:"statement"`
	Params             interface{}         `json:"parameters"`
	ResultDataContents []ResultDataContent `json:"resultDataContents,omitempty"`
	IncludeStats       bool                `json:"includeStats,omitempty"`
}

type Params map[string]interface{}

type Commit struct {
	Statements []Statement `json:"statements"`
}
