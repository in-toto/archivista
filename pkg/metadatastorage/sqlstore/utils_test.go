// Copyright 2025 The Archivista Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigFromMySQL(t *testing.T) {
	t.Run("tcp host and port", func(t *testing.T) {
		c, user, password, err := ConfigFromMySQL("user:pass@tcp(127.0.0.1:3306)/testify")
		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", c.Host)
		assert.Equal(t, 3306, c.Port)
		assert.Equal(t, "testify", c.DB)
		assert.Equal(t, "user", user)
		assert.Equal(t, "pass", password)
	})

	t.Run("bracketed IPv6 host", func(t *testing.T) {
		// Splitting on ":" left Host as "[" and the port empty.
		c, _, _, err := ConfigFromMySQL("user:pass@tcp([::1]:3306)/testify")
		require.NoError(t, err)
		assert.Equal(t, "::1", c.Host)
		assert.Equal(t, 3306, c.Port)
	})

	t.Run("unix socket is rejected, not a panic", func(t *testing.T) {
		// ParseDSN only defaults a port for tcp, so Addr here is a bare path
		// with no colon at all.
		_, _, _, err := ConfigFromMySQL("user:pass@unix(/tmp/mysql.sock)/testify")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "only tcp is supported")
	})

	t.Run("unparseable dsn", func(t *testing.T) {
		_, _, _, err := ConfigFromMySQL("not a dsn")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parsing connection string")
	})
}
