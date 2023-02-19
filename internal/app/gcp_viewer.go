package app

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"gitlab.com/garzelli95/go-net-gen/internal/gcputils"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

func GetVPCsAndPeerings(hubProjectID string) ([]VPC, []VPCPeering, error) {
	// The function's results, i.e., all VPCs and peerings to be considered
	vpcs := []VPC{}
	peerings := []VPCPeering{}

	ctx := context.Background()
	networksClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return []VPC{}, []VPCPeering{}, err
	}
	defer networksClient.Close()

	// get all VPCs in the hub project
	listReq := &computepb.ListNetworksRequest{
		Project:    hubProjectID,
		MaxResults: proto.Uint32(10),
	}
	it := networksClient.List(ctx, listReq)

	hubNetworks := []*computepb.Network{}
	for {
		network, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []VPC{}, []VPCPeering{}, err
		}

		vpcs = append(vpcs, NewVPC(network))       // build each VPC
		hubNetworks = append(hubNetworks, network) // store each network pb
	}

	// find all VPCs peered to hub VPCs
	for _, hubNetwork := range hubNetworks {
		for _, peerNetwork := range hubNetwork.GetPeerings() {
			peerings = append(peerings, VPCPeering{
				VPC1SelfLink: hubNetwork.GetSelfLink(),
				VPC2SelfLink: peerNetwork.GetNetwork(),
			})
		}
	}

	// get information about peered VPCs
	for _, peering := range peerings {
		peerSelfLink := peering.VPC2SelfLink

		// TODO: parallellize Get requests
		getReq := &computepb.GetNetworkRequest{
			Network: gcputils.GetVPCFromVPC(peerSelfLink),
			Project: gcputils.GetProjectFromVPC(peerSelfLink),
		}

		network, err := networksClient.Get(ctx, getReq)
		if err != nil {
			return []VPC{}, []VPCPeering{}, err
		}

		vpcs = append(vpcs, NewVPC(network)) // build each VPC
	}

	// TODO: remove duplicates from vpcs and peerings

	return vpcs, peerings, nil
}
