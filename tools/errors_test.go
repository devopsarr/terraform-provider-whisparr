package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrDataNotFoundError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		kind, field, search, expected string
	}{
		"tag": {"whisparr_tag", "label", "test", "data source not found: no whisparr_tag with label 'test'"},
	}
	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, ErrDataNotFoundError(test.kind, test.field, test.search).Error())
		})
	}
}
