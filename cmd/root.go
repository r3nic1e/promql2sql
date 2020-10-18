package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "promql2sql",
	Short: "Convert prometheus metrics to PostgreSQL data",
	Run:   main,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file")

	rootCmd.Flags().String("prometheus-addr", "http://localhost:9090", "Prometheus address")
	rootCmd.Flags().String("postgres-addr", "postgres://postgres:@localhost:5432/postgres?sslmode=disable", "PostgreSQL address")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	opts := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeHookFunc("2006-01-02"),
		mapstructure.StringToTimeDurationHookFunc(),
	))
	err := viper.Unmarshal(&cfg, opts)
	if err != nil {
		panic(err)
	}
}
