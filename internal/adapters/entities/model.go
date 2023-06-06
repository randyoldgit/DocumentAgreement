package entities

type Message struct {
	Body string `json:"body"`
}

type UserAuth struct {
	UserName string `json:"login"`
	Password string `json:"password"`
}
