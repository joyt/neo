package neo

import (
	"crypto/tls"
	"errors"
	"net/http"
)

func Connect(addr string, creds Creds) (*DB, error) {
	db := &DB{addr: addr, creds: creds, c: &http.Transport{
		TLSClientConfig: tls.Config{InsecureSkipVerify: true},
	}}
	if err := db.get(addr+"/db/data/", &db.endpoints); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Commit(statements ...Statement) ([]TransactionResult, error) {
	tx := Commit{Statements: statements}
	var res TransactionResponse
	if err := db.post(db.endpoints.ImmediateCommit(), tx, &res); err != nil {
		return nil, err
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	if len(res.Results) != len(statements) {
		return nil, errors.New("Wrong number of results")
	}
	return res.Results, nil
}

type Transaction struct {
	CommitURL string
	// Expires   time.Time
}

// {
//   "commit": "http://localhost:7474/db/data/transaction/2484/commit",
//   "results": [],
//   "transaction": {
//     "expires": "Tue, 13 Sep 2016 13:22:13 +0000"
//   },
//   "errors": []
// }

type TransactionResp struct {
	CommitURL   string
	Results     interface{}
	Transaction struct{ Expires string }
	Errors      interface{}
}

func (db *DB) OpenTransaction() (*Transaction, error) {
	// db.endpoints.TransactionBase
	// Start transaction.
	return &Transaction{}, nil
}

// func (tx *Transaction) Commit(statement ...Statement) (TransactionResult, error) {
// 	return nil
// }
