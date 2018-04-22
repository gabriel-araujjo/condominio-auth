package config

import (
	"crypto/rsa"
	"github.com/gabriel-araujjo/condominio-auth/domain"
	"github.com/dgrijalva/jwt-go"
	"os"
	"io/ioutil"
	"log"
)

type Config struct {
	Dao struct {
		Driver          string
		DNS             string
		VersionStrategy string
		TokenDriver string
		TokenDNS string
	}
	Clients []*domain.Client
	Jwt struct {
		SignatureAlgorithm string
		VerifyKey interface{}
		SignKey   interface{}
	}
}

func getEnv(name string, fallback string) (env string) {
	env = os.Getenv(name)
	if len(env) == 0 {
		env = fallback
	}
	return
}

func getSignKey() *rsa.PrivateKey {
	signBytes, err := ioutil.ReadFile(getEnv("JWT_PRIVATE_KEY_FILE", ""))

	if err != nil {
		log.Fatalf("Can't read private key on path %q", getEnv("JWT_PRIVATE_KEY_FILE", ""))
	}

	key, err :=jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		log.Fatalf("Invalid key format: %e ", err)
	}
	return key
}

func getVerifyKey() *rsa.PublicKey {
	verifyBytes, err := ioutil.ReadFile(getEnv("JWT_PUBLIC_KEY_FILE", ""))

	if err != nil {
		log.Fatalf("Can't read private key on path %q", getEnv("JWT_PRIVATE_KEY_FILE", ""))
	}

	key, err :=jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		log.Fatalf("Invalid key format: %e ", err)
	}
	return key
}

func DefaultConfig() *Config {

	return &Config{
		Dao: struct {
			Driver          string
			DNS             string
			VersionStrategy string
			TokenDriver     string
			TokenDNS		string
		}{
			Driver:          getEnv("DB_DRIVER", "postgres"),
			DNS:             getEnv("DB_DATA_SOURCE_NAME", ""),
			VersionStrategy: getEnv("DB_VERSION_STRATEGY", "psql-versioning"),
			TokenDriver:     getEnv("DB_TOKEN_DRIVER", "mongo"),
			TokenDNS:		 getEnv("DB_TOKEN_DATA_SOURCE_NAME", ""),
		},
		Jwt: struct {
			SignatureAlgorithm string
			VerifyKey interface{}
			SignKey   interface{}
		}{
			SignatureAlgorithm: getEnv("JWT_ALG", "RS512"),
			VerifyKey: getVerifyKey(),
			SignKey: getSignKey(),
		},
		Clients: []*domain.Client{
			{
				Name:     "CondominiumWeb",
				PublicId: "7535b92fcac0ad06d03d",
				Secret:   "64db530fafdc40759c54e1a520a86d0e13e786b3ba215050dbc870fa781651b6",
			},
		},
	}
}
