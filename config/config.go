package config

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gabriel-araujjo/condominio-auth/domain"
)

// Config the app config
type Config struct {
	Dao     Dao
	Clients []*domain.Client
	Jwt     Jwt
	Session Session
}

// Dao store config about Dao
type Dao struct {
	Driver          string
	URI             string
	VersionStrategy string
	TokenDriver     string
	TokenURI        string
}

// Jwt store config about Jwt tokens
type Jwt struct {
	SignatureAlgorithm string
	VerifyKey          interface{}
	SignKey            interface{}
}

// Session store configuration about session
type Session struct {
	// StoreType is the type of the uri
	StoreType string
	// SoreURI is the identifier of the store
	StoreURI string
	// PoolSize is the maximum number of idle connections to store
	// Valid on redis store type
	PoolSize int
	// CookieName is the name of the cookie stored on client
	CookieName string
	// HashKey is used to authenticate the cookie value using HMAC.
	// It is recommended to use a key with 32 or 64 bytes.
	HashKey []byte
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

	key, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

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

	key, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		log.Fatalf("Invalid key format: %e ", err)
	}
	return key
}

func mustParseInt(intString string) int {
	value, err := strconv.Atoi(intString)
	if err != nil {
		panic(fmt.Sprintf("expecting an integer instead of %q", intString))
	}
	return value
}

func mustDecodeHex(hexString string) []byte {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		panic(fmt.Sprintf("can't decode hex string %q", hexString))
	}
	return data
}

// DefaultConfig returns the app default configuration
func DefaultConfig() *Config {

	return &Config{
		Dao: Dao{
			Driver:          getEnv("DATABASE_DRIVER", "postgres"),
			URI:             getEnv("DATABASE_URL", ""),
			VersionStrategy: getEnv("DATABASE_VERSION_STRATEGY", "psql-versioning"),
			TokenDriver:     getEnv("DATABASE_TOKEN_DRIVER", "mongo"),
			TokenURI:        getEnv("DATABASE_TOKEN_URL", ""),
		},
		Jwt: Jwt{
			SignatureAlgorithm: getEnv("JWT_ALG", "RS512"),
			VerifyKey:          getVerifyKey(),
			SignKey:            getSignKey(),
		},
		Clients: []*domain.Client{
			{
				Name:     "CondominiumWeb",
				PublicID: "7535b92fcac0ad06d03d",
				Secret:   "64db530fafdc40759c54e1a520a86d0e13e786b3ba215050dbc870fa781651b6",
			},
		},
		Session: Session{
			StoreType:  getEnv("SESSIONS_STORE_TYPE", "redis"),
			StoreURI:   getEnv("SESSIONS_STORE_URL", ""),
			PoolSize:   mustParseInt(getEnv("SESSIONS_STORE_POOL_SIZE", "10")),
			CookieName: "sessions",
			HashKey: mustDecodeHex(getEnv("COOKIE_CODEC_HASH_KEY",
				"75625f538a4a5431762b96263e2762fb1cd8af1a3326c4468aaa9a7f336ed0ccf27dfd59167f1dd64aa28074ef87726b0c1f7f7d68fedd6f825e5323dba23280")),
		},
	}
}
