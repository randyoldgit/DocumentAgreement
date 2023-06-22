package entities

type Message struct {
	Body string `json:"body"`
}

type UserAuth struct {
	UserID   int    `json:"userId"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
