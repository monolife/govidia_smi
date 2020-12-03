module ducao/govidia_smi

go 1.15

replace ducao/govidia_smi/proto => ./proto

require (
	ducao/govidia_smi/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.33.2
	gopkg.in/yaml.v2 v2.4.0
)