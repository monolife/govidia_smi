module ducao/govidia_smi/monitor

go 1.15

replace ducao/govidia_smi/proto => ../proto

require (
	ducao/govidia_smi/proto v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
)
