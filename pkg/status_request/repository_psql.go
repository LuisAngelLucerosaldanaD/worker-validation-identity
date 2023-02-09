package status_request

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	TxID string
}

func newStatusRequestPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Update actualiza un registro en la BD
func (s *psql) updateStatus(status string, description string, userID string) error {
	date := time.Now()
	m := StatusRequest{
		Status:      status,
		UserId:      userID,
		Description: description,
		UpdatedAt:   date,
	}
	const psqlUpdate = `UPDATE wf.status_request SET status = :status, description = :description, updated_at = :updated_at WHERE user_id = :user_id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
