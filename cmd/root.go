package cmd

import (
	"os"
	"path"

	"github.com/deifyed/water/cmd/water"
	"github.com/deifyed/water/pkg/config"
	"github.com/deifyed/water/pkg/logging"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	logLevel     string
	templatesDir string

	fs = &afero.Afero{Fs: afero.NewOsFs()}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "water <path to file or directory>",
	Short:        "Scaffold files and directories based on names and context",
	Example:      "touch Makefile && water Makefile\nmkdir go-project && water go-project",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE:         water.RunE(fs),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}

	cfgDir := path.Join(home, ".config", "water")
	templatesDir = path.Join(cfgDir, "templates")

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"set config file (default $HOME/.config/water/water.yaml)",
	)

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "set log level")

	err = viper.BindPFlag(config.LogLevel, rootCmd.PersistentFlags().Lookup("log-level"))
	cobra.CheckErr(err)

	rootCmd.Flags().StringVarP(&templatesDir, "templates", "t", templatesDir, "set templates directory")

	err = viper.BindPFlag(config.TemplatesDirectory, rootCmd.Flags().Lookup("templates"))
	cobra.CheckErr(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log := logging.GetLogger()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgDir := path.Join(home, ".config", "water")

		viper.AddConfigPath(home)
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("water")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Debug("No config file found, using defaults")
	}
}
