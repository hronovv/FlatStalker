package tgauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func sign(botToken string, params map[string]string) string {
	pairs := make([]string, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)

	secret := hmac.New(sha256.New, []byte("WebAppData"))
	secret.Write([]byte(botToken))
	mac := hmac.New(sha256.New, secret.Sum(nil))
	mac.Write([]byte(strings.Join(pairs, "\n")))
	return hex.EncodeToString(mac.Sum(nil))
}

func makeInitData(botToken string, params map[string]string) string {
	hash := sign(botToken, params)
	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("hash", hash)
	return q.Encode()
}

func TestValidateAndParseUser(t *testing.T) {
	token := "123456:ABC-DEF"
	authDate := strconv.FormatInt(time.Now().Unix(), 10)
	userJSON := `{"id":42,"first_name":"Ivan","username":"ivan"}`
	initData := makeInitData(token, map[string]string{
		"auth_date": authDate,
		"user":      userJSON,
		"query_id":  "AAE",
	})

	if err := Validate(initData, token, time.Hour); err != nil {
		t.Fatalf("validate: %v", err)
	}

	user, err := ParseUser(initData)
	if err != nil {
		t.Fatalf("parse user: %v", err)
	}
	if user.ID != 42 || user.FirstName != "Ivan" {
		t.Fatalf("user=%+v", user)
	}
}

func TestValidateRejectsBadHash(t *testing.T) {
	token := "123456:ABC-DEF"
	authDate := strconv.FormatInt(time.Now().Unix(), 10)
	initData := makeInitData(token, map[string]string{
		"auth_date": authDate,
		"user":      `{"id":1}`,
	})
	bad := strings.Replace(initData, "hash=", "hash=deadbeef", 1)
	if !strings.Contains(bad, "hash=deadbeef") {
		bad = initData + "x"
	}
	// Force invalid hash value.
	values, _ := url.ParseQuery(initData)
	values.Set("hash", "0000000000000000000000000000000000000000000000000000000000000000")
	bad = values.Encode()

	if err := Validate(bad, token, time.Hour); err == nil {
		t.Fatal("expected invalid hash")
	}
}

func TestValidateRejectsExpired(t *testing.T) {
	token := "123456:ABC-DEF"
	old := strconv.FormatInt(time.Now().Add(-48*time.Hour).Unix(), 10)
	initData := makeInitData(token, map[string]string{
		"auth_date": old,
		"user":      `{"id":1}`,
	})
	if err := Validate(initData, token, time.Hour); err == nil {
		t.Fatal("expected expired")
	}
}

func TestValidateRejectsMissing(t *testing.T) {
	if err := Validate("", "token", time.Hour); err == nil {
		t.Fatal("expected missing")
	}
}
