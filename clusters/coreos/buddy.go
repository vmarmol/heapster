// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	heapster "github.com/GoogleCloudPlatform/heapster/sources"
	fleetClient "github.com/coreos/fleet/client"
	"github.com/coreos/fleet/etcd"
	fleetPkg "github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/registry"
)

var argCadvisorPort = flag.Int("cadvisor_port", 4194, "cAdvisor port in current CoreOS cluster")
var argEndpoints = flag.String("endpoints", "http://127.0.0.1:4001", "Comma separated list of fleet server endpoints")

func getFleetRegistryClient() (fleetClient.API, error) {
	var dial func(string, string) (net.Conn, error)

	tlsConfig, err := fleetPkg.ReadTLSConfigFiles("", "", "")
	if err != nil {
		return nil, err
	}

	trans := &http.Transport{
		Dial:            dial,
		TLSClientConfig: tlsConfig,
	}

	timeout := 3 * 1000 * time.Millisecond

	machines := strings.Split(*argEndpoints, ",")
	eClient, err := etcd.NewClient(machines, trans, timeout)
	if err != nil {
		return nil, err
	}

	reg := registry.NewEtcdRegistry(eClient, "/_coreos.com/fleet/")

	return &fleetClient.RegistryClient{Registry: reg}, nil
}

func getMachines(client fleetClient.API, outMachines map[string]string) error {
	machines, err := client.Machines()
	if err != nil {
		return err
	}
	for _, machine := range machines {
		outMachines[machine.ID] = machine.PublicIP
	}
	return nil
}

func updateHeapsterHostsFile(outMachines map[string]string) error {
	hosts := &heapster.CadvisorHosts{
		Port:  *argCadvisorPort,
		Hosts: outMachines,
	}
	data, err := json.Marshal(hosts)
	if err != nil {
		return err
	}
	fmt.Printf("writing %s to hosts file\n", string(data))
	if err := ioutil.WriteFile(heapster.HostsFile, data, 0755); err != nil {
		return err
	}
	return nil
}

func doWork(client fleetClient.API) error {
	machines := make(map[string]string)
	err := getMachines(client, machines)
	if err != nil {
		return err
	}
	return updateHeapsterHostsFile(machines)

}

func main() {
	client, err := getFleetRegistryClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err = doWork(client)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
