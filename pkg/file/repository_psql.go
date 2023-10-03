package file

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	TxID string
}

func newFilePsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *File) error {
	const psqlInsert = `INSERT INTO cfg.file (path, name, type, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Path,
		m.Name,
		m.Type,
		m.UserId,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *File) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.file SET path = :path, name = :name, type = :type, user_id = :user_id, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id int64) error {
	const psqlDelete = `DELETE FROM cfg.file WHERE id = :id `
	m := File{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id int64) (*File, error) {
	const psqlGetByID = `SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.file WHERE id = $1 `
	mdl := File{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*File, error) {
	var ms []*File
	const psqlGetAll = ` SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.file `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByUserID(userID string) ([]*File, error) {
	var ms []*File
	const psqlGetAll = ` SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.file where user_id = $1 `

	err := s.DB.Select(&ms, psqlGetAll, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// Delete elimina un registro de la BD
func (s *psql) deleteByUserId(userId string) error {
	const psqlDelete = `DELETE FROM cfg.file WHERE user_id = :user_id `
	m := File{UserId: userId}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByTypeAndUserID(typeFile int, userID string) (*File, error) {
	const psqlGetByID = `SELECT id , path, name, type, user_id, created_at, updated_at FROM cfg.file WHERE type = $1 and user_id = $2`
	mdl := File{}
	err := s.DB.Get(&mdl, psqlGetByID, typeFile, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
