// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/cpu"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://github.com/shirou/gopsutil/blob/master/cpu/cpu.go
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/viewing_metrics_with_cloudwatch.html
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/metrics-collected-by-CloudWatch-agent.html
const (
	CPUUsageIdle   = "cpu_usage_idle"
	CPUUsageSystem = "cpu_usage_system"
	CPUUsageIOWait = "cpu_usage_iowait"
	CPUUsageUser   = "cpu_usage_user"
)

// CPU
type CPU struct{}

// Collect Swap usage
func (c CPU) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
	metrics, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}

	times, err := cpu.Times(true)
	if err != nil {
		log.Fatal(err)
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(strcase.ToCamel(name), value, unit, dime), namespace)
	}

	for _, m := range metrics {

		for _, t := range times {

			key1 := "InstanceId"
			key2 := "ImageId"
			key3 := "InstanceType"
			key4 := "cpu"
			val4 := fmt.Sprintf("cpu%d", m.CPU)
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
				{
					Name:  &key4,
					Value: &val4,
				},
			}

			publish(CPUUsageIdle, t.Idle, cloudwatch.StandardUnitPercent, dime)
			publish(CPUUsageSystem, t.System, cloudwatch.StandardUnitPercent, dime)
			publish(CPUUsageIOWait, t.Iowait, cloudwatch.StandardUnitPercent, dime)
			publish(CPUUsageUser, t.User, cloudwatch.StandardUnitPercent, dime)

			log.Printf("cpu - %s:%v%% %s:%v %s:%v %s:%v \n",
				CPUUsageIdle, t.Idle, CPUUsageSystem, t.System, CPUUsageIOWait, t.Iowait, CPUUsageUser, t.User,
			)

		}

	}

}
