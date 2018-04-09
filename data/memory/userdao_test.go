package memory

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/gabriel-araujjo/condominio-auth/domain"
)

func mustParse(s string) (u *url.URL) {
	u, _ = url.Parse(s)
	return
}

func TestUserTailor_Add(t *testing.T) {
	user := domain.User{
		ID:           1,
		Name:         "",
		CPF:          "",
		FbID:         "",
		Avatar:       nil,
		Phones:       nil,
		Emails:       nil,
		PasswordHash: "",
	}
	userType := reflect.TypeOf(&user)
	userValue := reflect.ValueOf(&user).Elem()

	tailor := userTailor{}

	tests := []struct {
		name        string
		jsonValue   interface{}
		expectValue interface{}
		field       string
		expectErr   bool
	}{{
		name:        "ID",
		jsonValue:   int64(2),
		expectValue: int64(2),
		expectErr:   true,
		field:       "ID",
	}, {
		name:        "Name",
		jsonValue:   "Carlos Silva",
		expectValue: "Carlos Silva",
		field:       "Name",
		expectErr:   false,
	}, {
		name:        "CPF",
		jsonValue:   "875.208.787-59",
		expectValue: "875.208.787-59",
		field:       "CPF",
		expectErr:   false,
	}, {
		name:        "FbID",
		field:       "FbID",
		jsonValue:   "12314345345",
		expectValue: "12314345345",
		expectErr:   false,
	}, {
		name:        "Avatar",
		field:       "Avatar",
		jsonValue:   "https://gravatar.com/asderer234dd34.jpg",
		expectValue: mustParse("https://gravatar.com/asderer234dd34.jpg"),
		expectErr:   false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			field, ok := userType.Elem().FieldByName(tt.field)
			if !ok {
				t.Fatalf("mal formatted test, field %q not found", tt.field)
			}
			jsonName := field.Tag.Get("json")
			err := tailor.Add(&user, fmt.Sprintf("/%s", jsonName), tt.jsonValue)

			if tt.expectErr {
				if err == nil {
					t.Errorf("test %q: should err be returned", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("test %q: should err be nil instead of %q", tt.name, err)
				}

				if !reflect.DeepEqual(userValue.FieldByIndex(field.Index).Interface(), tt.expectValue) {
					t.Errorf("test %q: %q field should be %v instead of %v",
						tt.name, field.Name, tt.expectValue, userValue.FieldByIndex(field.Index).Interface())
				}
			}
		})
	}
}

func TestUserTailor_Remove(t *testing.T) {
	user := domain.User{
		ID:     1,
		Name:   "Paulo Silva",
		CPF:    "646.112.228-10",
		FbID:   "213124",
		Avatar: mustParse("https://gravatar.com/avatar/234235345345323423423423"),
		Phones: []domain.Phone{
			{Phone: "(84) 9 5432-1111", Verified: true},
			{Phone: "(84) 9 5432-1132", Verified: false},
		},
		Emails: []domain.Email{
			{Email: "paulosilva@gmail.com", Verified: true},
		},
		PasswordHash: "",
	}

	userType := reflect.TypeOf(&user)
	userValue := reflect.ValueOf(&user).Elem()

	tailor := userTailor{}

	tests := []struct {
		name        string
		jsonValue   interface{}
		expectValue interface{}
		field       string
		expectErr   bool
	}{{
		name:        "ID",
		expectValue: int64(1),
		expectErr:   true,
		field:       "ID",
	}, {
		name:        "Name",
		expectValue: "",
		field:       "Name",
		expectErr:   false,
	}, {
		name:        "CPF",
		expectValue: "",
		field:       "CPF",
		expectErr:   false,
	}, {
		name:        "FbID",
		field:       "FbID",
		expectValue: "",
		expectErr:   false,
	}, {
		name:        "Avatar",
		field:       "Avatar",
		expectValue: (*url.URL)(nil),
		expectErr:   false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := userType.Elem().FieldByName(tt.field)
			if !ok {
				t.Fatalf("mal formatted test, field %q not found", tt.field)
			}
			jsonName := field.Tag.Get("json")
			err := tailor.Remove(&user, fmt.Sprintf("/%s", jsonName))

			if tt.expectErr {
				if err == nil {
					t.Errorf("test %q: should err be returned", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("test %q: should err be nil instead of %q", tt.name, err)
				}

				if !reflect.DeepEqual(userValue.FieldByIndex(field.Index).Interface(), tt.expectValue) {
					t.Errorf("test %q: %q field should be %v instead of %v",
						tt.name, field.Name, tt.expectValue, userValue.FieldByIndex(field.Index).Interface())
				}
			}
		})
	}
}

