package mock

import (
	"fmt"
	"os"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

// InvalidDBConfig returns a Config with Dao.Driver field set to "Invalid"
func InvalidDBConfig() *config.Config {
	return &config.Config{
		Dao: config.Dao{Driver: "invalid", URI: "", VersionStrategy: "version1"},
	}
}

// FakeDBConfig returns a Config with Dao.Driver set to a real postgres database used for testing
func FakeDBConfig() *config.Config {
	return &config.Config{
		Dao: config.Dao{
			Driver: "postgres",
			URI: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				os.Getenv("POSTGRES_HOST"),
				os.Getenv("POSTGRES_PORT"),
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("POSTGRES_DB"),
				"disable"),
			VersionStrategy: "psql-versioning",
		},
		Clients: []*domain.Client{
			{Name: "Fake Client 1", PublicId: "1", Secret: "1"},
			{Name: "Fake Client 2", PublicId: "2", Secret: "2"},
		},
	}
}
