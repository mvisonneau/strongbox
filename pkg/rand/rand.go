package rand

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// String : Generates a random string of a given length
func String(length int) string {
	b := make([]byte, base64.RawStdEncoding.DecodedLen(length))
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
