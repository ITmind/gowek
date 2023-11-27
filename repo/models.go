package repo

type Note struct {
	ID     uint
	Text   string
	UserID uint
}

type User struct {
	ID      uint
	Login   string
	Email   string
	Hash    string
	IsAdmin bool
}
