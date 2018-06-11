package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gabriel-araujjo/base62"

	"math/rand"

	"github.com/gabriel-araujjo/condominio-auth/domain"
)

var count = rand.Int31()

var clientDaoStmts = map[string]string{
	"get": `
			SELECT c.client_id, c.name, c.secret
			FROM "client" c
			WHERE c.client_id = $1
			LIMIT 1
		`,
	"insert": `
			INSERT INTO "client"(client_id, name, secret)
			VALUES ($1, $2, $3)
			RETURNING "client".client_id
		`,
	"permissionsByUser": `
			SELECT s.scope_id, s.name, s.description 
			FROM "authorization" a LEFT JOIN "scope" s ON a.scope_id = s.scope_id
			WHERE a.client_id = $1 AND a.user_id = $2
	`,
}

func newClientID() int64 {
	var (
		timestamp = uint64(time.Now().Unix() << 32)
		pid       = uint64((os.Getpid() & 0xffff) << 16)
		c         = uint64(atomic.AddInt32(&count, 1) & 0xffff)
	)
	return int64(timestamp | pid | c)
}

type clientDaoPG struct {
	db          *sql.DB
	getStmt     *sql.Stmt
	scopeByUser *sql.Stmt
	stmts       map[string]*sql.Stmt
}

func (d *clientDaoPG) lazyPrepare() {
	if d.stmts == nil {
		prepared := map[string]*sql.Stmt{}
		for k, v := range clientDaoStmts {
			stmt, e := d.db.Prepare(v)
			if e != nil {
				panic(fmt.Sprintf("can't prepare client dao statement %q. Err: %q", v, e))
			}
			prepared[k] = stmt
		}
	}
}

func (d *clientDaoPG) Create(c *domain.Client) error {
	if c.ID == 0 {
		return errors.New("pg: cliente already created")
	}
	d.lazyPrepare()
	clientID := newClientID()

	row := d.stmts["insert"].QueryRow(clientID, c.Name, c.Secret)
	if err := row.Scan(&c.ID); err != nil {
		return err
	}

	c.PublicID = base62.FormatUint(uint64(clientID))

	panic("implement me")
}

func (d *clientDaoPG) Delete(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoPG) Update(u *domain.Client) error {
	panic("implement me")
}

func (d *clientDaoPG) Get(publicID string) (*domain.Client, error) {
	d.lazyPrepare()

	clientID, err := base62.ParseUint(publicID)
	if err != nil {
		return nil, fmt.Errorf("pg: cannot find client with publicId = %q", publicID)
	}

	client := &domain.Client{}
	row := d.getStmt.QueryRow(clientID)

	err = row.Scan(&client.ID, &client.Name, &client.Secret)
	if err != nil {
		return nil, err
	}

	client.PublicID = publicID
	return client, nil
}

func (d *clientDaoPG) Auth(publicID string, secret string) (string, error) {
	d.lazyPrepare()
	client, err := d.Get(publicID)
	if err != nil || client.Secret != secret {
		return "", errors.New("unauthorized client")
	}
	return client.PublicID, nil
}

func (d *clientDaoPG) GetAuthorizedScopesByUser(publicID string, userID int64) domain.Scope {
	d.lazyPrepare()

	clientID, err := base62.ParseUint(publicID)
	if err != nil {
		return nil
	}

	rows, err := d.stmts["permissionsByUser"].Query(clientID, userID)
	if err != nil {
		return nil
	}

	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var p string
		rows.Scan(&p)
		permissions = append(permissions, p)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}

	return permissions
}
