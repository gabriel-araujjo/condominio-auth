package memory

import (
	"errors"

	"github.com/gabriel-araujjo/condominio-auth/domain"
)

type permissionDaoMemory []string

func (p permissionDaoMemory) Create(permission string) error {
	id, _ := p.ScopeIntoPermissionIDs([]string{permission})
	if len(id) > 0 {
		return errors.New("duplicate permission")
	}
	p = append(p)
	return nil
}

func (p permissionDaoMemory) ScopeIntoPermissionIDs(scope domain.Scope) ([]int64, error) {
	permissions := make([]int64, 0, len(scope))
	for _, s := range scope {
		for i, name := range p {
			if s == name {
				permissions = append(permissions, int64(i))
			}
		}
	}
	return permissions, nil
}
