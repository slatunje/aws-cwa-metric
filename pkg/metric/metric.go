// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
	"github.com/slatunje/aws-cwa-metric/pkg/utils"
	"github.com/spf13/viper"
)

const (
	KeyPrefix = "aws_metrics_"
)

const (
	KeyDisk    = "disk"
	KeyDocker  = "docker"
	KeyMemory  = "memory"
	KeyNetwork = "network"
	KeySwap    = "swap"
)

var registered = map[string]Gatherer{
	KeyDisk:    Disk{},
	KeyDocker:  Docker{},
	KeyMemory:  Memory{},
	KeyNetwork: Network{},
	KeySwap:    Swap{},
}

// Gatherer entity
type Gatherer interface {
	Collect(string, service.CloudWatch, string)
}

// NewDatum returns a slice of `[]cloudwatch.MetricDatum` data object
func NewDatum(
	name string,
	value float64,
	unit cloudwatch.StandardUnit,
	dimensions []cloudwatch.Dimension,
) []cloudwatch.MetricDatum {
	return []cloudwatch.MetricDatum{
		{
			MetricName: &name,
			Dimensions: dimensions,
			Unit:       unit,
			Value:      &value,
		},
	}
}

// Execute starts the service
func Execute() {

	var cf = config()
	var cm = chosen()
	var ns = viper.GetString(utils.CWANamespaceKey)

	var cw = service.NewCloudWatch(cf)
	var md = service.NewEC2MetaData(cf)

	var id, err = md.InstanceID()
	if err != nil {
		log.Fatal(err)
	}

	// handle one time execution?

	if viper.GetBool(utils.CWAOnceKey) {
		collect(cm, id, cw, ns)
		return
	}

	// handle continuous execution? then trap Ctrl+C and call cancel on the context

	ctx := context.Background()
	ctx, cancel := OnSignal(ctx, os.Interrupt, os.Kill)
	defer cancel()

	forever(ctx, id, cm, cw, ns)
}

// OnSignal will listen to signals and gracefully shutdown
func OnSignal(ctx context.Context, s ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, s...)
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			cancel()
		}
		signal.Stop(c)
	}()
	return ctx, cancel
}

// config returns an aws.Config object
func config() (cfg aws.Config) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config")
	}
	cfg.Region = viper.GetString(utils.CWARegionKey)
	return
}

// chosen returns a slice of chosen metrics
func chosen() (cm []Gatherer) {

	keys := viper.AllKeys()
	sort.Strings(keys)
	settings := viper.AllSettings()

	for _, k := range keys {
		if !strings.HasPrefix(k, KeyPrefix) || settings[k] == false {
			continue
		}
		if val, ok := registered[strings.TrimPrefix(k, KeyPrefix)]; ok {
			cm = append(cm, val)
			log.Printf("selected %v: %+v", k, val)
		}
	}

	return
}

// collect enabled metrics
func collect(metrics []Gatherer, id string, cw service.CloudWatch, namespace string) {
	for _, m := range metrics {
		m.Collect(id, cw, namespace)
	}
}

// forever will forever collect metrics unless interrupted
func forever(ctx context.Context, id string, cm []Gatherer, cw service.CloudWatch, ns string) {
	var tt = time.NewTicker(time.Duration(viper.GetInt(utils.CWAIntervalKey)) * time.Minute)
	loop:
	for {
		select {
		case <-tt.C:
			collect(cm, id, cw, ns)
		case <-ctx.Done():
			log.Printf("ok stopping forever task due to: %s...", ctx.Err())
			break loop
		}
	}
	log.Println("shutdwon completed.")
}
