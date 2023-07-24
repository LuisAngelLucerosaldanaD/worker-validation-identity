package validation_request

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

func newValidationRequestPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *ValidationRequest) error {
	const psqlInsert = `INSERT INTO cfg.validation_request (client_id, max_num_validation, request_id, expired_at, user_identification, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.ClientId,
		m.MaxNumValidation,
		m.RequestId,
		m.ExpiredAt,
		m.UserIdentification,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *ValidationRequest) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.validation_request SET client_id = :client_id, max_num_validation = :max_num_validation, request_id = :request_id, expired_at = :expired_at, user_identification = :user_identification, status = :status, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM cfg.validation_request WHERE id = :id `
	m := ValidationRequest{ID: id}
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
func (s *psql) getByID(id int64) (*ValidationRequest, error) {
	const psqlGetByID = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_identification, status, created_at, updated_at FROM cfg.validation_request WHERE id = $1 `
	mdl := ValidationRequest{}
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
func (s *psql) getAll() ([]*ValidationRequest, error) {
	var ms []*ValidationRequest
	const psqlGetAll = ` SELECT id , client_id, max_num_validation, request_id, expired_at, user_identification, status, created_at, updated_at FROM cfg.validation_request `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByClientIDAndRequestID consulta un registro por su ID
func (s *psql) getByClientIDAndRequestID(clientIid int64, requestID string) (*ValidationRequest, error) {
	const psqlGetByClientID = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_identification, status, created_at, updated_at FROM cfg.validation_request WHERE client_id = %d and request_id = '%s' limit 1`
	mdl := ValidationRequest{}
	query := fmt.Sprintf(psqlGetByClientID, clientIid, requestID)
	err := s.DB.Get(&mdl, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) updateStatus(m *ValidationRequest) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.validation_request SET status = :status, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

func (s *psql) getPending() ([]*ValidationRequest, error) {
	var ms []*ValidationRequest
	const psqlGetAll = `SELECT id , client_id, max_num_validation, request_id, expired_at, user_identification, status, created_at, updated_at FROM cfg.validation_request where status = 'callback' or  status = 'refused';`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}
