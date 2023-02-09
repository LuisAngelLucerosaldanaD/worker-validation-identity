package files

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	TxID string
}

func newFilesPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

func (s *psql) getByUserID(userID string) ([]*Files, error) {
	var ms []*Files
	const psqlGetAll = ` SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.files where user_id = $1 `

	err := s.DB.Select(&ms, psqlGetAll, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByTypeAndUserID(typeFile int, userID string) (*Files, error) {
	const psqlGetByID = `SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.files WHERE type = $1 and user_id = $2`
	mdl := Files{}
	err := s.DB.Get(&mdl, psqlGetByID, typeFile, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
