package cli

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type InvalidFixture struct {
	args            []string
	expectedMessage string
}

func TestInvalidCommands(t *testing.T) {
	for _, fixture := range []InvalidFixture{
		{
			args:            []string{"unknown"},
			expectedMessage: "unknown subcommand: unknown",
		},
		{
			args:            []string{"generate", "--unknown"},
			expectedMessage: "flag provided but not defined: -unknown",
		},
		{
			args:            []string{"generate", "unknown", "unknown"},
			expectedMessage: "expected no positional arguments, got 2",
		},
	} {
		t.Run(strings.Join(fixture.args, " "), func(t *testing.T) {
			err := Run(fixture.args)
			assert.NotNilf(t, err, "expected error")
			assert.Equal(t, fixture.expectedMessage, err.Error())
		})
	}
}
