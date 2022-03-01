package hash_test

import (
	"shorturl/pkg/hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlShortHash(t *testing.T) {
	t.Parallel() // 允许与其他测试平行
	newAssert := assert.New(t)
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "abcd",
			want:  "bphWu6",
		},
		{
			input: "suuuuuuuuuuuuuuperlong",
			want:  "cDyAQZ",
		},
	}
	t.Run("短链接哈希测试", func(t *testing.T) {
		for _, item := range tests {
			getshort := hash.UrlShortHash(item.input)
			newAssert.Equal(item.want, getshort)
		}
	})
}
