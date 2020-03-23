module ExifParser

go 1.14

require (
	cloud.google.com/go/pubsub v1.3.1
	cloud.google.com/go/storage v1.6.0
	github.com/barasher/go-exiftool v1.1.1
	github.com/barsanuphe/goexiftool v0.0.0-20180509224600-0e25a2871ba9
	github.com/google/uuid v1.1.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/stretchr/testify v1.4.0
	mikenimer.com/services/core/GcpUtils v1.0.0
)

replace mikenimer.com/services/core/GcpUtils => ../../core/GcpUtils
