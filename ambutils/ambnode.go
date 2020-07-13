package ambutils

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/managedblockchain"
)

// Create Peer Node for a member in the network
// In AWS, it should be EC2 instance
func (ambNode *NetworkConfig) CreateNode(networkId string, memberId string) string {

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

	return *nodeOut.NodeId
}

func (ambNode *NetworkConfig) DeleteNode(networkId string, memberId string, nodeId string) {
	var deleteNodeInput managedblockchain.DeleteNodeInput

	deleteNodeInput.SetNetworkId(networkId)
	deleteNodeInput.SetMemberId(memberId)
	deleteNodeInput.SetNodeId(nodeId)

	_, err := ambNode.Amb.DeleteNode(&deleteNodeInput)
	if err != nil {
		log.Fatalf("Unable to Delete Node (%s) : %s\n", nodeId, err)
	}
	fmt.Printf("Deleted Node Successfully NodeId : %s\n", nodeId)
}

func (ambNode *NetworkConfig) GetNode(networkId string, memberId string, nodeId string) *managedblockchain.Node {
	var getNodeInput managedblockchain.GetNodeInput

	getNodeInput.SetNetworkId(networkId)
	getNodeInput.SetMemberId(memberId)
	getNodeInput.SetNodeId(nodeId)

	getNodeOutput, err := ambNode.Amb.GetNode(&getNodeInput)
	if err != nil {
		log.Printf("Unable to get Node : %s\n", err)
		return nil
	}
	return getNodeOutput.Node
}

func (ambNode *NetworkConfig) UpdateNode(networkId string, memberId string, nodeId string) {
	var updateNodeInput managedblockchain.UpdateNodeInput
	updateNodeInput.SetNetworkId(networkId)
	updateNodeInput.SetMemberId(memberId)
	updateNodeInput.SetNodeId(nodeId)

	_, err := ambNode.Amb.UpdateNode(&updateNodeInput)
	if err != nil {
		log.Printf("Unable to Update the Node : %s\n", err)
	}
	fmt.Printf("Updated Node Successfully NodeId : %s\n", nodeId)
}

func (ambNode *NetworkConfig) ListNodes(networkId string, memberId string) []*managedblockchain.NodeSummary {
	var listNodeInput managedblockchain.ListNodesInput
	listNodeInput.SetNetworkId(networkId)
	listNodeInput.SetMemberId(memberId)
	listNodeInput.SetMaxResults(5) //AWS AMB supports max 5 peer nodes
	listNodeOutput, err := ambNode.Amb.ListNodes(&listNodeInput)
	if err != nil {
		log.Printf("Unable to list the Nodes : %s\n", err)
	}
	return listNodeOutput.Nodes
}
