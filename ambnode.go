package ambutils

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/managedblockchain"
)

// Create Peer Node for a member in the network
// In AWS, it should be EC2 instance
func (ambNode *NetworkConfig) CreateNode(networkId string, memberId string) {

	var createNodeInput managedblockchain.CreateNodeInput
	createNodeInput.SetMemberId(memberId)
	createNodeInput.SetNetworkId(networkId)

	var nodeConfig managedblockchain.NodeConfiguration
	nodeConfig.SetAvailabilityZone("us-east-1a")
	nodeConfig.SetInstanceType("bc.t3.medium")
	createNodeInput.SetNodeConfiguration(&nodeConfig)

	nodeOut, err := ambNode.Amb.CreateNode(&createNodeInput)
	if err != nil {
		log.Fatalf("Unable to CreateNode : %s\n", err)
	}
	fmt.Println(nodeOut)
}