func TestUserDaoMemory(t *testing.T) {

	dao, _ := NewDao(nil)

	avatar, _ := url.Parse("https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d53")

	t.Run("Create", func(t *testing.T) {
		cases := []struct {
			name      string
			user      *domain.User
			expectErr bool
		}{{
			name: "ValidUser",
			user: &domain.User{
				Name:   "Fulano",
				CPF:    "61772443514",
				FbID:   "1111111111",
				Avatar: avatar,
				Phones: []domain.Phone{{
					Phone:    "447588164927",
					Verified: true,
				}, {
					Phone:    "554365128899",
					Verified: false,
				}},
				Emails: []domain.Email{{
					Email:    "fulano@email.com",
					Verified: true,
				}, {
					Email:    "fulano@email2.com",
					Verified: false,
				}},
				PasswordHash: "senha",
			},
			expectErr: false,
		},
			{
				name:      "NilUser",
				user:      nil,
				expectErr: true,
			}}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				err := dao.User.Create(tt.user)
				if tt.expectErr {
					if err == nil {
						t.Errorf("test %q: should err be returned", tt.name)
					}
				} else {
					if err != nil {
						t.Errorf("test %q: should err be nil instead of %q", tt.name, err)
					}

					if tt.user.ID <= 0 {
						t.Error("should an ID be set on user")
					}
				}
			})
		}
	})

	t.Run("Auth", func(t *testing.T) {
		cases := []struct {
			name       string
			credential string
			password   string
			expectID   int64
			expectErr  bool
		}{
			{
				name:       "cpf",
				credential: "61772443514",
				password:   "senha",
				expectID:   1,
				expectErr:  false,
			},
			{
				name:       "email",
				credential: "fulano@email.com",
				password:   "senha",
				expectID:   1,
				expectErr:  false,
			},
			{
				name:       "secondary_email",
				credential: "fulano@email2.com",
				password:   "senha",
				expectID:   1,
				expectErr:  false,
			},
			{
				name:       "phone",
				credential: "447588164927",
				password:   "senha",
				expectID:   1,
				expectErr:  false,
			},
			{
				name:       "secondary_phone",
				credential: "554365128899",
				password:   "senha",
				expectID:   1,
				expectErr:  false,
			},
			{
				name:       "wrong_password",
				credential: "554365128899",
				password:   "senhaerrada",
				expectErr:  true,
			},
			{
				name:       "wrong_credential",
				credential: "554365128890",
				password:   "senha",
				expectErr:  true,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				id, err := dao.User.Auth(tt.credential, tt.password)
				if tt.expectErr {
					if err == nil {
						t.Errorf("test %q: err should be returned", tt.name)
					}
				} else {
					if err != nil {
						t.Errorf("test %q: err should be nil instead of %q", tt.name, err)
					}

					if tt.expectID != id {
						t.Errorf("test %q: %d should be returned instead of %d", tt.name, tt.expectID, id)
					}
				}
			})
		}
	})

	t.Run("Get", func(t *testing.T) {
		user, err := dao.User.Get(1)
		if err != nil {
			t.Errorf("nonexpected error %q", err)
		}
		if user == nil {
			t.Error("expecting user != nil")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		cases := []struct {
			name      string
			id        int64
			expectErr bool
		}{{
			name:      "ValidID",
			id:        1,
			expectErr: false,
		}, {
			name:      "DoubleDeleting",
			id:        1,
			expectErr: true,
		}, {
			name:      "InvalidID",
			id:        20,
			expectErr: true,
		}, {
			name:      "NonexistentID",
			id:        400000,
			expectErr: true,
		}}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				err := dao.User.Delete(tt.id)
				if tt.expectErr {
					if err == nil {
						t.Errorf("test %q: should err be returned", tt.name)
					}
				} else {
					if err != nil {
						t.Errorf("test %q: should err be nil instead of %q", tt.name, err)
					}
					if u, _ := dao.User.Get(tt.id); u != nil {
						t.Error("expecting user does not exist")
					}
				}
			})
		}
	})
}
