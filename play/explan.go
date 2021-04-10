package play

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type User struct {
	Name string
	Type int
	Room int
}
type Message struct {
	User `json:"user"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
