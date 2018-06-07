package sessions

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/boj/redistore"
	"github.com/gabriel-araujjo/condominio-auth/config"
)

const redisDefaultPort = "6379"

var pathDBRegexp = regexp.MustCompile(`/(\d+)\z`)

// NewRedisStore creates a middleware for gin that handle session
// config: The configuration of the redis tore
func NewRedisStore(config *config.Config) (Store, error) {
	network, _, address, password, db, err := parseURI(config.Session.StoreURI)

	if err != nil {
		return nil, err
	}

	s, err := redistore.NewRediStoreWithDB(config.Session.PoolSize, network, address, password, db, config.Session.HashKey)

	if err != nil {
		return nil, err
	}
	return AdaptGorillaStore(s, s), nil
}

func parseURI(rawURI string) (network string, tls bool, address string, password string, db string, err error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return
	}

	if u.Scheme != "redis" && u.Scheme != "rediss" {
		err = fmt.Errorf("invalid redis URL scheme: %s", u.Scheme)
		return
	}

	match := pathDBRegexp.FindStringSubmatch(u.Path)
	if len(match) == 2 {
		db = match[1]
	} else if u.Path == "" || u.Path == "/" {
		db = "0"
	} else {
		err = fmt.Errorf("invalid database: %s", u.Path[1:])
		return
	}

	network = "tcp"
	tls = u.Scheme == "rediss"
	address = canonicalHost(u.Host)
	password, _ = u.User.Password()
	return
}

func canonicalHost(host string) string {
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		// assume port is missing
		port = redisDefaultPort
	}
	if host == "" {
		host = "localhost"
	}
	return net.JoinHostPort(host, port)
}
