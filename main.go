package main

import (
	"ambutils"
	"fmt"
)

func main() {

	amb := ambutils.NetworkConfig{
		Region:        "us-east-1",
		FabricVersion: "1.2",
		NetworkName:   "ambnetwork",
		OrgList:       []string{"supplier", "manufacturer"},
		PeersPerOrg:   2,
	}
	// Create Network
	fmt.Println("Creating Network....")
	networkOut := amb.CreateNetwork()
	fmt.Println(networkOut)
	fmt.Println("CreateNetwork DONE")

}
