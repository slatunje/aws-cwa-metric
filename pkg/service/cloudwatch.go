// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package service

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/spf13/viper"
)

// CloudWatch stores an aws configuration
type CloudWatch struct {
	Config aws.Config
}

// NewCloudWatch creates and instance of service.CloudWatch
func NewCloudWatch() CloudWatch {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config")
	}
	cfg.Region = viper.GetString("aws_default_region")
	return CloudWatch{Config: cfg}
}

// Publish saves metric data to cloud watch using AWS CloudWatch API
func (c CloudWatch) Publish(data []cloudwatch.MetricDatum, namespace string) {
	svc := cloudwatch.New(c.Config)
	req := svc.PutMetricDataRequest(&cloudwatch.PutMetricDataInput{
		MetricData: data,
		Namespace:  &namespace,
	})
	_, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}
}



