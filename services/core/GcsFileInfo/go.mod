module github.com/mnimer/dam/services/core/GcsFileInfo

go 1.14

require (
    github.com/mnimer/dam/services/core/GcpUtils v1.0.0
	cloud.google.com/go/pubsub v1.3.1
	cloud.google.com/go/storage v1.6.0
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	golang.org/x/sys v0.0.0-20200316230553-a7d97aace0b0 // indirect
	golang.org/x/tools v0.0.0-20200316212524-3e76bee198d8 // indirect
	google.golang.org/genproto v0.0.0-20200316142031-303a05041dad // indirect
)

replace github.com/mnimer/dam/services/core/GcpUtils => ../GcpUtils