package guard

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"hash"
	"strconv"
	"time"
)

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateSignature(t time.Time, algorithm string, data []byte, secret []byte) (string, error) {
	var mac hash.Hash
	switch algorithm {
	case "HMAC-SHA256":
		mac = hmac.New(sha256.New, secret)
	case "HMAC-SHA1":
		mac = hmac.New(sha1.New, secret)
	case "HMAC-SHA512":
		mac = hmac.New(sha512.New, secret)
	default:
		return "", errors.New("unsupported algorithm: " + algorithm)
	}
	mac.Write(data)
	mac.Write([]byte(strconv.Itoa(int(t.Unix())))) // Append timestamp to the data
	return hex.EncodeToString(mac.Sum(nil)), nil
}
