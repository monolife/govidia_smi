package main

//Standard Library
import (
    "bytes"
    "context"
    "errors"
    "fmt"
    "log"
    "net"
    "os"
    "os/exec"
    "strings"
    "strconv"
)

//Third Party
import (
	"gopkg.in/yaml.v2"

	"google.golang.org/grpc"
	pb "ducao/govidia_smi/proto"
)

type Configuration struct{
	AgentPort 			int 		`yaml:"agentPort"`// gRPC recieving port for Ingest process server
}
var _config = Configuration{}

func GetHostname() (hostname string, err error) {
	cmd := exec.Command("hostname");
	var out bytes.Buffer;
  cmd.Stdout = &out;

  if err = cmd.Run(); err != nil {
      log.Fatal(err);
      return "", err;
  }

  hostname = out.String();
  hostname = strings.Trim(hostname,"\n");
  return hostname, nil;
}

func GetGpuCount() (count int, err error) {
	cmd := exec.Command("nvidia-smi", 
  	"--query-gpu=count",
  	"--format=csv,noheader");
	var out bytes.Buffer;
  cmd.Stdout = &out;

  if err = cmd.Run(); err != nil {
      log.Fatal(err);
      return -1, err;
  }

  // countStr := out.String();
  // countStr = strings.Split(countStr,"\n");
  splitStr := strings.Split(out.String(),"\n");
  count,_ = strconv.Atoi(splitStr[0])
  return count, nil;
}

func QuerySmi(index int) (gpuInfo pb.GpuInfo, err error) {
	cmd := exec.Command("nvidia-smi", 
		"--id=" + strconv.Itoa(index),
  	"--query-gpu=timestamp,name,pci.bus_id,driver_version,pstate,pcie.link.gen.max,pcie.link.gen.current,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,index",
  	"--format=csv,noheader");

	var out bytes.Buffer;
  cmd.Stdout = &out;
  
  if err = cmd.Run(); err != nil {
		err = errors.New("Cannot get device info with index " + strconv.Itoa(index));
    log.Fatal(err);
    return pb.GpuInfo{}, err;
  }

  outStr := strings.Trim(out.String(),"\n");
  splitStr := strings.Split(outStr,", ");
  dex,_ := strconv.Atoi(splitStr[13]);
  // dex32 := 
  gpuInfo = pb.GpuInfo{
  	Timestamp         : splitStr[0],
    GpuIndex          : int32(dex),
  	Name              : splitStr[1],
  	PciBusId          : splitStr[2],
  	DriverVersion     : splitStr[3],
  	Pstate            : splitStr[4],
  	PcieLinkGenMax    : splitStr[5],
  	PcieLinkGenCurrent: splitStr[6],
  	TemperatureGpu    : splitStr[7],
  	UtilizationGpu    : splitStr[8][:len(splitStr[8])-2], //trimming " %"
  	UtilizationMemory : splitStr[9][:len(splitStr[9])-2], //trimming " %"
  	MemoryTotal       : splitStr[10][:len(splitStr[10])-4], //trimming " MiB"
  	MemoryFree        : splitStr[11][:len(splitStr[11])-4], //trimming " MiB"
  	MemoryUsed        : splitStr[12][:len(splitStr[12])-4], //trimming " MiB"
	}
  return;
}

/**
 * server is used to implement transcoder server.
 */
type server struct {
  pb.UnimplementedNvidiaQueryServiceServer}

func (s *server) QueryGpu(ctx context.Context, gpuTarget *pb.DataCue) (*pb.NvidiaQueryResponse, error){
  log.Printf("Received: %v\n", gpuTarget.GetGpuIndex())
  // TODO: just queries every GPU now

  //------------------ Print info --------------------------------
  hostname,_ := GetHostname();
  gpuCount,_ := GetGpuCount();

  gpuInfos := []*pb.GpuInfo{};

  fmt.Printf("Querying %d GPUs\n", gpuCount);
  //TODO: The better way is to just use the block output of nvidia-smi, split by "\n"
  for i:=0; i < gpuCount; i++{
    gpuInfo,err := QuerySmi(i);
    if err != nil{
      log.Fatal(err);
    }
    fmt.Println(gpuInfo);
    gpuInfos = append(gpuInfos, &gpuInfo);
  }
  //------------- (end) Print info -------------------------------

  return &pb.NvidiaQueryResponse{
    Hostname: hostname,
    Infos: gpuInfos,
  }, nil
}
/**/

func main() {

  /*------------------ Load config -------------------------------*/
  file, err := os.Open("./config.yaml")
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

  /*--------------- GRPC Server ----------------*/
  port := ":"+strconv.Itoa(_config.AgentPort);
  lis, err := net.Listen("tcp", port);
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  s := grpc.NewServer()
  pb.RegisterNvidiaQueryServiceServer(s, &server{})
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
  /*---------------(end) GRPC Server -----------*/

  /*------------------ Print info --------------------------------*
  // hostname,_ := GetHostname();
  gpuCount,_ := GetGpuCount();

  gpuInfos := []*pb.GpuInfo{};

  fmt.Printf("Querying %d GPUs\n", gpuCount);
  for i:=0; i < gpuCount; i++{
    gpuInfo,err := QuerySmi(i);
    if err != nil{
      log.Fatal(err);
    }
    fmt.Println(gpuInfo);
    gpuInfos = append(gpuInfos, &gpuInfo);
  }
  /*------------- (end) Print info -------------------------------*/
}
