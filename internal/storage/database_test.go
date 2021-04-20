package storage

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

const dbConnArg = "db-conn"

var connectionString string

func TestDatabase(t *testing.T) {
	if flag.Lookup(dbConnArg) == nil {
		flag.StringVar(&connectionString, dbConnArg, "", "Database connection string to run database tests")
	}
	flag.Parse()

	if connectionString == "" {
		t.Skipf("Skipping database tests as no '%s' provided", dbConnArg)
	}

	t.Run("check add banner slot and group", func(t *testing.T) {
		require.True(t, false, "write me")
	})

	t.Run("check add rotation", func(t *testing.T) {
		require.True(t, false, "write me")
	})

	t.Run("check delete rotation", func(t *testing.T) {
		require.True(t, false, "write me")
	})

	t.Run("check delete cascade", func(t *testing.T) {
		require.True(t, false, "write me")
	})

	t.Run("check total shows", func(t *testing.T) {
		require.True(t, false, "write me")
	})
}
