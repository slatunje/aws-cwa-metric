// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
	"github.com/shirou/gopsutil/mem"
)

// Swap metric entity
type Swap struct{}

// Collect Swap usage
func (c Swap) Collect(instanceID string, cw service.CloudWatch, namespace string) {
	m, err := mem.SwapMemory()
	if err != nil {
		log.Fatal(err)
	}

	key := "InstanceId"
	d := []cloudwatch.Dimension{
		{
			Name:  &key,
			Value: &instanceID,
		},
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(name, value, unit, dime), namespace)
	}

	publish("SwapUtilization", m.UsedPercent, cloudwatch.StandardUnitPercent, d)
	publish("SwapUsed", float64(m.Used), cloudwatch.StandardUnitBytes, d)
	publish("SwapFree", float64(m.Free), cloudwatch.StandardUnitBytes, d)

	log.Printf("swap - utilization:%v%% used:%v free:%v\n", m.UsedPercent, m.Used, m.Free)
}
