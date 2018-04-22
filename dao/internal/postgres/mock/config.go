package mock

import (
	"fmt"
	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gabriel-araujjo/condominio-auth/domain"
	"os"
)

func InvalidDBConfig() *config.Config {
	return &config.Config{
		Dao: struct {
			Driver          string
			DNS             string
			VersionStrategy string
		}{Driver: "invalid", DNS: "", VersionStrategy: "version1"},
	}
}

func FakeDBConfig() *config.Config {
	return &config.Config{
		Dao: struct {
			Driver          string
			DNS             string
			VersionStrategy string
		}{
			Driver: "postgres",
			DNS: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
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
