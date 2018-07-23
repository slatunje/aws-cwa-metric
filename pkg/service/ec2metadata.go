// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
)

// EC2MetaData retrieve ec2 metadata
type EC2MetaData struct {
	Config aws.Config
	Meta   *ec2metadata.EC2Metadata
}

// NewEC2MetaData returns an instance of `EC2MetaData`
func NewEC2MetaData(cfg aws.Config) EC2MetaData {
	return EC2MetaData{Config: cfg, Meta: ec2metadata.New(cfg)}
}

// InstanceID return the instance id from meta data
func (e *EC2MetaData) InstanceID() (string, error) {
	return e.Meta.GetMetadata("instance-id")
}
