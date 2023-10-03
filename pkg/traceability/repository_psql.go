package traceability

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	TxID string
}

func newTraceabilityPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Traceability) error {
	const psqlInsert = `INSERT INTO trx.traceability (action, type, description, user_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.Action,
		m.Type,
		m.Description,
		m.UserId,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *Traceability) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE trx.traceability SET action = :action, type = :type, description = :description, user_id = :user_id, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM trx.traceability WHERE id = :id `
	m := Traceability{ID: id}
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
func (s *psql) getByID(id int64) (*Traceability, error) {
	const psqlGetByID = `SELECT id , action, type, description, user_id, created_at, updated_at FROM trx.traceability WHERE id = $1 `
	mdl := Traceability{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// getAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*Traceability, error) {
	var ms []*Traceability
	const psqlGetAll = ` SELECT id , action, type, description, user_id, created_at, updated_at FROM trx.traceability `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getByUserID(userId string) ([]*Traceability, error) {
	var ms []*Traceability
	const psqlGetAllByUserID = ` SELECT id , action, type, description, user_id, created_at, updated_at FROM trx.traceability where user_id = $1;`

	err := s.DB.Select(&ms, psqlGetAllByUserID, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// Delete elimina un registro de la BD
func (s *psql) deleteByUserID(userId string) error {
	const psqlDelete = `DELETE FROM trx.traceability WHERE user_id = :user_id `
	m := Traceability{UserId: userId}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
