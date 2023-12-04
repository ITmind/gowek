package repo

type Note struct {
	ID      uint
	Note    string
	User_id uint
}

type User struct {
	ID      uint
	Login   string
	Email   string
	Hash    string
	IsAdmin bool
}
