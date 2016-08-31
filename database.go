package neo

import (
	"errors"
	"net/http"
)

func Connect(addr string, creds Creds) (*DB, error) {
	db := &DB{addr: addr, creds: creds, c: &http.Client{}}
	var endpoints Endpoints
	if err := db.get(addr+"/db/data", &endpoints); err != nil {
		return nil, err
	}
	db.endpoints = &endpoints
	return db, nil
}

func (db *DB) Commit(statement Statement) (TransactionResult, error) {
	tx := Transaction{Statements: []Statement{statement}}
	var res TransactionResponse
	if err := db.post(db.endpoints.ImmediateCommit(), tx, &res); err != nil {
		return TransactionResult{}, err
	}
	if err := res.Err(); err != nil {
		return TransactionResult{}, err
	}
	if len(res.Results) != 1 {
		return TransactionResult{}, errors.New("Wrong number of results")
	}
	return res.Results[0], nil
}
