package util_test

import (
	"testing"

	"github.com/night1010/everhealth/util"
	"github.com/stretchr/testify/assert"
)

func TestIsMemberOf(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		nums := []int{1, 2, 3}
		num := 2

		isMemberOf := util.IsMemberOf(nums, num)

		assert.True(t, isMemberOf)
	})
	t.Run("false", func(t *testing.T) {
		nums := []int{1, 2, 3}
		num := 4

		isMemberOf := util.IsMemberOf(nums, num)

		assert.False(t, isMemberOf)
	})
}
