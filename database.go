package neo

import (
	"bytes"
	"encoding/json"
	"strings"
	// "io/ioutil"
	"errors"
	// "log"
	"net/http"
	"reflect"
)

type ResultContentType string

const (
	ResultContentTypeRows ResultContentType = "rows"
)

type Database struct {
	addr  string
	creds Creds
	c     *http.Client
}

type Creds struct {
	Username string
	Password string
}

func Connect(addr string, creds Creds) (*Database, error) {
	return &Database{
		addr:  addr,
		creds: creds,
		c:     &http.Client{},
	}, nil
}

type TransactionResponse struct {
	Results []TransactionResult
	Errors  []interface{}
}

type TransactionResult struct {
	Columns []string
	Data    []struct {
		Row   []*json.RawMessage
		Meta  [1][]Meta
		Graph Graph
	}
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
	Statement          string                 `json:"statement"`
	Parameters         map[string]interface{} `json:"parameters"`
	ResultDataContents []string               `json:"resultDataContents,omitempty"`
	IncludeStats       bool                   `json:"includeStats,omitempty"`
}

type Transaction struct {
	Statements []Statement `json:"statements"`
}

func (r TransactionResult) UnmarshalRows(res interface{}) error {
	// TODO: This will panic if this is not a pointer to slice!
	arrType := reflect.TypeOf(res).Elem()
	arr := reflect.MakeSlice(arrType, len(r.Data), len(r.Data))
	mapping := map[string]string{}
	elemType := arrType.Elem()
	for f := 0; f < elemType.NumField(); f++ {
		field := elemType.Field(f)
		tag := field.Tag.Get("json")
		if len(tag) == 0 {
			continue
		}
		jsonName := strings.TrimSpace(strings.Split(tag, ",")[0])
		if len(jsonName) > 0 && jsonName != "-" {
			mapping[jsonName] = field.Name
		}
	}
	for i, datum := range r.Data {
		for j, c := range r.Columns {
			if _, ok := mapping[c]; !ok {
				continue
			}
			v := arr.Index(i).FieldByName(mapping[c]).Addr().Interface()
			if err := json.Unmarshal(*datum.Row[j], &v); err != nil {
				return err
			}
		}
	}
	reflect.Indirect(reflect.ValueOf(res)).Set(arr)
	return nil
}

func (db *Database) Commit(statement Statement) (TransactionResult, error) {
	tx := Transaction{Statements: []Statement{statement}}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(tx); err != nil {
		return TransactionResult{}, nil
	}
	resp, err := db.c.Post(db.addr+"/db/data/transaction/commit", "application/json", &b)
	if err != nil {
		return TransactionResult{}, err
	}
	defer resp.Body.Close()
	var res TransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return TransactionResult{}, err
	}
	if len(res.Results) != 1 {
		return TransactionResult{}, errors.New("Wrong number of results")
	}
	return res.Results[0], nil
}
