package client

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

func newClientPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Client) error {
	const psqlInsert = `INSERT INTO cfg.client (full_name, nit, banner, logo_small, main_color, second_color, url_redirect, url_api) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.FullName,
		m.Nit,
		m.Banner,
		m.LogoSmall,
		m.MainColor,
		m.SecondColor,
		m.UrlRedirect,
		m.UrlApi,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *Client) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.client SET full_name = :full_name, nit = :nit, banner = :banner, logo_small = :logo_small, main_color = :main_color, second_color = :second_color, url_redirect = :url_redirect, url_api = :url_api, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM cfg.client WHERE id = :id `
	m := Client{ID: id}
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
func (s *psql) getByID(id int64) (*Client, error) {
	const psqlGetByID = `SELECT id , full_name, nit, banner, logo_small, main_color, second_color, url_redirect, url_api, created_at, updated_at FROM cfg.client WHERE id = $1 `
	mdl := Client{}
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
func (s *psql) getAll() ([]*Client, error) {
	var ms []*Client
	const psqlGetAll = ` SELECT id , full_name, nit, banner, logo_small, main_color, second_color, url_redirect, url_api, created_at, updated_at FROM cfg.client `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getByNit(nit string) (*Client, error) {
	const psqlGetByNit = `SELECT id , full_name, nit, banner, logo_small, main_color, second_color, url_redirect, url_api, created_at, updated_at FROM cfg.client WHERE nit = $1 `
	mdl := Client{}
	err := s.DB.Get(&mdl, psqlGetByNit, nit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
