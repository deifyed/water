package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	iofs "io/fs"
	"path"
	"text/template"

	"github.com/deifyed/water/pkg/context"
	"github.com/spf13/afero"
)

// Discover returns a relevant template for the given context
func Discover(log logger, fs *afero.Afero, templateDir string, context context.Context) (io.Reader, error) {
	metadatas, err := gatherMetadataForTemplateEntities(log, fs, templateDir)
	if err != nil {
		return nil, fmt.Errorf("gathering metadata: %w", err)
	}

	log.Debugf("Assessing %d metadatas", len(metadatas))
	topHit := getMostRelevantMetadata(log, context.Tags, metadatas)

	content, err := fs.ReadFile(path.Clean(topHit.Metadata.Target))
	if err != nil {
		return nil, fmt.Errorf("reading top hit target %s: %w", topHit.Metadata.Target, err)
	}

	t, err := template.New("tophit").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func getMostRelevantMetadata(log logger, tags map[string]string, metadatas []metadata) metadataHit {
	topHit := metadataHit{}
	logEntry := make([]string, len(metadatas))

	for index, metadata := range metadatas {
		hitrate := calculateHitrate(tags, metadata.Tags)

		if hitrate > topHit.Hitrate {
			topHit = metadataHit{Metadata: metadata, Hitrate: hitrate}
		}

		logEntry[index] = fmt.Sprintf("%s relevance: %f", metadata.Target, hitrate)
	}

	log.Debugf("%v", logEntry)

	return topHit
}

func gatherMetadataForTemplateEntities(log logger, fs *afero.Afero, templateDir string) ([]metadata, error) {
	allMetadatas := make([]metadata, 0)

	err := fs.Walk(templateDir, func(targetPath string, info iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != "metadata" {
			return nil
		}

		rawMetadata, err := fs.ReadFile(targetPath)
		if err != nil {
			return fmt.Errorf("reading target %s: %w", path.Join(templateDir, targetPath), err)
		}

		var metadatas []metadata

		err = json.Unmarshal(rawMetadata, &metadatas)
		if err != nil {
			return fmt.Errorf("unmarshalling %s: %w", path.Join(templateDir, targetPath), err)
		}

		allMetadatas = append(allMetadatas, enrichMetadatas(metadatas, path.Dir(targetPath))...)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", templateDir, err)
	}

	return allMetadatas, nil
}

func enrichMetadatas(items []metadata, baseDir string) []metadata {
	enrichedMetas := make([]metadata, len(items))

	for index, metadata := range items {
		enrichedMetas[index] = metadata

		enrichedMetas[index].Target = path.Clean(path.Join(baseDir, metadata.Target))
	}

	return enrichedMetas
}

func calculateHitrate(a map[string]string, b map[string]string) float32 {
	var hitrate float32

	for key, value := range a {
		if b[key] == value {
			hitrate++
		}
	}

	return hitrate / float32(len(a))
}
