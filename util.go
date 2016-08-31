package neo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

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
			// Try using the field name with first letter lowercased.
			// Note: Name must not start with multi-byte rune
			tag = strings.ToLower(field.Name[:0]) + field.Name[1:]
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

func (db *DB) newReq(method, endpoint string, payload interface{}) (*http.Request, error) {
	var req *http.Request
	var err error
	if payload != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, endpoint, &buf)
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}
	if err != nil {
		return nil, err
	}
	if len(db.creds.Username) > 0 {
		req.SetBasicAuth(db.creds.Username, db.creds.Password)
	}
	req.Header.Set("X-Stream", "true")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (db *DB) get(endpoint string, res interface{}) error {
	req, err := db.newReq("GET", endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := db.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(res)
}

func (db *DB) post(endpoint string, payload, res interface{}) error {
	req, err := db.newReq("POST", endpoint, payload)
	if err != nil {
		return err
	}
	resp, err := db.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(res)
}
