package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "course-fetcher",
		Short: "Fetch course from coursera",
		Long:  ``,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .config.json)")
	rootCmd.PersistentFlags().Bool("db_debug", false, "log sql to console")
	viper.BindPFlag("db_debug", rootCmd.PersistentFlags().Lookup("db_debug"))

	viper.SetDefault("db_network", "tcp")
	viper.SetDefault("db_addr", "localhost:5432")
	viper.SetDefault("db_user", "course-fetcher")
	viper.SetDefault("db_password", "course-fetcher")
	viper.SetDefault("db_database", "course-fetcher")

	viper.SetDefault("log_level", "debug")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
		viper.SetConfigName(".config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
