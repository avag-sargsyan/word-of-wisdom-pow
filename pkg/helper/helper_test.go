package helper_test

import (
	"testing"

	"github.com/avag-sargsyan/word-of-wisdom-pow/pkg/helper"
	"github.com/stretchr/testify/assert"
)

// Black-box testing helpers

func TestRandomString(t *testing.T) {
	len := 20
	str1 := helper.RandomString(len)
	str2 := helper.RandomString(len)

	assert.NotEqual(t, str1, str2)
}
