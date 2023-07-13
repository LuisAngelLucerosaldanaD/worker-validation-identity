package onboarding

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

func newOnboardingPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Onboarding) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.onboarding (id ,client_id, request_id, user_id, status, created_at, updated_at) VALUES (:id ,:client_id, :request_id, :user_id, :status,:created_at, :updated_at) `
	rs, err := s.DB.NamedExec(psqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *Onboarding) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.onboarding SET client_id = :client_id, request_id = :request_id, user_id = :user_id, status = :status, updated_at = :updated_at WHERE id = :id `
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
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM auth.onboarding WHERE id = :id `
	m := Onboarding{ID: id}
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
func (s *psql) getByID(id string) (*Onboarding, error) {
	const psqlGetByID = `SELECT id , client_id, request_id, user_id, status, created_at, updated_at FROM auth.onboarding WHERE id = $1 `
	mdl := Onboarding{}
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
func (s *psql) getAll() ([]*Onboarding, error) {
	var ms []*Onboarding
	const psqlGetAll = ` SELECT id , client_id, request_id, user_id, status, created_at, updated_at FROM auth.onboarding `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getPending() ([]*Onboarding, error) {
	var ms []*Onboarding
	const psqlGetAll = ` SELECT id , client_id, request_id, user_id, status, created_at, updated_at FROM auth.onboarding where status = 'pending';`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
