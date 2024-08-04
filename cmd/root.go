package cmd

import (
	"fmt"
	"github.com/maxihafer/pretix-apprise-daemon/pkg/daemon"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pretix-apprise-daemon",
	Short: "A simple daemon mapping pretix webhooks to apprise api notifications.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := &daemon.Config{}

		if err := viper.Unmarshal(config); err != nil {
			return err
		}

		server := daemon.NewServer(config)

		return server.Run()

	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("apprise-host") {
			return fmt.Errorf("apprise host not set")
		}
		if !viper.IsSet("apprise-key") {
			return fmt.Errorf("apprise key not set")
		}

		return nil
	},
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

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.pretix-apprise-daemon.yaml)")

	rootCmd.Flags().String("pretix-host", "", "Pretix API Host (PRETIX_HOST)")
	viper.BindPFlag("pretix-host", rootCmd.Flags().Lookup("pretix-host"))

	rootCmd.Flags().String("pretix-token", "", "Pretix API Token (PRETIX_TOKEN)")
	viper.BindPFlag("pretix-token", rootCmd.Flags().Lookup("pretix-token"))

	rootCmd.Flags().Int("port", 8080, "Port the daemon should listen on for webhook calls (DAEMON_PORT)")
	viper.BindPFlag("daemon-port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().String("bind", "0.0.0.0", "Address the daemon should bind to (DAEMON_BIND_ADDRESS)")
	viper.BindPFlag("daemon-bind-address", rootCmd.Flags().Lookup("bind"))

	rootCmd.Flags().String("host", "http://localhost:8000", "Apprise API Host (APPRISE_HOST)")
	viper.BindPFlag("apprise-host", rootCmd.Flags().Lookup("host"))

	rootCmd.Flags().StringP("key", "k", "", "Apprise API Configuration Key (APPRISE_KEY)")
	viper.BindPFlag("apprise-key", rootCmd.Flags().Lookup("key"))
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pretix-apprise-daemon")
	}

	viper.BindEnv("pretix-host", "PRETIX_HOST")
	viper.BindEnv("pretix-token", "PRETIX_TOKEN")
	viper.BindEnv("daemon-port", "DAEMON_PORT")
	viper.BindEnv("daemon-bind-address", "DAEMON_BIND_ADDRESS")
	viper.BindEnv("apprise-host", "APPRISE_HOST")
	viper.BindEnv("apprise-key", "APPRISE_KEY")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
