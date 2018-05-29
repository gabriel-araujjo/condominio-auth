package redis

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gabriel-araujjo/condominio-auth/config"
	"github.com/gin-gonic/contrib/sessions"
)

type redisParams struct {
	network  string
	address  string
	password string
	database int
	tls      bool
}

var pathDBRegexp = regexp.MustCompile(`/(\d*)\z`)

func parseURI(rawURI string) (*redisParams, error) {
	u, err := url.Parse(rawURI)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return nil, fmt.Errorf("invalid redis URL scheme: %s", u.Scheme)
	}

	params := &redisParams{}

	// As per the IANA draft spec, the host defaults to localhost and
	// the port defaults to 6379.
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		// assume port is missing
		host = u.Host
		port = "6379"
	}
	if host == "" {
		host = "localhost"
	}
	params.address = net.JoinHostPort(host, port)

	if u.User != nil {
		password, isSet := u.User.Password()
		if isSet {
			params.password = password
		}
	}

	match := pathDBRegexp.FindStringSubmatch(u.Path)
	if len(match) == 2 {
		db := 0
		if len(match[1]) > 0 {
			db, err = strconv.Atoi(match[1])
			if err != nil {
				return nil, fmt.Errorf("invalid database: %s", u.Path[1:])
			}
		}
		if db != 0 {
			params.database = db
		}
	} else if u.Path != "" {
		return nil, fmt.Errorf("invalid database: %s", u.Path[1:])
	}

	params.tls = u.Scheme == "rediss"
	params.network = "tcp"

	return params, nil
}

// NewStore creates a middleware for gin that handle session
func NewStore(config *config.Config) (sessions.RedisStore, io.Closer, error) {
	params, err := parseURI(config.Sessions.StoreURI)

	if err != nil {
		return nil, nil, err
	}

	if params.database != 0 {
		return nil, nil, fmt.Errorf("redis: only redis uri without database are currently accepted %q", config.Sessions.StoreURI)
	}

	store, err := sessions.NewRedisStore(config.Sessions.MaxConnections, params.network, params.address, params.password, config.Sessions.CookieCodecHashKey)
	//FIXME: Find out how close redis store
	return store, nil, err
}
