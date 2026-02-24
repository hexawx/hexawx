package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitVersion(t *testing.T) {

	t.Run("mon-plugin - Doit retourner unknown", func(t *testing.T) {
		expected := "unknown"
		s := getVersionFromPath("mon-plugin")
		assert.Equal(t, expected, s)
	})

	t.Run("mon-plugin.exe - Doit retourner unknown", func(t *testing.T) {
		expected := "unknown"
		s := getVersionFromPath("mon-plugin.exe")
		assert.Equal(t, expected, s)
	})

	t.Run("plugins/mon-plugin - Doit retourner unknown", func(t *testing.T) {
		expected := "unknown"
		s := getVersionFromPath("plugins/mon-plugin")
		assert.Equal(t, expected, s)
	})

	t.Run("plugins/mon-plugin.exe - Doit retourner unknown", func(t *testing.T) {
		expected := "unknown"
		s := getVersionFromPath("plugins/mon-plugin.exe")
		assert.Equal(t, expected, s)
	})

	t.Run("mon-plugin_1.2.3 - Doit retourner 1.2.3", func(t *testing.T) {
		expected := "1.2.3"
		s := getVersionFromPath("mon-plugin_1.2.3")
		assert.Equal(t, expected, s)
	})

	t.Run("plugins/mon-plugin_1.2.3 - Doit retourner 1.2.3", func(t *testing.T) {
		expected := "1.2.3"
		s := getVersionFromPath("plugins/mon-plugin_1.2.3")
		assert.Equal(t, expected, s)
	})

	t.Run("mon-plugin_1.2.3.exe - Doit retourner 1.2.3", func(t *testing.T) {
		expected := "1.2.3"
		s := getVersionFromPath("mon-plugin_1.2.3.exe")
		assert.Equal(t, expected, s)
	})

	t.Run("plugins/mon-plugin_1.2.3.exe - Doit retourner 1.2.3", func(t *testing.T) {
		expected := "1.2.3"
		s := getVersionFromPath("plugins/mon-plugin_1.2.3.exe")
		assert.Equal(t, expected, s)
	})

}
