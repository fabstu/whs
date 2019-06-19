module aduu.dev

go 1.12

require (
	github.com/adlio/trello v0.0.0-20190408222002-5b31feeddcaa
	github.com/atotto/clipboard v0.1.2
	github.com/aws/aws-sdk-go v1.17.14
	github.com/bradfitz/iter v0.0.0-20190303215204-33e6a9893b0c
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/improbable-eng/grpc-web v0.9.5
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.0.0
	github.com/magiconair/properties v1.8.0
	github.com/mattn/go-sqlite3 v1.10.0
	github.com/miekg/dns v1.1.6
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/otiai10/copy v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3-0.20190313112143-fa4aa9000d28
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/common v0.2.0
	github.com/prometheus/procfs v0.0.0-20190315082738-e56f2e22fc76 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190306220146-200a235640ff // indirect
	github.com/spf13/cobra v0.0.4-0.20190311125509-ba1052d4cbce
	github.com/spf13/viper v1.3.1
	github.com/stretchr/testify v1.3.0
	github.com/tealeg/xlsx v1.0.3
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/tools v0.0.0-20190424031103-cb2dda6eabdf
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	google.golang.org/grpc v1.19.0
	gopkg.in/ini.v1 v1.42.0 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb

	k8s.io/apimachinery v0.0.0-20190311155258-f9b45bc4494d
)

//replace github.com/DataDog/dd-trace-go => gopkg.in/DataDog/dd-trace-go.v1/ddtrace
