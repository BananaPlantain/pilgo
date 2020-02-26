package osfs_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/osfs"
)

var _ fs.FileSystem = new(osfs.FileSystem)

func TestFileSystem(t *testing.T) {
	t.Run("Info", testFileSystemInfo)
	t.Run("ReadDir", testFileSystemReadDir)
}

func testFileSystemInfo(t *testing.T) {
	testCases := []struct {
		filename string
		exists   bool
		isDir    bool
		linkname string
		err      error
	}{
		{
			filename: "exists",
			exists:   true,
			isDir:    false,
			linkname: "",
			err:      nil,
		},
		{
			filename: "not_exists",
			exists:   false,
			isDir:    false,
			linkname: "",
			err:      nil,
		},
		{
			filename: "directory",
			exists:   true,
			isDir:    true,
			linkname: "",
			err:      nil,
		},
		{
			filename: "symlink",
			exists:   true,
			isDir:    false,
			linkname: "directory",
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       osfs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			fi, err := fs.Info(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			t.Run("Exists", func(t *testing.T) {
				if want, got := tc.exists, fi.Exists(); want != got {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
			t.Run("IsDir", func(t *testing.T) {
				if want, got := tc.isDir, fi.IsDir(); want != got {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
			t.Run("Linkname", func(t *testing.T) {
				if want, got := tc.linkname, fi.Linkname(); want != got {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		})
	}
}

func testFileSystemReadDir(t *testing.T) {
	testCases := []struct {
		filename string
		want     []string
		err      error
	}{
		{
			filename: "directory",
			want:     []string{"bar", "foo"},
			err:      nil,
		},
		{
			filename: "symlink",
			want:     []string{"bar", "foo"},
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       osfs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			files, err := fs.ReadDir(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, files; !cmp.Equal(got, want) {
				t.Errorf(
					"(*FileSystem).ReadDir mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
		})
	}
}