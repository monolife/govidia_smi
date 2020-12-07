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

	"github.com/gorilla/mux"
)

type Configuration struct{
	MonitorUiPort	int			`yaml:"monitorUiPort"`// gRPC recieving port for Ingest process server
	AgentPort		int			`yaml:"agentPort"`// gRPC recieving port for Ingest process server
	AgentHosts		[]string 	`yaml:"agentHosts,flow"`// gRPC recieving port for Ingest process server
}
var _config = Configuration{}

var _count = 0

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
		log.Println(err)
		return
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
		log.Println(err)
		return
	}
	if( response.Infos[0].GpuIndex == 0){
		response.Infos[0].GpuIndex = 0;
	}
	// log.Info("Ingest of ID ", r.GetId()," complete")
	return;
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println(r)
    fmt.Fprintf(w, "Hi there, %s [%v]!\n\n", r.URL.Path[1:], _count)
    
    for _,hostname := range _config.AgentHosts{
	    response, err := QueryGpus(hostname)
	    if err != nil {
			log.Println(err)
	    	fmt.Fprintf(w, "No response from %s\n", hostname)
			break
		}


		marshalOpts := pjson.MarshalOptions{
			Multiline:true,
			EmitUnpopulated:true,
		}
		jsonByte,err := marshalOpts.Marshal(response)
		log.Println(_count)
		log.Println(string(jsonByte))
	    fmt.Fprintf(w, "%s\n", string(jsonByte))
	    _count++
	}

	jsonByte,err := marshalOpts.Marshal(response)

  fmt.Fprintf(w, "%s\n", string(jsonByte))
}

func handleOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	response, err := QueryGpus(params["host"])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&pb.NvidiaQueryResponse{})
	}else{

		marshalOpts := pjson.MarshalOptions{
			Multiline:true,
			EmitUnpopulated:true,
		}
		jsonByte,err := marshalOpts.Marshal(response)
		// fmt.Fprintf(w, "%s\n", string(jsonByte))
		json.NewEncoder(w).Encode(response)

	}

}
func doNothing(w http.ResponseWriter, r *http.Request){}

func main(){
	/*------------------ Load config -------------------------------*/
	configName := "./config.yaml"
	if( len(os.Args) > 1 ){
		configName = os.Args[1];
	}
	file, err := os.Open(configName);
	defer file.Close();
	if err != nil {
		log.Println("!!! Failed to open config file %v (Using hardcoded defaults) !!!", err)
		_config.MonitorUiPort = 8080;
		_config.AgentHosts = []string{"localhost"};
		_config.AgentPort = 1234;
	} else {
		decoder := yaml.NewDecoder(file);
		err = decoder.Decode(&_config);
		if err != nil {
			log.Printf("Failed to decode config: %v", err);
			return;
		}
	}
	log.Println("===========");
	log.Println("Config is: ");
	log.Println(_config);
	log.Println("===========");
	/*------------------(end) Load config --------------------------*/

<<<<<<< HEAD
	uiPort := ":"+strconv.Itoa(_config.MonitorUiPort);
	http.HandleFunc("/", handler);
	http.HandleFunc("/favicon.ico", doNothing)
    log.Println(http.ListenAndServe(uiPort, nil));
=======

	// if _,err = QueryGpus(_config.AgentHost); err != nil {
	// 	log.Fatal(err)
	// }

  router := mux.NewRouter();
	router.HandleFunc("/", handler)
  router.HandleFunc("/query", handler).Methods("GET");
  router.HandleFunc("/query/{host}", handleOne).Methods("GET")

  log.Fatal(http.ListenAndServe(":8080", router))
>>>>>>> ef9608bbd8162b8654d34dcbbd389a6ce9c75f60
}