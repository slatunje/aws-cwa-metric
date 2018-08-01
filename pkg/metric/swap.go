// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/shirou/gopsutil/mem"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://github.com/shirou/gopsutil/blob/master/mem/mem.go#L75
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/metrics-collected-by-CloudWatch-agent.html
const (
	SwapFreeMemory  = "swap_free"
	SwapUsedMemory  = "swap_used"
	SwapUsedPercent = "swap_used_percent"
	SwapTotalMemory = "swap_total"
)

// Swap metric entity
type Swap struct{}

// Collect Swap usage
func (c Swap) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
	m, err := mem.SwapMemory()
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
		cw.Publish(NewDatum(name, value, unit, dime), namespace)
	}

	publish(SwapFreeMemory, float64(m.Free), cloudwatch.StandardUnitBytes, dime)
	publish(SwapUsedMemory, float64(m.Used), cloudwatch.StandardUnitBytes, dime)
	publish(SwapUsedPercent, m.UsedPercent, cloudwatch.StandardUnitPercent, dime)
	publish(SwapTotalMemory, float64(m.Total), cloudwatch.StandardUnitBytes, dime)

	log.Printf("swap - utilization:%v%% used:%v free:%v total:%v \n", m.UsedPercent, m.Used, m.Free, m.Total)
}
