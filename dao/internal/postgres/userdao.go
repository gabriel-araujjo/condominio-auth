package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gabriel-araujjo/condominio-auth/domain"
	jsonpointer "github.com/gabriel-araujjo/go-jsonpointer"
	patcher "github.com/gabriel-araujjo/json-patcher"
)

var userdaoStmts = map[string]string{
	"insert": `
			INSERT INTO "user"(name, cpf, fb_id, avatar, hash, phone,
				phone_verified, email, email_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "user".user_id
		`,
	"findByID": `
			SELECT u.user_id, u.name, u.cpf, u.fb_id, u.avatar, u.hash, u.phone,
				u.phone_verified, u.email, u.email_verified FROM "user" u
			WHERE u.user_id = $1 LIMIT 1
		`,
	"queryEmails": `
			SELECT e.email, e.verified FROM "user_email" e
			WHERE e.user_id = $1 ORDER BY e.email ASC
		`,
	"queryPhones": `
			SELECT p.phone, p.verified FROM "user_phone" p
			WHERE p.user_id = $1 ORDER BY p.phone ASC
		`,
	"mapCPFIntoID": `
			SELECT u.user_id FROM "user" u
		  	WHERE u.cpf = $1 LIMIT 1
		`,
	"mapFBIDIntoID": `
			SELECT u.user_id FROM "user" u
			WHERE u.fb_id = $1 LIMIT 1
		`,
	"mapEmailIntoID": `
			SELECT l.user_id FROM "email_lookup" l
			WHERE l.email = $1 LIMIT 1
		`,
	"mapPhoneIntoID": `
			SELECT l.user_id FROM "phone_lookup" l
			WHERE l.phone = $1 LIMIT 1
		`,
	"updateByID": `
			UPDATE "user"
			SET
				name = $2,
				cpf = $3,
				fb_id = $4,
				avatar = $5,
				hash = $6,
				phone= $7,
				phone_verified = $8,
				email = $9,
				email_verified = $10
			WHERE user_id = $1
		`,
	"addEmail": `
			INSERT INTO "user_email"(user_id, email, verified)
			VALUES ($1, $2, $3) RETURNING user_id, email
		`,
	"addPhone": `
			INSERT INTO "user_phone"(user_id, phone, verified)
			VALUES ($1, $2, $3) RETURNING user_id, phone
		`,
	"removeEmail": `
			DELETE FROM "user_email" WHERE user_id = $1 AND email = $2 RETURNING email
		`,
	"removePhone": `
			DELETE FROM "user_phone" WHERE user_id = $1 AND phone = $2 RETURNING phone
		`,
	"deleteByID": `
			DELETE FROM "user" WHERE user_id = $1 RETURNING user_id
		`,
	"deleteAllEmails": `
			DELETE FROM "user_email" WHERE user_id = $1 RETURNING email
		`,
	"deleteAllPhones": `
			DELETE FROM "user_phone" WHERE user_id = $1 RETURNING phone
		`,
	"auth": `
			SELECT u.user_id
            FROM "user" u
				LEFT JOIN "user_email" e ON u.user_id = e.user_id
				LEFT JOIN "user_phone" P ON p.user_id = p.user_id
            WHERE
				( u.email = $1 OR
                  u.cpf = $2 OR
                  u.phone = $1 OR
                  e.email = $1 OR
                  p.phone = $1
				)
			AND u.hash = crypt($3, u.hash);
	`,
}

type updateStmt struct {
	updates []string
	args    []interface{}
}

type patch struct {
	user          domain.User
	userPatch     updateStmt
	addedEmails   []string
	deletedEmails []string
	addedPhones   []string
	deletedPhones []string
}

type userTailor struct{}

func (userTailor) Add(obj interface{}, path string, value interface{}) error {

	pointer, err := jsonpointer.NewJSONPointerFromString(path)
	if err != nil {
		return err
	}

	if pointer.Depth() > 2 || pointer.Depth() < 1 {
		return fmt.Errorf("postgres_userdao: invalid path %q", path)
	}

	updateStmt := obj.(*updateStmt)

	//TODO: Use reflections

	switch pointer.Tokens()[0] {
	default:
		return fmt.Errorf("postgres_userdao: invalid path %q", path)
	case "id":
		return errors.New("postgres_userdao: can't edit user id")
	case "name":
		updateStmt.updates = append(updateStmt.updates, "name = $")
		updateStmt.args = append(updateStmt.args, value)
	case "cpf":
		updateStmt.updates = append(updateStmt.updates, "cpf = $")
		updateStmt.args = append(updateStmt.args, value)
	case "fb_id":
		updateStmt.updates = append(updateStmt.updates, "fb_id = $")
		updateStmt.args = append(updateStmt.args, value)
	case "avatar":
		updateStmt.updates = append(updateStmt.updates, "avatar = $")
		updateStmt.args = append(updateStmt.args, value)
	case "password":
		updateStmt.updates = append(updateStmt.updates, "hash = $")
		updateStmt.args = append(updateStmt.args, value)
	}
	return nil
}

