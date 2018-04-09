package data

import (
	"testing"
	"reflect"
	"github.com/gabriel-araujjo/condominio-auth/mock/io"
	"github.com/golang/mock/gomock"
	"github.com/gabriel-araujjo/condominio-auth/mock/data"
)

func TestDao(t *testing.T) {

	t.Run("client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_data.NewMockClientDao(ctrl)

		tests := []struct {
			name string
			dao  Dao
			want ClientDao
		}{{
			name: "not nil client dao",
			dao:  Dao{Client:client},
			want: client,
		}, {
			name: "nil client dao",
			dao:  Dao{Client: nil},
			want: nil,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if !reflect.DeepEqual(tt.dao.Client, tt.want) {
					t.Errorf("daoPG: user got=%v, want=%v", tt.dao.Client, tt.want)
				}
			})
		}
	})

	t.Run("user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userDao := mock_data.NewMockUserDao(ctrl)

		tests := []struct {
			name string
			dao  Dao
			want UserDao
		}{{
			name: "not nil user dao",
			dao:  Dao{User: userDao},
			want: userDao,
		}, {
			name: "nil user dao",
			dao:  Dao{User: nil},
			want: nil,
		}}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if !reflect.DeepEqual(tt.dao.User, tt.want) {
					t.Errorf("daoPG: user got=%v, want=%v", tt.dao.User, tt.want)
				}
			})
		}
	})

	t.Run("close", func(t *testing.T) {
		t.Run("not nil", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			closer := mock_io.NewMockCloser(ctrl)
			closer.EXPECT().Close()

			dao := Dao{closer: closer}

			dao.Close()
		})

		t.Run("nil", func(t *testing.T) {
			dao := Dao{closer: nil}

			if err := dao.Close(); err != nil {
				t.Errorf("should not return err")
			}
		})

	})
}

func TestNewDao(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_data.NewMockClientDao(ctrl)
	user := mock_data.NewMockUserDao(ctrl)
	closer := mock_io.NewMockCloser(ctrl)

	dao := NewDao(user, client, closer)

	if !reflect.DeepEqual(dao.User, user) {
		t.Errorf("expect user is %q instead of %q", user, dao.User)
	}

	if !reflect.DeepEqual(dao.Client, client) {
		t.Errorf("expect user is %q instead of %q", client, dao.Client)
	}

}