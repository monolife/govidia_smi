package main

//Standard Library
import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Third Party
import (
	"gopkg.in/yaml.v2"
	pjson "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/grpc"
	pb "ducao/govidia_smi/proto"
)

type Configuration struct{
	AgentHost 			string 		`yaml:"agentHost"`// gRPC recieving port for Ingest process server
	AgentPort 			int 		`yaml:"agentPort"`// gRPC recieving port for Ingest process server
}
var _config = Configuration{}

//------------------------------------------------------------------------------

func QueryGpus(hostname string)(response *pb.NvidiaQueryResponse, err error){
	// Set up a connection to the server.
	address := hostname+":"+strconv.Itoa(_config.AgentPort)
	// log.Debug("Connecting to %v", address)
	conn, err := grpc.Dial(
		address, 
		grpc.WithInsecure(), 
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	);
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewNvidiaQueryServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// log.WithFields(logrus.Fields{
	// 					"Filename": filename,
	// 					"Id": id,
	// 				}).Info("Calling Ingest")
	response, err = client.QueryGpu(ctx, &pb.DataCue{GpuIndex: 0})
	if err != nil {
		log.Fatal(err)
	}
	if( response.Infos[0].GpuIndex == 0){
		response.Infos[0].GpuIndex = 0;
	}
	// log.Info("Ingest of ID ", r.GetId()," complete")
	return;
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, %s!\n\n", r.URL.Path[1:])
    
    response, err := QueryGpus(_config.AgentHost)
    if err != nil {
		log.Fatal(err)
	}

	marshalOpts := pjson.MarshalOptions{
		Multiline:true,
		EmitUnpopulated:true,
	}
	jsonByte,err := marshalOpts.Marshal(response)

    fmt.Fprintf(w, "%s\n", string(jsonByte))
}

func main(){
	/*------------------ Load config -------------------------------*/
	file, err := os.Open("./config.yaml");
	defer file.Close();
	if err != nil {
		log.Println("!!! Failed to open config file %v (Using defaults) !!!", err)
		_config.AgentPort = 1234;
	} else {
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&_config)
		if err != nil {
			log.Printf("Failed to decode config: %v", err)
			return
		}
	}
	log.Println("===========")
	log.Println("Config is: ")
	log.Println(_config)
	log.Println("===========")
	/*------------------(end) Load config --------------------------*/


	// if _,err = QueryGpus(_config.AgentHost); err != nil {
	// 	log.Fatal(err)
	// }

	http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}