package work_validation

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

func newWorkValidationPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

func (s *psql) getPending() ([]*WorkValidation, error) {
	var ms []*WorkValidation
	const psqlGetAll = ` SELECT id , status, user_id, created_at, updated_at FROM wf.work_validation where status = 'pending'`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// Update actualiza un registro en la BD
func (s *psql) updateStatus(status string, userID string) error {
	date := time.Now()
	m := WorkValidation{
		Status:    status,
		UserId:    userID,
		UpdatedAt: date,
	}
	const psqlUpdate = `UPDATE wf.work_validation SET status = :status, updated_at = :updated_at WHERE user_id = :user_id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}
