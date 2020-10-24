package spotify

import (
	"errors"
	"net/url"

	"github.com/musicmash/artisync/internal/log"
)

var (
	errAuthFailed  = errors.New("auth failed")
	errCodeIsEmpty = errors.New("query arg 'code' is empty")
)

func validateQuery(values url.Values) error {
	if e := values.Get("error"); e != "" {
		log.Errorf("auth failed: %v", e)

		return errAuthFailed
	}

	code := values.Get("code")
	if code == "" {
		return errCodeIsEmpty
	}

	return nil
}
