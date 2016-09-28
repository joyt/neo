package neo

type Cqlizer interface {
	ToCQL() (string, map[string]interface{}, error)
}
