// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/disk"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/metrics-collected-by-CloudWatch-agent.html
const (
	DiskUsedPercent   = "disk_used_percent"
	DiskUsed          = "disk_used"
	DiskFree          = "disk_free"
	DiskInodesPercent = "disk_inodes_used_percent"
	DiskInodesUsed    = "disk_inodes_used"
	DiskInodesFree    = "disk_inodes_free"
)

const (
	DiskIoIoTimes        = "diskio_io_time"
	DiskIOPsInProgress   = "diskio_iops_in_progress"
	DiskIOWrites         = "diskio_writes"
	DiskIOReads          = "diskio_reads"
	DiskWriteBytes       = "diskio_write_bytes"
	DiskReadBytes        = "diskio_read_bytes"
	DiskWriteTimes       = "diskio_write_time"
	DiskReadTimes        = "diskio_read_time"
	DiskWeightedIO       = "diskio_weighted_io"
	DiskMergedWriteCount = "diskio_merged_write"
	DiskMergedReadCount  = "diskio_merged_read"
)

const (
	PartitionDeviceCGroup  = "cgroup"
	PartitionDeviceOverlay = "overlay"
)

// Disk metric entity
type Disk struct{}

// Collect Disk used & free space
func (c Disk) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
	partitions, err := disk.Partitions(true)
	if err != nil || len(partitions) == 0 {
		log.Fatal(err)
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(strcase.ToCamel(name), value, unit, dime), namespace)
	}

	key1 := "InstanceId"
	key2 := "ImageId"
	key3 := "InstanceType"

	// handle usage

	for i, p := range partitions {

		key4 := "device"
		key5 := "fstype"
		key6 := "path"

		switch p.Device {
		case PartitionDeviceCGroup, PartitionDeviceOverlay:
			continue
		}

		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Fatal(err)
		}
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
				Value: &p.Device,
			},
			{
				Name:  &key5,
				Value: &p.Fstype,
			},
			{
				Name:  &key6,
				Value: &p.Mountpoint,
			},
		}

		publish(DiskUsedPercent, u.UsedPercent, cloudwatch.StandardUnitPercent, dime)
		publish(DiskUsed, float64(u.Used), cloudwatch.StandardUnitBytes, dime)
		publish(DiskFree, float64(u.Free), cloudwatch.StandardUnitBytes, dime)

		publish(DiskInodesPercent, float64(u.InodesUsedPercent), cloudwatch.StandardUnitPercent, dime)
		publish(DiskInodesUsed, float64(u.InodesUsed), cloudwatch.StandardUnitCount, dime)
		publish(DiskInodesFree, float64(u.InodesFree), cloudwatch.StandardUnitCount, dime)

		log.Printf("disk: [%d] device: %s,  mountpoint: %s, fstype: %s", i, p.Device, p.Mountpoint, p.Fstype)

	}

	ioc, err := disk.IOCounters()
	if err != nil {
		log.Fatal(err)
	}

	// handle counter

	for _, i := range ioc {

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
				Value: &i.Name,
			},
		}

		publish(DiskIoIoTimes, float64(i.IoTime), cloudwatch.StandardUnitMilliseconds, dime)
		publish(DiskIOPsInProgress, float64(i.IopsInProgress), cloudwatch.StandardUnitMilliseconds, dime)
		publish(DiskIOWrites, float64(i.WriteCount), cloudwatch.StandardUnitCount, dime)
		publish(DiskIOReads, float64(i.ReadCount), cloudwatch.StandardUnitCount, dime)
		publish(DiskWriteBytes, float64(i.WriteBytes), cloudwatch.StandardUnitBytes, dime)
		publish(DiskReadBytes, float64(i.ReadBytes), cloudwatch.StandardUnitBytes, dime)
		publish(DiskWriteTimes, float64(i.WriteTime), cloudwatch.StandardUnitMilliseconds, dime)
		publish(DiskReadTimes, float64(i.ReadBytes), cloudwatch.StandardUnitMilliseconds, dime)
		publish(DiskWeightedIO, float64(i.WeightedIO), cloudwatch.StandardUnitCount, dime)
		publish(DiskMergedWriteCount, float64(i.MergedWriteCount), cloudwatch.StandardUnitCount, dime)
		publish(DiskMergedReadCount, float64(i.MergedReadCount), cloudwatch.StandardUnitCount, dime)

		log.Printf("disk - %d ms bytes(read/write): %v/%v count(read/write): %v/%v\n",
			i.IoTime, i.ReadBytes, i.WriteBytes, i.ReadCount, i.WriteCount,
		)

	}

}
