// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/net"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://github.com/shirou/gopsutil/blob/master/net/net.go#L17
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/metrics-collected-by-CloudWatch-agent.html
const (
	NetworkBytesIn   = "net_bytes_recv"
	NetworkBytesOut  = "net_bytes_sent"
	NetworkPacketIn  = "net_packets_recv"
	NetworkPacketOut = "net_packets_sent"
	NetworkErrorsIn  = "net_err_in"
	NetworkErrorsOut = "net_err_out"
	NetworkDropIn    = "net_drop_in"
	NetworkDropOut   = "net_drop_out"
)

// Network metric entity
type Network struct{}

// Collect Network Traffic metrics
func (c Network) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
	metrics, err := net.IOCounters(false)
	if err != nil {
		log.Fatal(err)
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(strcase.ToCamel(name), value, unit, dime), namespace)
	}

	for _, ioc := range metrics {

		key1 := "InstanceId"
		key2 := "ImageId"
		key3 := "InstanceType"
		key4 := "IOCounter"

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
				Value: &ioc.Name,
			},
		}

		publish(NetworkBytesIn, float64(ioc.BytesRecv), cloudwatch.StandardUnitBytes, dime)
		publish(NetworkBytesOut, float64(ioc.BytesSent), cloudwatch.StandardUnitBytes, dime)
		publish(NetworkPacketIn, float64(ioc.PacketsRecv), cloudwatch.StandardUnitCount, dime)
		publish(NetworkPacketOut, float64(ioc.PacketsSent), cloudwatch.StandardUnitCount, dime)
		publish(NetworkErrorsIn, float64(ioc.Errin), cloudwatch.StandardUnitCount, dime)
		publish(NetworkErrorsOut, float64(ioc.Errout), cloudwatch.StandardUnitCount, dime)
		publish(NetworkDropIn, float64(ioc.Dropin), cloudwatch.StandardUnitCount, dime)
		publish(NetworkDropOut, float64(ioc.Dropout), cloudwatch.StandardUnitCount, dime)

		log.Printf("network - %s bytes in/out: %v/%v packets in/out: %v/%v errors in/out: %v/%v\n",
			ioc.Name, ioc.BytesRecv, ioc.BytesSent, ioc.Errin, ioc.Errout, ioc.PacketsRecv, ioc.PacketsSent,
		)
	}
}
