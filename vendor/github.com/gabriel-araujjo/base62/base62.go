package base62

import (
	"fmt"
	"errors"
	"math"
	"strings"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"


func FormatInt(n int64) string {
	if n < 0 {
		return fmt.Sprintf("-%s", FormatUint(uint64(-n)))
	} else {
		return FormatUint(uint64(n))
	}
}

func FormatUint(i uint64) string {
	if i == 0 {
		return "0"
	}
	bytes := make([]byte, 16)
	pos := len(bytes)
	for ; i > 0; pos, i = pos-1,  i/62 {
		bytes[pos-1] = alphabet[i%62]
	}

	return string(bytes[pos:])
}

const maxBeforeOverflow uint64 = math.MaxUint64 / 62
const maxAbsNegativeInt64 = 9223372036854775808

func ParseInt(s string) (int64, error) {
	s = strings.TrimSpace(s)
	var (
		n uint64
		e error
	)
	if s == "" {
		return 0, errors.New(fmt.Sprintf("invalid string %q", s))
	}
	if s[0] == '-' {
		if n, e = ParseUint(s[1:]); e != nil {
			return 0, e
		}
		if n > maxAbsNegativeInt64 {
			return 0, errors.New(fmt.Sprintf("negative overflow with number %q", s))
		}
		return -int64(n), nil
	} else {
		if n, e = ParseUint(s); e != nil {
			return 0, e
		}
		if n > uint64(math.MaxInt64) {
			return 0, errors.New(fmt.Sprintf("positive overflow with number %q", s))
		}
		return int64(n), nil
	}
}

func ParseUint(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	var n uint64
	if s == "" {
		return 0, errors.New(fmt.Sprintf("invalid string %q", s))
	}
	for _, c := range s {
		if n > maxBeforeOverflow {
			return 0, errors.New(fmt.Sprintf("number %q overflows", s))
		}
		n *= 62
		if c >= '0' && c <= '9' {
			n = n + uint64(c) - '0'
		} else if c >= 'A' && c <= 'Z' {
			n = n + uint64(c) - 'A' + 10
		} else if c >= 'a' && c <= 'z' {
			n = n + uint64(c) - 'a' + 36
		} else {
			return 0, errors.New(fmt.Sprintf("unexpected token %q", c))
		}
	}
	return n, nil
}
