package neo

type Query struct {
	Match  string
	With   string
	Where  []string
	Return []string
}

type Match struct {
	PathVar string
}
