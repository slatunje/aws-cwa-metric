// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/slatunje/aws-cwa-metric/pkg/metric"
	"github.com/slatunje/aws-cwa-metric/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	app     = "cwametric"
	version = "0.0.1"
)

var (
	region    string
	namespace string
	interval  int
	once      bool
	memory    bool
	swap      bool
	cpu       bool
	disk      bool
	network   bool
	docker    bool
)

// rootCmd represents the base command when called without any sub commands
var rootCmd = &cobra.Command{
	Use:   app,
	Short: fmt.Sprintf("=> %s sends aws metric to cloud watch", app),
	Long: fmt.Sprintf(`
Description:
  %s sends aws metric to cloud watch.
`, app),
	Run: func(cmd *cobra.Command, args []string) {
		metric.Execute()
	},
}

// init is called in alphabetic order within this package
func init() {
	setEmptyTimezone()
	cobra.OnInitialize(initConfig)
	rootCmd.Version = version
	// === settings === //
	rootCmd.PersistentFlags().
		StringVar(&region, "region", utils.CWARegion, "set aws region value.")
	rootCmd.PersistentFlags().
		StringVar(&namespace, "namespace", utils.CWANamespace, "set metric label.")
	rootCmd.PersistentFlags().
		IntVarP(&interval, "interval", "i", utils.CWAInterval, "set time interval value.")
	rootCmd.PersistentFlags().
		BoolVarP(&once, "once", "o", false, "execute once and stop. (i.e. never repeat.")
	// === metrics === //
	rootCmd.PersistentFlags().
		BoolVarP(&disk, metric.KeyCPU, "c", false, "collect cpu metrics.")
	rootCmd.PersistentFlags().
		BoolVarP(&disk, metric.KeyDisk, "d", false, "collect disk metrics.")
	rootCmd.PersistentFlags().
		BoolVar(&docker, metric.KeyDocker, false, "collect docker container metrics.")
	rootCmd.PersistentFlags().
		BoolVarP(&memory, metric.KeyMemory, "m", false, "collect memory metrics.")
	rootCmd.PersistentFlags().
		BoolVarP(&network, metric.KeyNetwork, "n", false, "collect network metrics.")
	rootCmd.PersistentFlags().
		BoolVarP(&swap, metric.KeySwap, "s", false, "collect swap metrics.")
}

// setEmptyTimezone
func setEmptyTimezone() {
	if err := os.Setenv("TZ", ""); err != nil {
		log.Fatalln(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setDefaults()
}

// setDefaults
func setDefaults() {
	viper.AutomaticEnv()
	viper.SetDefault(utils.CWARegionKey, region)
	viper.SetDefault(utils.CWANamespaceKey, namespace)
	viper.SetDefault(utils.CWAIntervalKey, interval)
	viper.SetDefault(utils.CWAOnceKey, once)
	viper.SetDefault("aws_metrics_cpu", cpu)
	viper.SetDefault("aws_metrics_memory", memory)
	viper.SetDefault("aws_metrics_swap", swap)
	viper.SetDefault("aws_metrics_disk", disk)
	viper.SetDefault("aws_metrics_network", network)
	viper.SetDefault("aws_metrics_docker", docker)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(utils.ExitExecute)
	}
}
