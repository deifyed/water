package water

import (
	"fmt"

	"github.com/deifyed/water/pkg/config"
	"github.com/deifyed/water/pkg/context"
	"github.com/deifyed/water/pkg/logging"
	"github.com/deifyed/water/pkg/template"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunE(fs *afero.Afero) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log := logging.GetLogger()

		if len(args) == 0 {
			return fmt.Errorf("no arguments provided. See --help for usage")
		}

		targetPath := args[0]
		templateDir := viper.GetString(config.TemplatesDirectory)

		log.Debug(map[string]string{
			"targetPath":  targetPath,
			"templateDir": templateDir,
		})

		exists, err := fs.Exists(targetPath)
		if err != nil {
			return fmt.Errorf("checking target path existence: %w", err)
		}

		if !exists {
			return fmt.Errorf("target path \"%s\" does not exist, please create it before watering", targetPath)
		}

		targetContext, err := context.GatherContext(log, fs, targetPath)
		if err != nil {
			return fmt.Errorf("gathering context: %w", err)
		}

		log.Debug(targetContext)

		template, err := template.Discover(log, fs, templateDir, targetContext)
		if err != nil {
			return fmt.Errorf("discovering template: %w", err)
		}

		if targetContext.TargetType != context.TargetTypeFile {
			log.Debugf("target path is a directory, don't know how to handle that yet")

			return nil
		}

		err = fs.WriteReader(targetPath, template)
		if err != nil {
			return fmt.Errorf("writing template: %w", err)
		}

		return nil
	}
}
