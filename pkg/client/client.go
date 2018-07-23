// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package client

import (
	"net/http"
	"io/ioutil"
	"io"
	"time"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"fmt"
)

// Application default settings
const (
	AppContentTypeKey   = "Content-Type"
	AppContentTypeValue = "application/json"
	AppTraceKey         = "X-Trace-ID"
	AppHTTPTimeOut      = 30
)

// New creates and returns an instances of `*http.Client` and
// `*http.Request` for the purpose of making a requests to retrieve EC2 instance data
func New(method string, url string, body io.Reader, trace string) (cli *http.Client, req *http.Request, err error) {
	cli = &http.Client{Timeout: time.Duration(AppHTTPTimeOut * time.Second)}
	if trace == "" {
		trace = uuid.NewV4().String()
	}
	req, err = http.NewRequest(method, url, body)
	req.Header.Set(AppContentTypeKey, AppContentTypeValue)
	req.Header.Set(AppTraceKey, trace)
	if err != nil {
		return
	}
	return
}

// InstanceID return EC2 instance id
func InstanceID() (string, error) {
	value := viper.GetString("aws_instance_id")
	if len(value) > 0 {
		return value, nil
	}
	return metadata("instance-id")
}

// ImageID returns EC2 image id
func ImageID() (string, error) {
	value := viper.GetString("aws_image_id")
	if len(value) > 0 {
		return value, nil
	}
	return metadata("ami-id")
}

// InstanceType returns the type of instance used
func InsatnceType() (string, error) {
	value := viper.GetString("aws_instance_type")
	if len(value) > 0 {
		return value, nil
	}
	return metadata("ami-id")
}

func metadata(path string) (string, error) {
	url := fmt.Sprintf("http://169.254.169.254/latest/meta-data/%s", path)
	cli, req, err := New("GET", url, nil, "")
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}


