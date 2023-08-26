package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fxpixel",
	Short: "REST service for RGB(W) LED lighting effects",
	Long:  `REST service for RGB(W) LED lighting effects`,
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/fxpixel.yaml)")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Err(err).Msg("Couldn't find home directory")
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(home)      // adding home directory as first search path
		viper.AddConfigPath(".")       // also look in the working directory
		viper.SetConfigName("fxpixel") // name the config file (without extension)
	}

	viper.AutomaticEnv() // read in environment variables that match

	//	Set our defaults
	viper.SetDefault("datastore.system", path.Join(home, "fxpixel", "db", "fxpixel.db"))
	viper.SetDefault("server.port", "3050")

	// If a config file is found, read it in
	viper.ReadInConfig()

}
