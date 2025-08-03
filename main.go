package main

type User struct {
	id   int
	name string
	//conn  conection type
}

type Channel struct {
	activeUsers []User
	messages    []string
}
type Server struct {
	channels []Channel
}

func (c *Channel) Broadcast(message string) {
	c.messages = append(c.messages, message)
}
