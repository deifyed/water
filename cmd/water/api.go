package water

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunE(fs *afero.Afero) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no arguments provided. See --help for usage")
		}

		fmt.Println("Hello, world!")
		return nil
	}
}
