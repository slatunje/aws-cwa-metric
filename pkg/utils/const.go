// Copyright © 2018 Sylvester La-Tunje. All rights reserved.

package utils

// ExitXXX represents various exit code within the system
const (
	_                       = iota
	ExitExecute
	ExitRequireKeys
	ExitCredentialsFailure
	ExitCommandlineFailure
	ExitShareConfigFailure
	ExitBase64DecodeFailure
	ExitOnDebug
)

const (
	CWARegion    = "eu-west-1"
	CWANamespace = "CustomMetrics"
	CWAInterval  = 5
)

const (
	CWARegionKey    = "aws_cwa_region"
	CWANamespaceKey = "aws_cwa_namespace"
	CWAIntervalKey  = "aws_cwa_interval"
	CWAOnceKey      = "aws_cwa_once"
)
