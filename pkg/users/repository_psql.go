package users

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

func newUsersPsqlRepository(db *sqlx.DB, txID string) *psql {
	return &psql{
		DB:   db,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Users) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO auth.users (id ,type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at) VALUES (:id ,:type_document, :document_number, :expedition_date, :email, :first_name, :second_name, :second_surname, :age, :gender, :nationality, :civil_status, :first_surname, :birth_date, :country, :department, :city, :real_ip, :cellphone, :created_at, :updated_at) `
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
func (s *psql) update(m *Users) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE auth.users SET type_document = :type_document, document_number = :document_number, expedition_date = :expedition_date, email = :email, first_name = :first_name, second_name = :second_name, second_surname = :second_surname, age = :age, gender = :gender, nationality = :nationality, civil_status = :civil_status, first_surname = :first_surname, birth_date = :birth_date, country = :country, department = :department, city = :city, real_ip = :real_ip, cellphone = :cellphone, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM auth.users WHERE id = :id `
	m := Users{ID: id}
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
func (s *psql) getByID(id string) (*Users, error) {
	const psqlGetByID = `SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users WHERE id = $1 `
	mdl := Users{}
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
func (s *psql) getAll() ([]*Users, error) {
	var ms []*Users
	const psqlGetAll = ` SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByEmail consulta un registro por su ID
func (s *psql) getByEmail(email string) (*Users, error) {
	const psqlGetByEmail = `SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users WHERE email = $1 `
	mdl := Users{}
	err := s.DB.Get(&mdl, psqlGetByEmail, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

func (s *psql) getLasted(email string, limit, offset int) ([]*Users, error) {
	var ms []*Users
	const psqlGetAll = ` SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users where email <> $1 order by created_at desc limit $2 OFFSET $3`

	err := s.DB.Select(&ms, psqlGetAll, email, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getNotStarted() ([]*Users, error) {
	var ms []*Users
	const psqlGetAll = ` SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users where document_number = 0 or document_number is null`

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getNoUploadFile(fileType int) ([]*Users, error) {
	var ms []*Users
	const psqlGetAll = `SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at 
 from auth.users u where (select f.id from cfg.files f where f.user_id = u.id and f."type" = $1) is null ;`

	err := s.DB.Select(&ms, psqlGetAll, fileType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByEmail consulta un registro por su ID
func (s *psql) getByIdentityNumber(identityNumber int64) (*Users, error) {
	const psqlGetByEmail = `SELECT id , type_document, document_number, expedition_date, email, first_name, second_name, second_surname, age, gender, nationality, civil_status, first_surname, birth_date, country, department, city, real_ip, cellphone, created_at, updated_at FROM auth.users WHERE document_number = $1 `
	mdl := Users{}
	err := s.DB.Get(&mdl, psqlGetByEmail, identityNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}
