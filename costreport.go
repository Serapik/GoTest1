package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
)

type Info struct {
	//Name       string `json:"name"`
	Properties struct {
		//	Cluster string `json:"cluster"`
		//Node      string `json:"node"`
		Namespace string `json:"namespace"`
		//NamespaceLabels struct {
		//KubernetesIoMetadataName string `json:"kubernetes_io_metadata_name"`
		//} `json:"namespaceLabels"`
	} `json:"properties"`
	Window struct {
		//	Start time.Time `json:"start"`
		//	End   time.Time `json:"end"`
	} `json:"window"`
	//Start       time.Time `json:"start"`
	//End         time.Time `json:"end"`
	//Minutes     float64   `json:"minutes"`
	CPUCost     float64 `json:"cpuCost"`
	GpuCost     float64 `json:"gpuCost"`
	NetworkCost float64 `json:"networkCost"`
	RAMCost     float64 `json:"ramCost"`
	SharedCost  float64 `json:"sharedCost"`
	TotalCost   float64 `json:"totalCost"`
}
type OpenCostMain struct {
	Code   int                `json:"code"`
	Status string             `json:"status"`
	Data   []map[string]*Info `json:"data"`
}

func main() {
	//windowsArg := flag.String("w", "30d", "the range to fetch OpenCost data in")
	//flag.Parse()
	window, windowFound := os.LookupEnv("WINDOW")
	if !windowFound {
		window = "30d"
	}
	Ip, ipFound := os.LookupEnv("URL_IP")
	if !ipFound {
		panic("IP needs to be written")
	}
	buildNumber, buildnumberFound := os.LookupEnv("BUILD_NUMBER")
	if !buildnumberFound {
		panic("Build number is required")
	}

	response, err := http.Get(fmt.Sprintf("%s/model/allocation/compute?window=%s&accumulate=false&aggregate=namespace", Ip, window))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var x OpenCostMain
	err = json.Unmarshal(bytes, &x)
	if err != nil {
		panic(err)
	}
	if x.Code != 200 {
		panic(x.Code)
	}
	if x.Data == nil || len(x.Data) == 0 {
		fmt.Println("No data to write to CSV.")
		return
	}
	Infos := make([]Info, 0, len(x.Data))
	for _, v := range x.Data[0] {
		Infos = append(Infos, *v)
	}
	// Open a file for writing

	outFile, err := os.Create("output" + buildNumber + ".csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(outFile.Name())

	defer outFile.Close() // Ensure the file gets closed
	// Write the CSV data to the file
	err = gocsv.Marshal(&Infos, outFile)
	if err != nil {
		panic(err)
	}
	// Explicitly close the file
	err = outFile.Close()
	if err != nil {
		panic(err)
	}
	// Print the raw JSON if you still want to
	//fmt.Println(string(bytes))
}
