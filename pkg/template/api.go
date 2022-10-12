package template

import (
	"io"

	"github.com/deifyed/water/pkg/context"
	"github.com/spf13/afero"
)

func Discover(fs *afero.Afero, templateDir string, context context.Context) (io.Reader, error) {
	return nil, nil
}
