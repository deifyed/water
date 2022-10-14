package context

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type testFile struct {
	Path    string
	Content io.Reader
}

func TestAcquireMainLanguageForFile(t *testing.T) {
	testCases := []struct {
		name           string
		withFile       testFile
		expectLanguage string
	}{
		{
			name: "Should work with golang files",
			withFile: testFile{
				Path:    "/main.go",
				Content: strings.NewReader("package main\nfunc main() {}"),
			},
			expectLanguage: "golang",
		},
		{
			name: "Should work with java files",
			withFile: testFile{
				Path:    "/main.java",
				Content: strings.NewReader("public static void main() {}"),
			},
			expectLanguage: "java",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			fs.WriteReader(tc.withFile.Path, tc.withFile.Content)

			language, err := acquireMainLanguageForFile(fs, tc.withFile.Path)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectLanguage, language)
		})
	}
}

func TestAcquireMainLanguageForDir(t *testing.T) {
	testCases := []struct {
		name           string
		withFiles      []testFile
		expectLanguage string
	}{
		{
			name: "Should return golang with only golang files",
			withFiles: []testFile{
				{
					Path:    "/main.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/tests.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/types.go",
					Content: strings.NewReader("package main"),
				},
			},
			expectLanguage: "golang",
		},
		{
			name: "Should return golang with mostly golang files",
			withFiles: []testFile{
				{
					Path:    "/main.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/tests.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/Makefile",
					Content: strings.NewReader(".PHONY: all"),
				},
			},
			expectLanguage: "golang",
		},
		{
			name: "Should return golang with 3/5 files being golang",
			withFiles: []testFile{
				{
					Path:    "/monitor.java",
					Content: strings.NewReader("public static void main() {}"),
				},
				{
					Path:    "/main.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/tests.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/types.go",
					Content: strings.NewReader("package main"),
				},
				{
					Path:    "/lib.java",
					Content: strings.NewReader("public static void main() {}"),
				},
			},
			expectLanguage: "golang",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			for _, item := range tc.withFiles {
				err := fs.WriteReader(item.Path, item.Content)
				assert.NoError(t, err)
			}

			language, err := acquireMainLanguageForDir(fs, "/")
			assert.NoError(t, err)

			assert.Equal(t, tc.expectLanguage, language)
		})
	}
}
