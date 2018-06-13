package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gabriel-araujjo/condominio-auth/domain"
)

var permissionStmts = map[string]string{
	"insert": `
			INSERT INTO "scope"(name) VALUES ($1) RETURNING "scope".user_id
	`,
	"mapIDs": `
			SELECT s.scope_id FROM "scope" s 
			WHERE s.name = ANY ($1)
	`,
}

// PermissionDao manage queries related to permissions
type pgPermissionDao struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func (d *pgPermissionDao) Create(name string) error {
	_, err := d.stmts["insert"].Exec(name)
	return err
}

func (d *pgPermissionDao) ScopeIntoPermissionIDs(scope domain.Scope) ([]int64, error) {
	rows, err := d.stmts["mapIDs"].Query(scope)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(scope))
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func newPGPermissionDao(db *sql.DB) *pgPermissionDao {
	prepared := map[string]*sql.Stmt{}
	for k, v := range permissionStmts {
		stmt, e := db.Prepare(v)
		if e != nil {
			panic(fmt.Sprintf("can't prepare permission dao statement %q. Err: %q", v, e))
		}
		prepared[k] = stmt
	}
	return &pgPermissionDao{db, prepared}
}
