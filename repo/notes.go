package repo

import (
	"log/slog"
)

func GetAllNotesByUser(userID uint) ([]Note, error) {
	var entities []Note
	err := DB.Select(&entities, "SELECT * FROM notes WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
		return []Note{}, err
	}
	return entities, nil
}

func GetNote(userID uint, noteID uint) ([]Note, error) {
	var entities []Note
	err := DB.Select(&entities, "SELECT * FROM notes WHERE user_id = ? and id = ? ", userID, noteID)
	if err != nil {
		slog.Error(err.Error())
		return []Note{}, err
	}
	return entities, nil
}

func AddNote(note *Note) error {
	_, err := DB.NamedExec("INSERT INTO notes (Text, UserID) VALUES (:text, :userid)", &note)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