func (userTailor) Remove(obj interface{}, path string) error {
	return errors.New("implement me")
}

func (userTailor) Move(obj interface{}, path string, from uint64, to uint64) error {
	return errors.New("implement me")
}

func (userTailor) Replace(obj interface{}, path string, value interface{}) error {
	return errors.New("implement me")
}

type userDaoPG struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func (d *userDaoPG) lazyPrepare() {
	if d.stmts == nil {
		prepared := map[string]*sql.Stmt{}
		for k, v := range userdaoStmts {
			stmt, e := d.db.Prepare(v)
			if e != nil {
				panic(fmt.Sprintf("can't prepare user dao statement %q. Err: %q", v, e))
			}
			prepared[k] = stmt
		}
		d.stmts = prepared
	}
}

func (d *userDaoPG) Create(u *domain.User) error {
	if u == nil {
		return errors.New("postgres_userdao: Trying to create nil user")
	}
	if u.ID != 0 {
		return errors.New("postgres_userdao: Trying to create already created user")
	}

	d.lazyPrepare()

	name := normalizeString(u.Name)
	cpf := normalizeString(u.CPF)
	fbID := normalizeString(u.FbID)
	password := u.PasswordHash
	primaryEmail, verifiedPrimaryEmail := inflateEmail(u.PrimaryEmail())
	primaryPhone, verifiedPhone := inflatePhone(u.PrimaryPhone())
	avatar := safeString(u.Avatar)

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	row := d.stmts["insert"].QueryRow(name, cpf, fbID, avatar, password,
		primaryPhone, verifiedPhone,
		primaryEmail, verifiedPrimaryEmail)

	if err := row.Scan(&u.ID); err != nil {
		tx.Rollback()
		return err
	}

	if len(u.Emails) > 1 {
		for _, email := range u.Emails[1:] {
			_, err := d.stmts["addEmail"].Query(u.ID, email.Email, email.Verified)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	if len(u.Phones) > 1 {
		for _, phone := range u.Phones[1:] {
			_, err := d.stmts["addPhone"].Query(u.ID, phone.Phone, phone.Verified)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (d *userDaoPG) Delete(id int64) error {
	if id < 1 {
		return errors.New("postgres_userdao: Trying to delete user with invalid ID")
	}

	d.lazyPrepare()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	r, err := d.stmts["deleteByID"].Exec(id)

	if err != nil {
		tx.Rollback()
		return err
	}
	if rows, err := r.RowsAffected(); err != nil || rows != 1 {
		if err != nil {
			return err
		}
		return errors.New("postgres_userdao: no user was deleted")
	}
	tx.Commit()
	return err
}

func (d *userDaoPG) Update(id int64, patch patcher.Patch) error {
	d.lazyPrepare()

	return errors.New("unimplemented method")
}

func (d *userDaoPG) Get(id int64) (*domain.User, error) {
	u := domain.User{}
	p := domain.Phone{}
	e := domain.Email{}

	d.lazyPrepare()

	tx, err := d.db.Begin()

	if err != nil {
		return nil, err
	}

	{
		var avatarString string
		err = d.stmts["findByID"].QueryRow(id).Scan(&u.ID, &u.Name, &u.CPF, &u.FbID, &avatarString, &u.PasswordHash,
			&p.Phone, &p.Verified, &e.Email, &e.Verified)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
		u.Avatar, err = url.Parse(avatarString)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	emails := []domain.Email{e}
	phones := []domain.Phone{p}

	rows, err := d.stmts["queryEmails"].Query(id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		email := domain.Email{}
		rows.Scan(&email.Email, &email.Verified)
		emails = append(emails, email)
	}

	rows, err = d.stmts["queryPhones"].Query(id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		phone := domain.Phone{}
		rows.Scan(&phone.Phone, &phone.Verified)
		phones = append(phones, phone)
	}

	u.Emails = emails
	u.Phones = phones

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *userDaoPG) Authenticate(credential string, password string) (int64, error) {
	d.lazyPrepare()
	var id int64
	cpf, _ := strconv.ParseInt(credential, 10, 64)

	err := d.stmts["auth"].QueryRow(credential, cpf, password).Scan(&id)
	return id, err
}

func (d *userDaoPG) Authorize(*domain.ClientAuthorizationRequest) error {
	//TODO
	return nil
}

func safeString(url *url.URL) *string {
	if url == nil {
		return nil
	}

	s := url.String()
	return &s
}

func normalizeString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func inflateEmail(email *domain.Email) (*string, bool) {
	if email != nil {
		return &email.Email, email.Verified
	}
	return nil, false
}

func inflatePhone(phone *domain.Phone) (*string, bool) {
	if phone != nil {
		return &phone.Phone, phone.Verified
	}
	return nil, false
}
