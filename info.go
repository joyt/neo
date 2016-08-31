package neo

// Labels returns all the distinct labels in this database.
func (db *DB) Labels() ([]string, error) {
	var res []string
	return res, db.get(db.endpoints.NodeLabels, &res)
}

// Relationships returns all the distinct relationship types in this database
func (db *DB) Relationships() ([]string, error) {
	var res []string
	return res, db.get(db.endpoints.RelationshipTypes, &res)
}

type IndexInfo struct {
	PropertyKeys []string `json:"property_keys"`
	Label        string
}

// Indexes returns all the indexes created on this database.
func (db *DB) Indexes() ([]IndexInfo, error) {
	var res []IndexInfo
	return res, db.get(db.endpoints.Indexes, &res)
}

// Version returns the version of this database.
func (db *DB) Version() (string, error) {
	var res string
	return res, db.get(db.endpoints.Version, &res)
}

// TODO: Constraints
