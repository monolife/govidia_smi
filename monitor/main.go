package main

//Standard Library
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Third Party
import (
	"gopkg.in/yaml.v2"
	
	"google.golang.org/grpc"
	pb "ducao/govidia_smi/proto"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type Configuration struct{
	MonitorPort	int			`yaml:"monitorPort"`// gRPC recieving port for Ingest process server
	AgentPort	int			`yaml:"agentPort"`// gRPC recieving port for Ingest process server
	AgentHosts	[]string 	`yaml:"agentHosts,flow"`// gRPC recieving port for Ingest process server
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
		log.Println(err)
		return
	}
	defer conn.Close()
	client := pb.NewNvidiaQueryServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err = client.QueryGpu(ctx, &pb.DataCue{GpuIndex: 0})
	if err != nil {
		log.Println(err)
		return
	}
	if( response.Infos[0].GpuIndex == 0){
		response.Infos[0].GpuIndex = 0;
	}

	return;
}

func queryRound(w http.ResponseWriter, r *http.Request) {

	var resps []*pb.NvidiaQueryResponse

	for _,hostname := range _config.AgentHosts{
		response, err := QueryGpus(hostname)
		if err != nil {
			log.Println(err)
			break
		}
		resps = append(resps, response)
	}
	json.NewEncoder(w).Encode(resps)

}

func handleOne(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	log.Println(params)

	response, err := QueryGpus(params["host"])
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(&pb.NvidiaQueryResponse{})
	}else{
		json.NewEncoder(w).Encode(response)
	}

}

func doNothing(w http.ResponseWriter, r *http.Request){
	//NO-OP function to deal with things like favicon fetching
}

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
		_config.MonitorPort = 8080;
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
	log.Printf("%+v\n",_config);
	log.Println("===========");
	/*------------------(end) Load config --------------------------*/

	monitorPort := ":"+strconv.Itoa(_config.MonitorPort);

  router := mux.NewRouter();
  router.HandleFunc("/", doNothing)
  router.HandleFunc("/favicon.ico", doNothing)
  router.HandleFunc("/query", queryRound).Methods("GET");
  // router.HandleFunc("/query/{host}", handleOne).Methods("GET")

  log.Fatal( http.ListenAndServe(monitorPort, 
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), 
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), 
			handlers.AllowedOrigins([]string{"*"}))(router)))
}