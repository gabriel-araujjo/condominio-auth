package postgres

import (
	"database/sql"
	"net/url"
	"testing"

	"github.com/gabriel-araujjo/condominio-auth/dao/internal/postgres/mock"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

func TestUserDaoPg(t *testing.T) {
	conf := mock.FakeDBConfig()

	db, err := sql.Open(conf.Dao.Driver, conf.Dao.URI)
	if err != nil {
		t.Fatalf("can't prepare db for test %e", err)
	}

	cleanDB(t, db)

	userDao, _, _, err := NewDao(conf)
	if err != nil {
		t.Fatalf("can't prepare dao for test %e", err)
	}
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
		}, {
			name: "InvalidCPF",
			user: &domain.User{
				Name:         "Fulano",
				CPF:          "12312312312",
				PasswordHash: "senha",
			},
			expectErr: true,
		}, {
			name: "WithoutCPF",
			user: &domain.User{
				Name:         "Fulano",
				PasswordHash: "senha",
			},
			expectErr: false,
		}, {
			name: "ManyWithoutCPF",
			user: &domain.User{
				Name:         "Fulano 2",
				PasswordHash: "senha",
			},
			expectErr: false,
		}, {
			name: "ManyWithoutFBID",
			user: &domain.User{
				Name:         "Fulano 3",
				PasswordHash: "senha",
			},
			expectErr: false,
		}, {
			name: "ExistentUser",
			user: &domain.User{
				ID:           1,
				Name:         "Fulano",
				CPF:          "04781539459",
				PasswordHash: "senha",
			},
			expectErr: true,
		}, {
			name: "WithoutPassword",
			user: &domain.User{
				Name: "Fulano",
				CPF:  "81763529509",
			},
			expectErr: false,
		}, {
			name: "DuplicatePrimaryEmail",
			user: &domain.User{
				Name: "Fulano",
				CPF:  "71308718144",
				Emails: []domain.Email{{
					Email:    "fulano@email.com",
					Verified: false,
				}},
			},
			expectErr: true,
		}, {
			name: "DuplicateSecondaryEmail",
			user: &domain.User{
				Name: "fulano",
				CPF:  "53828721559",
				Emails: []domain.Email{{
					Email:    "fulano@email2.com",
					Verified: false,
				}},
			},
			expectErr: true,
		}, {
			name: "DuplicatePrimaryPhone",
			user: &domain.User{
				Name: "Fulano",
				CPF:  "71308718144",
				Phones: []domain.Phone{{
					Phone:    "447588164927",
					Verified: false,
				}},
			},
			expectErr: true,
		}, {
			name: "DuplicateSecondaryPhone",
			user: &domain.User{
				Name: "Fulano",
				CPF:  "43730922300",
				Phones: []domain.Phone{{
					Phone:    "554365128899",
					Verified: false,
				}},
			},
			expectErr: true,
		}, {
			name:      "NilUser",
			user:      nil,
			expectErr: true,
		}}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				err := userDao.Create(tt.user)
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
				id, err := userDao.Auth(tt.credential, tt.password)
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
		user, err := userDao.Get(1)
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
				err := userDao.Delete(tt.id)
				if tt.expectErr {
					if err == nil {
						t.Errorf("test %q: should err be returned", tt.name)
					}
				} else {
					if err != nil {
						t.Errorf("test %q: should err be nil instead of %q", tt.name, err)
					}
					if u, _ := userDao.Get(tt.id); u != nil {
						t.Error("expecting user does not exist")
					}
				}
			})
		}
	})
}
