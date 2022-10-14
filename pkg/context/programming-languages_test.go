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
