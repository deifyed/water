package context

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/afero"
)

var extensionToLanguage = map[string]string{
	".go":   "golang",
	".java": "java",
	".js":   "javascript",
	".py":   "python",
	".rb":   "ruby",
	".sh":   "shell",
	".ts":   "typescript",
}

// acquireMainLanguageForDir attempts to determine the main programming language of the target directory
func acquireMainLanguageForDir(fs *afero.Afero, targetDir string) (string, error) {
	numberOfFiles := 0
	languageCount := make(map[string]int)

	fs.Walk(targetDir, func(targetPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		numberOfFiles++

		lang, err := acquireMainLanguageForFile(fs, targetPath)
		if err != nil {
			return fmt.Errorf("acquiring language for file %s: %w", targetPath, err)
		}

		_, ok := extensionToLanguage[lang]
		if !ok {
			languageCount[lang] = 0
		}

		languageCount[lang]++

		return nil
	})

	strongestLanguage := ""
	strongestLanguageCount := 0

	for key, val := range languageCount {
		if val > strongestLanguageCount {
			strongestLanguage = key
			strongestLanguageCount = val
		}
	}

	return strongestLanguage, nil
}

// acquireMainLanguageForFile attempts to determine the main programming language of the target file
func acquireMainLanguageForFile(fs *afero.Afero, targetFile string) (string, error) {
	ext := path.Ext(targetFile)

	if language, ok := extensionToLanguage[ext]; ok {
		return language, nil
	}

	return "", nil
}
