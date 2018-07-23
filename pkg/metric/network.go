// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metrics/pkg/service"
	"github.com/shirou/gopsutil/net"
)

// Network metric entity
type Network struct{}

// Collect Network Traffic metrics
func (c Network) Collect(instanceID string, cw service.CloudWatch, namespace string) {
	metrics, err := net.IOCounters(false)
	if err != nil {
		log.Fatal(err)
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(name, value, unit, dime), namespace)
	}

	for _, ioc := range metrics {

		dimensions := make([]cloudwatch.Dimension, 0)

		key1 := "InstanceId"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &key1,
			Value: &instanceID,
		})
		key2 := "IOCounter"
		dimensions = append(dimensions, cloudwatch.Dimension{
			Name:  &key2,
			Value: &ioc.Name,
		})

		publish("BytesIn", float64(ioc.BytesRecv), cloudwatch.StandardUnitBytes, dimensions)
		publish("BytesOut", float64(ioc.BytesSent), cloudwatch.StandardUnitBytes, dimensions)
		publish("PacketsIn", float64(ioc.PacketsRecv), cloudwatch.StandardUnitBytes, dimensions)
		publish("PacketsOut", float64(ioc.PacketsSent), cloudwatch.StandardUnitBytes, dimensions)
		publish("ErrorsIn", float64(ioc.Errin), cloudwatch.StandardUnitBytes, dimensions)
		publish("ErrorsOut", float64(ioc.Errout), cloudwatch.StandardUnitBytes, dimensions)

		log.Printf("network - %s bytes in/out: %v/%v packets in/out: %v/%v errors in/out: %v/%v\n",
			ioc.Name, ioc.BytesRecv, ioc.BytesSent, ioc.Errin,ioc.Errout, ioc.PacketsRecv, ioc.PacketsSent,
		)
	}
}
