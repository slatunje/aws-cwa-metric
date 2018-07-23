// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metrics/pkg/service"
	"log"
	"github.com/slatunje/aws-cwa-metrics/pkg/client"
	"github.com/spf13/viper"
	"sort"
	"strings"
	"github.com/slatunje/aws-cwa-metrics/pkg/utils"
	"context"
	"os"
	"os/signal"
	"time"
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

// Collect metrics about enabled metric
func Collect(metrics []Gatherer, cw service.CloudWatch, namespace string) {
	id, err := client.InstanceID()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range metrics {
		m.Collect(id, cw, namespace)
	}
}

// Execute starts the service
func Execute() {

	var cm = chosen()
	var cw = service.NewCloudWatch()
	var ns = viper.GetString(utils.CWANamespaceKey)

	// handle one time execution?

	if viper.GetBool(utils.CWAOnceKey) {
		Collect(cm, cw, ns)
		return
	}

	// handle continuous execution? then trap Ctrl+C and call cancel on the context

	ctx := context.Background()
	ctx, cancel := OnSignal(ctx, os.Interrupt, os.Kill)
	defer cancel()

	forever(ctx, cm, cw, ns)
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

// forever will forever collect metrics unless interrupted
func forever(ctx context.Context, cm []Gatherer, cw service.CloudWatch, ns string) {
	var tt = time.NewTicker(time.Duration(viper.GetInt(utils.CWAIntervalKey)) * time.Minute)
	loop:
	for {
		select {
		case <-tt.C:
			Collect(cm, cw, ns)
		case <-ctx.Done():
			log.Printf("ok stopping forever task due to: %s...", ctx.Err())
			break loop
		}
	}
	log.Println("shutdwon completed.")
}
