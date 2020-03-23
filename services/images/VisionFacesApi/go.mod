module VisionApi

go 1.14

require (
	cloud.google.com/go v0.55.0
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8 // indirect
	golang.org/x/sys v0.0.0-20200321134203-328b4cd54aae // indirect
	google.golang.org/genproto v0.0.0-20200319113533-08878b785e9c
	mikenimer.com/services/core/GcpUtils v1.0.0
)

replace mikenimer.com/services/core/GcpUtils => ../../core/GcpUtils
