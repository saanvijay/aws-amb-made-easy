package ambutils

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/managedblockchain"
)

type NetworkConfig struct {
	Region        string
	FabricVersion string
	NetworkName   string
	ChannelName   string
	OrgList       []string
	PeersPerOrg   int
	Amb           *managedblockchain.ManagedBlockchain
}

func (nConfig *NetworkConfig) CreateSession() *managedblockchain.ManagedBlockchain {
	awsConfig := aws.Config{Region: aws.String(nConfig.Region)}
	sess, err := session.NewSession(&awsConfig)
	if err != nil {
		log.Fatalf("Unable to create session : %s\n", err)
	}

	amb := managedblockchain.New(sess)
	fmt.Printf("Service Name %s\n", amb.ServiceName)
	nConfig.Amb = amb
	return amb
}

func (nConfig *NetworkConfig) GetApprovalThresholdPolicy() *managedblockchain.ApprovalThresholdPolicy {
	var appThreshPoly managedblockchain.ApprovalThresholdPolicy

	appThreshPoly.SetProposalDurationInHours(24)
	appThreshPoly.SetThresholdPercentage(20)
	appThreshPoly.SetThresholdComparator("GREATER_THAN")

	return &appThreshPoly
}

func (nConfig *NetworkConfig) GetMemberFabricConfiguration() *managedblockchain.MemberFabricConfiguration {
	var memberFabricConfiguration managedblockchain.MemberFabricConfiguration
	// AdminPassword must be at least 8 characters long and must contain at least one
	// uppercase character, one lowercase character, and one digit. It must not
	// contain ', ", \, /, @ or spaces. It must not exceed 32 characters in length.
	memberFabricConfiguration.SetAdminUsername("admin")
	memberFabricConfiguration.SetAdminPassword("Adminpwd1!")

	return &memberFabricConfiguration
}

func (nConfig *NetworkConfig) GetNetworkFrameworkConfiguration() *managedblockchain.NetworkFrameworkConfiguration {
	var frameworkConfiguration managedblockchain.NetworkFrameworkConfiguration
	var networkFabricConfiguration managedblockchain.NetworkFabricConfiguration

	networkFabricConfiguration.SetEdition("STARTER")
	frameworkConfiguration.SetFabric(&networkFabricConfiguration)

	return &frameworkConfiguration
}
func (nConfig *NetworkConfig) GetMemberFrameworkConfiguration() *managedblockchain.MemberFrameworkConfiguration {
	var memeberFrameworkConfiguration managedblockchain.MemberFrameworkConfiguration
	memeberFrameworkConfiguration.SetFabric(nConfig.GetMemberFabricConfiguration())

	return &memeberFrameworkConfiguration
}

func (nConfig *NetworkConfig) GetMemberConfiguration() *managedblockchain.MemberConfiguration {
	var memberConfiguration managedblockchain.MemberConfiguration

	memberConfiguration.SetDescription("member description")
	memberConfiguration.SetFrameworkConfiguration(nConfig.GetMemberFrameworkConfiguration())
	memberConfiguration.SetName(nConfig.OrgList[0])

	return &memberConfiguration
}

func (nConfig *NetworkConfig) GetVotingPolicy() *managedblockchain.VotingPolicy {

	var votPolicy managedblockchain.VotingPolicy
	votPolicy.SetApprovalThresholdPolicy(nConfig.GetApprovalThresholdPolicy())

	return &votPolicy
}

func (nConfig *NetworkConfig) GetNetworkStatus(networkId string) string {
	var networkInput managedblockchain.GetNetworkInput
	networkInput.SetNetworkId(networkId)
	networkOut, err := nConfig.Amb.GetNetwork(&networkInput)
	if err != nil {
		log.Fatalf("Unable to get the network : %s\n", err)
	}
	return *networkOut.Network.Status
}

func (nConfig *NetworkConfig) CreateNetwork() *managedblockchain.CreateNetworkOutput {

	// Any AWS service session should be created
	amb := nConfig.CreateSession()

	var networkInput managedblockchain.CreateNetworkInput
	networkInput.SetName(nConfig.NetworkName)
	networkInput.SetFrameworkVersion(nConfig.FabricVersion)
	networkInput.SetVotingPolicy(nConfig.GetVotingPolicy())
	networkInput.SetFramework("HYPERLEDGER_FABRIC")
	networkInput.SetFrameworkConfiguration(nConfig.GetNetworkFrameworkConfiguration())
	networkInput.SetMemberConfiguration(nConfig.GetMemberConfiguration())

	fabnetwork, err := amb.CreateNetwork(&networkInput)
	if err != nil {
		log.Fatalf("Unable to create the network : %s\n", err)
	}

	// Until network is available don't create peer nodes
	for {
		status := nConfig.GetNetworkStatus(*fabnetwork.NetworkId)
		fmt.Printf("Network status %s...\n", status)
		if status == "AVAILABLE" {
			break
		}
		time.Sleep(30 * time.Second)
	}
	fmt.Println("Network Created Successfully")
	// Create peer nodes
	for i := 0; i < nConfig.PeersPerOrg; i++ {
		fmt.Println("Creating Peer Node...")
		peerNode := nConfig.CreateNode(*fabnetwork.NetworkId, *fabnetwork.MemberId)
		if err != nil {
			log.Fatalf("Unable to create Peer Node %s\n", err)
		}
		fmt.Printf("Created Peer Node DONE")
	}

	return fabnetwork
}
