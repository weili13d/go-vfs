package vfst

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	vfs "github.com/twpayne/go-vfs/v4"
)

func TestWalk(t *testing.T) {
	fileSystem, cleanup, err := NewTestFS(map[string]interface{}{
		"/home/user/.bashrc":  "# .bashrc contents\n",
		"/home/user/skip/foo": "bar",
		"/home/user/symlink":  &Symlink{Target: "baz"},
	})
	require.NoError(t, err)
	defer cleanup()
	pathTypeMap := make(map[string]fs.FileMode)
	require.NoError(t, vfs.Walk(fileSystem, "/", func(path string, info fs.FileInfo, err error) error {
		assert.NoError(t, err)
		pathTypeMap[filepath.ToSlash(path)] = info.Mode() & fs.ModeType
		if filepath.Base(path) == "skip" {
			return vfs.SkipDir
		}
		return nil
	}))
	expectedPathTypeMap := map[string]fs.FileMode{
		"/":                  fs.ModeDir,
		"/home":              fs.ModeDir,
		"/home/user":         fs.ModeDir,
		"/home/user/.bashrc": 0,
		"/home/user/skip":    fs.ModeDir,
		"/home/user/symlink": fs.ModeSymlink,
	}
	assert.Equal(t, expectedPathTypeMap, pathTypeMap)
}
