package tgauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrMissing      = errors.New("tgauth: init data missing")
	ErrInvalidFormat = errors.New("tgauth: invalid init data format")
	ErrMissingHash   = errors.New("tgauth: hash missing")
	ErrInvalidHash   = errors.New("tgauth: invalid hash")
	ErrExpired       = errors.New("tgauth: init data expired")
	ErrMissingUser   = errors.New("tgauth: user missing")
)

// User is the Telegram user embedded in WebApp initData.
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Validate checks Telegram WebApp initData signature and freshness.
// See https://core.telegram.org/bots/webapps#validating-data-received-via-the-mini-app
func Validate(initData, botToken string, maxAge time.Duration) error {
	if strings.TrimSpace(initData) == "" {
		return ErrMissing
	}
	if strings.TrimSpace(botToken) == "" {
		return ErrInvalidFormat
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidFormat, err)
	}

	hash := values.Get("hash")
	if hash == "" {
		return ErrMissingHash
	}

	pairs := make([]string, 0, len(values))
	var authDate time.Time
	for key, vals := range values {
		if key == "hash" || len(vals) == 0 {
			continue
		}
		if key == "auth_date" {
			sec, err := strconv.ParseInt(vals[0], 10, 64)
			if err != nil {
				return fmt.Errorf("%w: bad auth_date", ErrInvalidFormat)
			}
			authDate = time.Unix(sec, 0)
		}
		pairs = append(pairs, key+"="+vals[0])
	}
	sort.Strings(pairs)

	if maxAge > 0 {
		if authDate.IsZero() {
			return ErrInvalidFormat
		}
		if time.Since(authDate) > maxAge {
			return ErrExpired
		}
	}

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(botToken))

	mac := hmac.New(sha256.New, secret.Sum(nil))
	mac.Write([]byte(strings.Join(pairs, "\n")))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(hash)) {
		return ErrInvalidHash
	}
	return nil
}

// ParseUser extracts the Telegram user from validated initData.
func ParseUser(initData string) (*User, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFormat, err)
	}
	raw := values.Get("user")
	if raw == "" {
		return nil, ErrMissingUser
	}
	var user User
	if err := json.Unmarshal([]byte(raw), &user); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFormat, err)
	}
	if user.ID == 0 {
		return nil, ErrMissingUser
	}
	return &user, nil
}
