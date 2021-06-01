package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.NotEqual(t, String(10), String(10), "strings should be different")
	assert.Len(t, String(10), 10, "length of string should be 10")
}
