module github.com/510909033/redisx

go 1.16

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/gin-gonic/gin v1.8.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redsync/redsync/v4 v4.5.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/kafka-go v0.4.30
	github.com/sigurn/crc16 v0.0.0-20211026045750-20ab5afb07e3
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.7.1
	go.uber.org/zap v1.21.0
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
)

//replace go.etcd.io/bbolt@v1.3.6=>github.com/coreos/bbolt@v1.3.6