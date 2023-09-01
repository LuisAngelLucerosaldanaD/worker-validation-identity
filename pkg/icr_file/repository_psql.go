package icr_file

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

func newIcrFilePsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *IcrFile) error {
	const psqlInsert = `INSERT INTO icr.icr_file (user_id, text, x, y, w, h) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.UserId,
		m.Text,
		m.X,
		m.Y,
		m.W,
		m.H,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *IcrFile) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE icr.icr_file SET user_id = :user_id, text = :text, x = :x, y = :y, w = :w, h = :h, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM icr.icr_file WHERE id = :id `
	m := IcrFile{ID: id}
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
func (s *psql) getByID(id int64) (*IcrFile, error) {
	const psqlGetByID = `SELECT id , user_id, text, x, y, w, h, created_at, updated_at FROM icr.icr_file WHERE id = $1 `
	mdl := IcrFile{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*IcrFile, error) {
	var ms []*IcrFile
	const psqlGetAll = ` SELECT id , user_id, text, x, y, w, h, created_at, updated_at FROM icr.icr_file `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
