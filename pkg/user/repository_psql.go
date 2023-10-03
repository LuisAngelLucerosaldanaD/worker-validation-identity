package user

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

func newUserPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *User) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.user (id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at) VALUES (:id, :nickname, :email, :password, :first_name, :second_name, :first_surname, :second_surname, :age, :type_document, :document_number, :cellphone, :gender, :nationality, :country, :department, :city, :real_ip, :status_id, :failed_attempts, :block_date, :disabled_date, :last_login, :last_change_password, :birth_date, :verified_code, :is_deleted,
:deleted_at, :created_at, :updated_at) `
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
func (s *psql) update(m *User) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.user SET nickname = :nickname, email = :email, password = :password, first_name = :first_name, second_name = :second_name, first_surname = :first_surname, second_surname = :second_surname, age = :age, type_document = :type_document, document_number = :document_number, cellphone = :cellphone, gender = :gender, nationality = :nationality, country = :country, 
                     department = :department, city = :city, real_ip = :real_ip, status_id = :status_id, failed_attempts = :failed_attempts, block_date = :block_date, disabled_date = :disabled_date, last_login = :last_login, last_change_password = :last_change_password, birth_date = :birth_date, verified_code = :verified_code, is_deleted = :is_deleted, deleted_at = :deleted_at, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.user WHERE id = :id `
	m := User{ID: id}
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
func (s *psql) getByID(id string) (*User, error) {
	const psqlGetByID = `SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user WHERE id = $1 `
	mdl := User{}
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
func (s *psql) getAll() ([]*User, error) {
	var ms []*User
	const psqlGetAll = ` SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByEmail consulta un registro por su ID
func (s *psql) getByEmail(email string) (*User, error) {
	const psqlGetByEmail = `SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user WHERE email = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getLasted(email string, limit, offset int) ([]*User, error) {
	var ms []*User
	const psqlGetAll = ` SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user where email <> $1 order by created_at desc limit $2 OFFSET $3`

	err := s.DB.Select(&ms, psqlGetAll, email, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getNotStarted() ([]*User, error) {
	var ms []*User
	const psqlGetAll = ` SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user where document_number = ''`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getNoUploadFile(fileType int) ([]*User, error) {
	var ms []*User
	const psqlGetAll = `SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at 
 from auth.user u where (select f.id from cfg.file f where f.user_id = u.id and f."type" = $1) is null ;`

	err := s.DB.Select(&ms, psqlGetAll, fileType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByEmail consulta un registro por su ID
func (s *psql) getByIdentityNumber(identityNumber string) (*User, error) {
	const psqlGetByEmail = `SELECT id, nickname, email, password, first_name, second_name, first_surname, second_surname, age, type_document, document_number, cellphone, gender, nationality, country, department, city, real_ip, status_id, failed_attempts, block_date, disabled_date, last_login, last_change_password, birth_date, verified_code, is_deleted,
deleted_at, created_at, updated_at FROM auth.user WHERE document_number = $1 `
	mdl := User{}
	err := s.DB.Get(&mdl, psqlGetByEmail, identityNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
