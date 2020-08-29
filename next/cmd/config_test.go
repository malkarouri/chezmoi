package cmd

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vfs "github.com/twpayne/go-vfs"
	xdg "github.com/twpayne/go-xdg/v3"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
	"github.com/twpayne/chezmoi/next/internal/chezmoitest"
)

func TestAddTemplateFuncPanic(t *testing.T) {
	c := newTestConfig(t, nil)
	assert.NotPanics(t, func() {
		c.addTemplateFunc("func", nil)
	})
	assert.Panics(t, func() {
		c.addTemplateFunc("func", nil)
	})
}

func TestUpperSnakeCaseToCamelCase(t *testing.T) {
	for s, want := range map[string]string{
		"BUG_REPORT_URL":   "bugReportURL",
		"ID":               "id",
		"ID_LIKE":          "idLike",
		"NAME":             "name",
		"VERSION_CODENAME": "versionCodename",
		"VERSION_ID":       "versionID",
	} {
		assert.Equal(t, want, upperSnakeCaseToCamelCase(s))
	}
}

func TestValidateKeys(t *testing.T) {
	for _, tc := range []struct {
		data    interface{}
		wantErr bool
	}{
		{
			data:    nil,
			wantErr: false,
		},
		{
			data: map[string]interface{}{
				"foo":                    "bar",
				"a":                      0,
				"_x9":                    false,
				"ThisVariableIsExported": nil,
				"αβ":                     "",
			},
			wantErr: false,
		},
		{
			data: map[string]interface{}{
				"foo-foo": "bar",
			},
			wantErr: true,
		},
		{
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar-bar": "baz",
				},
			},
			wantErr: true,
		},
		{
			data: map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"bar-bar": "baz",
					},
				},
			},
			wantErr: true,
		},
	} {
		if tc.wantErr {
			assert.Error(t, validateKeys(tc.data, identifierRegexp))
		} else {
			assert.NoError(t, validateKeys(tc.data, identifierRegexp))
		}
	}
}

func newTestConfig(t *testing.T, fs vfs.FS, options ...configOption) *Config {
	system := chezmoi.NewRealSystem(fs, chezmoitest.NewPersistentState())
	c, err := newConfig(append(
		[]configOption{
			withBaseSystem(system),
			withDestSystem(system),
			withSourceSystem(system),
			withTestFS(fs),
			withTestUser("user"),
		},
		options...,
	)...)
	require.NoError(t, err)
	return c
}

func withBaseSystem(baseSystem chezmoi.System) configOption {
	return func(c *Config) {
		c.baseSystem = baseSystem
	}
}

func withDestSystem(destSystem chezmoi.System) configOption {
	return func(c *Config) {
		c.destSystem = destSystem
	}
}

func withSourceSystem(sourceSystem chezmoi.System) configOption {
	return func(c *Config) {
		c.sourceSystem = sourceSystem
	}
}

func withTestFS(fs vfs.FS) configOption {
	return func(c *Config) {
		c.fs = fs
	}
}

func withTestUser(username string) configOption {
	return func(c *Config) {
		var homeDir string
		switch runtime.GOOS {
		case "windows":
			homeDir = "C:\\home\\user"
		default:
			homeDir = "/home/user"
		}
		c.SourceDir = osPath(filepath.Join(homeDir, ".local", "share", "chezmoi"))
		c.DestDir = osPath(homeDir)
		c.Umask = 0o22
		c.bds = &xdg.BaseDirectorySpecification{
			ConfigHome: filepath.Join(homeDir, ".config"),
			DataHome:   filepath.Join(homeDir, ".local"),
			CacheHome:  filepath.Join(homeDir, ".cache"),
			RuntimeDir: filepath.Join(homeDir, ".run"),
		}
	}
}