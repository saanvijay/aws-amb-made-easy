package ambutils

import (
	"log"

	"github.com/aws/aws-sdk-go/service/managedblockchain"
)

func (ambMember *NetworkConfig) CreateMember(networkId string, memberName string, clientToken string, invitationId string) string {
	var memberInput managedblockchain.CreateMemberInput

	memberInput.SetInvitationId(invitationId)
	memberInput.SetNetworkId(networkId)
	memberInput.SetMemberConfiguration(ambMember.GetMemberConfiguration(memberName))
	memberInput.SetClientRequestToken(clientToken)
	memberOutput, err := ambMember.Amb.CreateMember(&memberInput)
	if err != nil {
		log.Fatalf("Unable to Create Member %s\n", err)
	}
	return *memberOutput.MemberId
}

func (ambMember *NetworkConfig) DeleteMember(networkId string, memberId string) {
	var memberInput managedblockchain.DeleteMemberInput
	memberInput.SetNetworkId(networkId)
	memberInput.SetMemberId(memberId)

	_, err := ambMember.Amb.DeleteMember(&memberInput)
	if err != nil {
		log.Printf("Unable to delete the member : %s\n", err)
		return
	}
	log.Printf("Member deleted successfully : %s\n", memberId)
}

func (ambMember *NetworkConfig) GetMember(networkId string, memberId string) *managedblockchain.Member {
	var getMemberInput managedblockchain.GetMemberInput
	getMemberInput.SetNetworkId(networkId)
	getMemberInput.SetMemberId(memberId)
	getMemberOut, err := ambMember.Amb.GetMember(&getMemberInput)
	if err != nil {
		log.Printf("Unable to get the member : %s\n", err)
		return nil
	}
	return getMemberOut.Member
}

func (ambMember *NetworkConfig) UpdateMember(networkId string, memberId string) {
	var updateMemberInput managedblockchain.UpdateMemberInput
	updateMemberInput.SetNetworkId(networkId)
	updateMemberInput.SetMemberId(memberId)
	_, err := ambMember.Amb.UpdateMember(&updateMemberInput)
	if err != nil {
		log.Printf("Unable to get the member : %s\n", err)
		return
	}
	log.Printf("Mamber Updated Successfully : %s\n", memberId)
}
