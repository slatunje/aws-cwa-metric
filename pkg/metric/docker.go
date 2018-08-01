// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package metric

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/shirou/gopsutil/docker"
	"github.com/slatunje/aws-cwa-metric/pkg/service"
)

// https://github.com/shirou/gopsutil/blob/master/docker/docker.go
const (
	DockerContainerMemory    = "docker_container_mem"
	DockerContainerCPUUser   = "docker_container_cpu_user"
	DockerContainerCPUSystem = "docker_container_cpu_system"
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
func (c Docker) Collect(doc ec2metadata.EC2InstanceIdentityDocument, cw service.CloudWatch, namespace string) {
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

		key1 := "InstanceId"
		key2 := "ImageId"
		key3 := "InstanceType"
		key4 := "ContainerId"
		key5 := "ContainerName"
		key6 := "DockerImage"

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
				Value: &container.ContainerID,
			},
			{
				Name:  &key5,
				Value: &container.Name,
			},
			{
				Name:  &key6,
				Value: &container.Image,
			},
		}

		mem, err := docker.CgroupMem(container.ContainerID, fmt.Sprintf("%s/mem/docker", base))
		if err != nil {
			log.Fatal(err)
		}
		cpu, err := docker.CgroupCPU(container.ContainerID, fmt.Sprintf("%s/cpuacct/docker", base))
		if err != nil {
			log.Fatal(err)
		}

		publish(DockerContainerMemory, float64(mem.MemUsageInBytes), cloudwatch.StandardUnitBytes, dime)
		publish(DockerContainerCPUUser, float64(cpu.User), cloudwatch.StandardUnitSeconds, dime)
		publish(DockerContainerCPUSystem, float64(cpu.System), cloudwatch.StandardUnitSeconds, dime)

		log.Printf("docker - container:%s memory:%v user:%v system:%v\n",
			container.Name, mem.MemMaxUsageInBytes, cpu.User, cpu.System,
		)
	}
}
