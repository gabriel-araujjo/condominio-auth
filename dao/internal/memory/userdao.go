package memory

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gabriel-araujjo/condominio-auth/domain"
	"github.com/gabriel-araujjo/go-jsonpointer"
	"github.com/gabriel-araujjo/json-patcher"
)

type userTailor struct{}

func (userTailor) Add(obj interface{}, path string, value interface{}) error {

	pointer, err := jsonpointer.NewJSONPointerFromString(path)
	if err != nil {
		return err
	}

	if pointer.Depth() > 2 || pointer.Depth() < 1 {
		return fmt.Errorf("memory_userdao: invalid path %q", path)
	}

	user := obj.(*domain.User)

	//TODO: Use reflections

	switch pointer.Tokens()[0] {
	default:
		return fmt.Errorf("memory_userdao: invalid path %q", path)
	case "id":
		return errors.New("memory_userdao: can't edit user id")
	case "name":
		user.Name = value.(string)
	case "cpf":
		user.CPF = value.(string)
	case "fb_id":
		user.FbID = value.(string)
	case "avatar":
		url, err := url.Parse(value.(string))
		if err != nil {
			return err
		}
		user.Avatar = url
	case "phones":
		idx, err := getIndex(pointer, len(user.Phones))
		if err != nil {
			return err
		}
		// TODO: check whether phone already exists
		user.Phones = append(
			user.Phones[:idx],
			append(
				[]domain.Phone{
					{
						Phone:    value.(string),
						Verified: false,
					},
				},
				user.Phones[idx:]...,
			)...,
		)
	case "emails":
		idx, err := getIndex(pointer, len(user.Emails))
		if err != nil {
			return err
		}
		// TODO: check whether phone already exists
		user.Emails = append(
			user.Emails[:idx],
			append(
				[]domain.Email{
					{
						Email:    value.(string),
						Verified: false,
					},
				},
				user.Emails[idx:]...,
			)...,
		)
	case "password":
		user.PasswordHash = value.(string)
	}
	return nil
}

func (userTailor) Remove(obj interface{}, path string) error {
	pointer, err := jsonpointer.NewJSONPointerFromString(path)
	if err != nil {
		return err
	}

	if pointer.Depth() > 2 || pointer.Depth() < 1 {
		return fmt.Errorf("memory_userdao: invalid path %q", path)
	}

	user := obj.(*domain.User)

	switch pointer.Tokens()[0] {
	default:
		return fmt.Errorf("memory_userdao: invalid path %q", path)
	case "id":
		return errors.New("memory_userdao: can't remove id")
	case "name":
		user.Name = ""
	case "cpf":
		user.CPF = ""
	case "fb_id":
		user.FbID = ""
	case "avatar":
		user.Avatar = nil
	case "phones":
		idx, err := getIndex(pointer, len(user.Phones))
		if err != nil || idx == len(user.Phones) {
			if pointer.Depth() != 1 {
				return err
			}
			// clear all phones
			user.Phones = nil
			return nil
		}
		// TODO: check whether phone already exists
		user.Phones = append(user.Phones[:idx], user.Phones[idx+1:]...)
	case "emails":
		idx, err := getIndex(pointer, len(user.Emails))
		if err != nil || idx == len(user.Emails) {
			if pointer.Depth() != 1 {
				return err
			}
			// clear all phones
			user.Phones = nil
			return nil
		}
		// TODO: check whether phone already exists
		user.Phones = append(user.Phones[:idx], user.Phones[idx+1:]...)
	}
	return nil
}

func (userTailor) Move(obj interface{}, path string, from uint64, to uint64) error {
	panic("implement me")
}

func (userTailor) Replace(obj interface{}, path string, value interface{}) error {
	panic("implement me")
}

type userDaoMemory []*domain.User

func (d *userDaoMemory) Create(u *domain.User) error {
	if u == nil {
		return errors.New("memory_userdao: trying to create a nil user")
	}
	*d = append(*d, u)
	u.ID = int64(len(*d))
	return nil
}

func (d *userDaoMemory) Delete(id int64) error {
	if id <= 0 || id > int64(len(*d)) || (*d)[id-1] == nil {
		return errors.New("memory_userdao: no user was deleted")
	}
	(*d)[id-1] = nil
	return nil
}

func (d *userDaoMemory) Update(id int64, patch json_patcher.Patch) error {
	if id <= 0 || id > int64(len(*d)) || (*d)[id-1] == nil {
		return errors.New("memory_userdao: no user found")
	}
	return json_patcher.Mend(nil, patch, (*d)[id-1])
}

func (d *userDaoMemory) Get(id int64) (*domain.User, error) {
	if id <= 0 || id > int64(len(*d)) {
		return nil, errors.New("memory_userdao: no user found")
	}
	return (*d)[id-1], nil
}

func (d *userDaoMemory) Authenticate(credential string, password string) (int64, error) {
	var user *domain.User
	for i := range *d {
		if (*d)[i].CPF == credential {
			user = (*d)[i]
			break
		}

		for j := range (*d)[i].Emails {
			if (*d)[i].Emails[j].Email == credential {
				user = (*d)[i]
				break
			}
		}

		for j := range (*d)[i].Phones {
			if (*d)[i].Phones[j].Phone == credential {
				user = (*d)[i]
				break
			}
		}
	}

	if user == nil || user.PasswordHash != password {
		return -1, errors.New("memory_userdao: authentication failed")
	}

	return user.ID, nil
}

func (d *userDaoMemory) AuthorizeClient(userID int64, clientPublicID string, scope domain.Scope) error {
	// TODO
	return nil
}

func getIndex(pointer *jsonpointer.JSONPointer, maxValue int) (idx int, err error) {
	if pointer.Depth() < 2 {
		err = fmt.Errorf("memory_userdao: invalid path %v", pointer)
		return
	}
	idx, err = strconv.Atoi(pointer.Tokens()[1])
	if err != nil {
		if pointer.Tokens()[1] != "-" {
			return
		}
		idx = maxValue
		err = nil
	}
	if idx > maxValue {
		err = fmt.Errorf("memory_usage: out of range index (%v) max is (%v)", idx, maxValue)
	}
	return idx, err
}
