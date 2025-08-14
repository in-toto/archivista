package sqlstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureMySQLConnectionString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "already has tcp protocol",
			input:       "user:pass@tcp(localhost:3306)/dbname",
			expected:    "user:pass@tcp(localhost:3306)/dbname",
			expectError: false,
		},
		{
			name:        "needs tcp protocol",
			input:       "user:pass@localhost:3306/dbname",
			expected:    "user:pass@tcp(localhost:3306)/dbname",
			expectError: false,
		},
		{
			name:        "with mysql:// prefix",
			input:       "mysql://user:pass@localhost:3306/dbname",
			expected:    "user:pass@tcp(localhost:3306)/dbname",
			expectError: false,
		},
		{
			name:        "invalid url format",
			input:       "invalid:url:format",
			expected:    "",
			expectError: true,
		},
		{
			name:        "with query parameters",
			input:       "user:pass@localhost:3306/dbname?param=value",
			expected:    "user:pass@tcp(localhost:3306)/dbname?param=value",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ensureMySQLConnectionString(tt.input)
			if tt.expectError {
				require.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
