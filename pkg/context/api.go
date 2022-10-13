package context

import (
	"fmt"
	"path"

	"github.com/spf13/afero"
)

func GatherContext(fs *afero.Afero, targetPath string) (Context, error) {
	info, err := fs.Stat(targetPath)
	if err != nil {
		return Context{}, fmt.Errorf("statting: %w", err)
	}

	ctx := Context{TargetPath: targetPath}
	var target string

	if info.IsDir() {
		ctx.TargetType = TargetTypeDirectory
		target = targetPath
	} else {
		ctx.TargetType = TargetTypeFile
		target = path.Dir(targetPath)
	}

	lang, err := acquireMainLanguageForDir(fs, target)
	if err != nil {
		return Context{}, fmt.Errorf("acquiring language for target %s: %w", target, err)
	}

	ctx.Tags["language"] = lang

	return ctx, nil
}
