// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
	"github.com/shirou/gopsutil/mem"
)

// Memory metric entity
type Memory struct{}

// Collect Memory utilization
func (c Memory) Collect(id string, cw service.CloudWatch, namespace string) {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	key := "InstanceId"
	dimensions := []cloudwatch.Dimension{
		{
			Name:  &key,
			Value: &id,
		},
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(name, value, unit, dime), namespace)
	}

	publish("MemoryUtilization", m.UsedPercent, cloudwatch.StandardUnitPercent, dimensions)
	publish("MemoryUsed", float64(m.Used), cloudwatch.StandardUnitBytes, dimensions)
	publish("MemoryAvailable", float64(m.Available), cloudwatch.StandardUnitBytes, dimensions)

	log.Printf("memory - utilization:%v%% used:%v available:%v\n", m.UsedPercent, m.Used, m.Available)
}
