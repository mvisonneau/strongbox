package rand

import (
  "math/rand"
  "time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-.?!{}[]|"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func String(length int) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}
