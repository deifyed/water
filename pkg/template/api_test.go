package template

import (
	"encoding/json"
	"io"
	"path"
	"testing"

	"github.com/deifyed/water/pkg/context"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDiscover(t *testing.T) {
	testCases := []struct {
		name            string
		withContext     context.Context
		withFs          *afero.Afero
		withTemplateDir string
		expectTemplate  string
	}{
		{
			name: "Should find relevant file",
			withContext: context.Context{
				TargetType: context.TargetTypeFile,
				TargetPath: "/Makefile",
				Tags:       map[string]string{"language": "golang"},
			},
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				err := fs.MkdirAll("/templates/Makefile", 0o700)
				assert.NoError(t, err)

				metadatas := []metadata{
					{
						Target: "./golang",
						Tags: map[string]string{
							"language": "golang",
						},
					},
				}

				raw, err := json.Marshal(metadatas)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/metadata.json", raw, 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/golang", []byte("golang makefile content"), 0o600)
				assert.NoError(t, err)

				return fs
			}(),
			withTemplateDir: "/templates",
			expectTemplate:  "golang makefile content",
		},
		{
			name: "Should find relevant file when same file alternatives exist",
			withContext: context.Context{
				TargetType: context.TargetTypeFile,
				TargetPath: "/Makefile",
				Tags:       map[string]string{"language": "golang"},
			},
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				err := fs.MkdirAll("/templates/Makefile", 0o700)
				assert.NoError(t, err)

				metadatas := []metadata{
					{Target: "./golang", Tags: map[string]string{"language": "golang"}},
					{Target: "./java", Tags: map[string]string{"language": "java"}},
				}

				raw, err := json.Marshal(metadatas)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/metadata.json", raw, 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/golang", []byte("golang makefile content"), 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/java", []byte("java makefile content"), 0o600)
				assert.NoError(t, err)

				return fs
			}(),
			withTemplateDir: "/templates",
			expectTemplate:  "golang makefile content",
		},
		{
			name: "Should find relevant file when same file alternatives exist and other templates",
			withContext: context.Context{
				TargetType: context.TargetTypeFile,
				TargetPath: "/Makefile",
				Tags:       map[string]string{"language": "golang"},
			},
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				err := fs.MkdirAll("/templates/Makefile", 0o700)
				assert.NoError(t, err)

				metadatas := []metadata{
					{Target: "./golang", Tags: map[string]string{"language": "golang"}},
					{Target: "./java", Tags: map[string]string{"language": "java"}},
				}

				raw, err := json.Marshal(metadatas)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/metadata.json", raw, 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/golang", []byte("golang makefile content"), 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/Makefile/java", []byte("java makefile content"), 0o600)
				assert.NoError(t, err)

				err = fs.MkdirAll("/templates/package.json", 0o700)
				assert.NoError(t, err)

				metadatas = []metadata{
					{
						Target: "./vue",
						Tags:   map[string]string{"language": "javascript", "framework": "vue"},
					},
				}

				raw, err = json.Marshal(metadatas)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/package.json/metadata.json", raw, 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile("/templates/package.json/vue", []byte("vue package.json content"), 0o600)
				assert.NoError(t, err)

				return fs
			}(),
			withTemplateDir: "/templates",
			expectTemplate:  "golang makefile content",
		},
		{
			name:            "Should interpolate project name",
			withTemplateDir: "/templates",
			withContext: context.Context{
				TargetType: context.TargetTypeFile,
				TargetPath: "/horse/Makefile",
				Tags:       map[string]string{"name": "Makefile", "language": "golang"},
			},
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				scaffoldableItemDir := path.Join("/", "templates", "Makefile")

				err := fs.MkdirAll(scaffoldableItemDir, 0o700)
				assert.NoError(t, err)

				metadatas := []metadata{{Target: "./golang", Tags: map[string]string{"name": "Makefile", "language": "golang"}}}

				raw, err := json.Marshal(metadatas)
				assert.NoError(t, err)

				err = fs.WriteFile(path.Join(scaffoldableItemDir, "metadata.json"), raw, 0o600)
				assert.NoError(t, err)

				err = fs.WriteFile(path.Join(scaffoldableItemDir, "golang"), []byte("golang {{ .ProjectName }} content"), 0o600)
				assert.NoError(t, err)

				return fs
			}(),
			expectTemplate: "golang horse content",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := Discover(&mockLogger{t: t}, tc.withFs, tc.withTemplateDir, tc.withContext)
			assert.NoError(t, err)

			raw, err := io.ReadAll(result)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectTemplate, string(raw))
		})
	}
}

type mockLogger struct {
	t *testing.T
}

func (m mockLogger) Debugf(format string, args ...interface{}) {
	m.t.Logf(format, args...)
}
