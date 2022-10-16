package context

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

func GatherContext(log logger, fs *afero.Afero, targetPath string) (Context, error) {
	info, err := fs.Stat(targetPath)
	if err != nil {
		return Context{}, fmt.Errorf("statting: %w", err)
	}

	ctx := Context{TargetType: TargetTypeFile, Tags: make(map[string]string)}

	ctx.TargetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return Context{}, fmt.Errorf("acquiring absolute path for %s: %w", targetPath, err)
	}

	languageSearchDir := path.Dir(ctx.TargetPath)

	if info.IsDir() {
		ctx.TargetType = TargetTypeDirectory
		languageSearchDir = ctx.TargetPath
	}

	lang, err := acquireMainLanguageForDir(fs, languageSearchDir)
	if err != nil {
		return Context{}, fmt.Errorf("acquiring main language: %w", err)
	}

	ctx.Tags["name"] = path.Base(ctx.TargetPath)
	ctx.Tags["language"] = lang

	return ctx, nil
}
