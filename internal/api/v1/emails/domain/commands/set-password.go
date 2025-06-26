package commands

type SetPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
