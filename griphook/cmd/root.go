package cmd

import (
	"fmt"
	"os"

	"github.com/cryptohoard/experimental/griphook/cmd/daemon"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "griphook",
	Short: "griphook is your personal crypto portfolio assistant.",
	Long: `griphook is your personal crypto portfolio assistant.
It will help you with following jobs:
    - view transactions
    - buy/sell coins
    - submit stop/limit buy/sell orders
    - reports on how your portfolio is performing
    - price alerts
Currently it supports Coinbase and GDAX.`,
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.griphook.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("daemon", "d", false, "Run as a service.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".griphook" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".griphook")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	if !cmd.Flags().Changed("daemon") {
		return
	}
	fmt.Println("run as daemon")
	daemon.Run()
}
