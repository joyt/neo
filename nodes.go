package neo

func (db *DB) NodeProperties(id string) (map[string]interface{}, error) {
	var res map[string]interface{}
	return res, db.get(db.endpoints.NodeProperties(id), &res)
}

// {
//   "extensions" : { },
//   "metadata" : {
//     "id" : 1172,
//     "labels" : [ "User" ]
//   },
//   "paged_traverse" : "http://localhost:7474/db/data/node/1172/paged/traverse/{returnType}{?pageSize,leaseTime}",
//   "outgoing_relationships" : "http://localhost:7474/db/data/node/1172/relationships/out",
//   "outgoing_typed_relationships" : "http://localhost:7474/db/data/node/1172/relationships/out/{-list|&|types}",
//   "create_relationship" : "http://localhost:7474/db/data/node/1172/relationships",
//   "labels" : "http://localhost:7474/db/data/node/1172/labels",
//   "traverse" : "http://localhost:7474/db/data/node/1172/traverse/{returnType}",
//   "all_relationships" : "http://localhost:7474/db/data/node/1172/relationships/all",
//   "all_typed_relationships" : "http://localhost:7474/db/data/node/1172/relationships/all/{-list|&|types}",
//   "property" : "http://localhost:7474/db/data/node/1172/properties/{key}",
//   "self" : "http://localhost:7474/db/data/node/1172",
//   "incoming_relationships" : "http://localhost:7474/db/data/node/1172/relationships/in",
//   "properties" : "http://localhost:7474/db/data/node/1172/properties",
//   "incoming_typed_relationships" : "http://localhost:7474/db/data/node/1172/relationships/in/{-list|&|types}",
//   "data" : {
