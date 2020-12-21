module ducao/govidia_smi

go 1.15

replace ducao/govidia_smi/proto => ./proto

require (
	ducao/govidia_smi/proto v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201211090839-8ad439b19e0f // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201211151036-40ec1c210f7a // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.4.0
)
