package postgres

import (
	"database/sql"
	"github.com/gabriel-araujjo/condominio-auth/domain"
	"errors"
)

type clientDaoPG struct {
	db      *sql.DB
	getStmt *sql.Stmt
}

func (d *clientDaoPG) lazyPrepare()  {
	if d.getStmt == nil {
		var e error
		d.getStmt, e = d.db.Prepare(`
			SELECT c.client_id, c.name, c.public_id, c.secret
			FROM "client" c
			WHERE c.public_id = $1
			LIMIT 1
		`)
		if e != nil {
			panic("can't prepare client dao ")
		}
	}
}

func (d *clientDaoPG) Create(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoPG) Delete(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoPG) Update(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoPG) Get(publicId string) (*domain.Client, error) {
	d.lazyPrepare()

	client := &domain.Client{}
	row := d.getStmt.QueryRow(publicId)

	err := row.Scan(&client.ID, &client.Name, &client.PublicId, &client.Secret)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *clientDaoPG) Auth(publicID string, secret string) (string, error) {
	client, err := d.Get(publicID)
	if err != nil || client.Secret != secret {
		return "", errors.New("unauthorized client")
	}
	return client.PublicId, nil
}
