package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func generateHashSignature(message string, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(message))
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

// VerifySignature verify sign
func VerifySignature(signature string, data string, secret string) bool {
	return signature == generateHashSignature(string(data), secret)
}
