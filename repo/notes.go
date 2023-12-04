package repo

import (
	"log/slog"
)

func GetAllNotesByUser(userID uint) ([]Note, error) {
	var entities []Note
	err := db.Select(&entities, "SELECT * FROM notes WHERE user_id = ?", userID)
	if err != nil {
		slog.Error(err.Error())
		return []Note{}, err
	}
	return entities, nil
}

func GetNote(userID uint, noteID uint) ([]Note, error) {
	var entities []Note
	err := db.Select(&entities, "SELECT * FROM notes WHERE user_id = ? and id = ? ", userID, noteID)
	if err != nil {
		slog.Error(err.Error())
		return []Note{}, err
	}
	return entities, nil
}

func AddNote(note *Note) error {
	_, err := db.NamedExec("INSERT INTO notes (note, user_id) VALUES (:note, :user_id)", &note)
	if err != nil {
		slog.Error(err.Error())
	}

	return err
}
