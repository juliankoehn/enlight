package enlight

import (
	"testing"

	testify "github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	e := New()
	c := e.NewContext()
	assert := testify.New(t)
	assert.Equal(e, c.Enlight())
}
