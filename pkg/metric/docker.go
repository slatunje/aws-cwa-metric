// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
	"github.com/shirou/gopsutil/docker"
)

// Docker metric entity
type Docker struct{}

// On older systems, the control groups might be mounted on /cgroup
func cGroupMountPath() (string, error) {
	out, err := exec.Command("grep", "-m1", "cgroup", "/proc/mounts").Output()
	if err != nil {
		return "", errors.New("cannot figure out where control groups are mounted")
	}
	res := strings.Fields(string(out))
	if strings.HasPrefix(res[1], "/cgroup") {
		return "/cgroup", nil
	}
	return "/sys/fs/cgroup", nil
}

// Collect CPU & Memory usage per Docker Container
func (c Docker) Collect(id string, cw service.CloudWatch, namespace string) {
	containers, err := docker.GetDockerStat()
	if err != nil {
		log.Fatal(err)
	}

	base, err := cGroupMountPath()
	if err != nil {
		log.Fatal(err)
	}

	var publish = func(name string, value float64, unit cloudwatch.StandardUnit, dime []cloudwatch.Dimension) {
		cw.Publish(NewDatum(name, value, unit, dime), namespace)
	}

	for _, container := range containers {
		d := make([]cloudwatch.Dimension, 0)

		// declare dimension keys

		key1 := "InstanceId"
		d = append(d, cloudwatch.Dimension{
			Name:  &key1,
			Value: &id,
		})
		key2 := "ContainerId"
		d = append(d, cloudwatch.Dimension{
			Name:  &key2,
			Value: &container.ContainerID,
		})
		key3 := "ContainerName"
		d = append(d, cloudwatch.Dimension{
			Name:  &key3,
			Value: &container.Name,
		})
		key4 := "DockerImage"
		d = append(d, cloudwatch.Dimension{
			Name:  &key4,
			Value: &container.Image,
		})

		// mem

		mem, err := docker.CgroupMem(container.ContainerID, fmt.Sprintf("%s/mem/docker", base))
		if err != nil {
			log.Fatal(err)
		}
		publish("ContainerMemory", float64(mem.MemUsageInBytes), cloudwatch.StandardUnitBytes, d)

		// cpu

		cpu, err := docker.CgroupCPU(container.ContainerID, fmt.Sprintf("%s/cpuacct/docker", base))
		if err != nil {
			log.Fatal(err)
		}
		publish("ContainerCPUUser", float64(cpu.User), cloudwatch.StandardUnitSeconds, d)

		// cpu system data

		publish("ContainerCPUSystem", float64(cpu.System), cloudwatch.StandardUnitSeconds, d)

		// log message

		msg := "docker - container:%s memory:%v user:%v system:%v\n"
		log.Printf(msg, container.Name, mem.MemMaxUsageInBytes, cpu.User, cpu.System)
	}
}
