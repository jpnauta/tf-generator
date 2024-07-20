package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type HasDirsFixture struct {
	filePath string
	hasDirs  bool
}

func TestHasDirsFixtures(t *testing.T) {
	for _, fixture := range []HasDirsFixture{
		{
			filePath: "test.txt",
			hasDirs:  false,
		},
		{
			filePath: "abc",
			hasDirs:  false,
		},
		{
			filePath: "./test.txt",
			hasDirs:  true,
		},
		{
			filePath: "test/test.txt",
			hasDirs:  true,
		},
	} {
		t.Run(fixture.filePath, func(t *testing.T) {
			assert.Equal(t, fixture.hasDirs, hasDirs(fixture.filePath))
		})
	}
}
