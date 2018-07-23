// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metrics/pkg/service"
	"github.com/shirou/gopsutil/disk"
)

// Disk metric entity
type Disk struct{}

// Collect Disk used & free space
func (c Disk) Collect(id string, cw service.CloudWatch, namespace string) {
	m, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	k := "InstanceId"
	d := []cloudwatch.Dimension{
		{
			Name:  &k,
			Value: &id,
		},
	}

	cw.Publish(
		NewDatum("DiskUtilization", m.UsedPercent, cloudwatch.StandardUnitPercent, d), namespace)
	cw.Publish(
		NewDatum("DiskUsed", float64(m.Used), cloudwatch.StandardUnitBytes, d), namespace)
	cw.Publish(
		NewDatum("DiskFree", float64(m.Free), cloudwatch.StandardUnitBytes, d), namespace)

	log.Printf("disk - utilization:%v%% used:%v free:%v\n", m.UsedPercent, m.Used, m.Free)
}
