// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/mem"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://github.com/shirou/gopsutil/blob/master/mem/mem.go#L15
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/metrics-collected-by-CloudWatch-agent.html
const (
	MemoryTotal       = "mem_total"
	MemoryAvailable   = "mem_available"
	MemoryUsed        = "mem_used"
	MemoryUsedPercent = "mem_used_percent"
	MemoryFree        = "mem_free"
	MemoryCached      = "mem_cached"
)

// Memory metric entity
type Memory struct{}

// Collect Memory utilization
func (c Memory) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	key1 := "InstanceId"
	key2 := "ImageId"
	key3 := "InstanceType"
	dime := []cloudwatch.Dimension{
		{
			Name:  &key1,
			Value: &doc.InstanceID,
		},
		{
			Name:  &key2,
			Value: &doc.ImageID,
		},
		{
			Name:  &key3,
			Value: &doc.InstanceType,
		},
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(strcase.ToCamel(name), value, unit, dime), namespace)
	}

	publish(MemoryTotal, float64(m.Total), cloudwatch.StandardUnitBytes, dime)
	publish(MemoryAvailable, float64(m.Available), cloudwatch.StandardUnitBytes, dime)
	publish(MemoryUsed, float64(m.Used), cloudwatch.StandardUnitBytes, dime)
	publish(MemoryUsedPercent, m.UsedPercent, cloudwatch.StandardUnitPercent, dime)
	publish(MemoryFree, float64(m.Free), cloudwatch.StandardUnitBytes, dime)
	publish(MemoryCached, float64(m.Cached), cloudwatch.StandardUnitBytes, dime)

	log.Printf("memory - utilization:%v%% used:%v available:%v\n", m.UsedPercent, m.Used, m.Available)
}
