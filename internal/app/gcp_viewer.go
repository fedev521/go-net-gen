package app

import (
	"context"
	"errors"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"gitlab.com/garzelli95/go-net-gen/internal/gcputils"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"
)

// Given a project id (of the "hub"), RetrieveVPCsAndPeerings returns two slices
// (after querying Google API). One contains all VPCs in the hub project as well
// as all VPCs peered to that project; the other information about peergings, if
// any.
func RetrieveVPCsAndPeerings(hubProjectID string) ([]VPC, []VPCPeering, error) {
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
			Network: gcputils.GetVPCName(peerSelfLink),
			Project: gcputils.GetVPCProject(peerSelfLink),
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

// Given a VPC structure, RetrieveVPCSubnets returns a slice of the subnets
// belonging to that VPC (after querying Google API).
func RetrieveVPCSubnets(vpc VPC) ([]Subnet, error) {
	subnets := []Subnet{}

	ctx := context.Background()
	subnetsClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return []Subnet{}, err
	}
	defer subnetsClient.Close()

	// retrieve information about subnets within each VPC (in parallel)
	n := len(vpc.Subnets)
	ch := make(chan *computepb.Subnetwork, n)
	for _, subnetSelfLink := range vpc.Subnets {
		go func(ssl string, ch chan *computepb.Subnetwork) {
			req := &computepb.GetSubnetworkRequest{
				Project:    gcputils.GetSubnetProject(ssl),
				Region:     gcputils.GetSubnetRegion(ssl),
				Subnetwork: gcputils.GetSubnetName(ssl),
			}
			subnet, err := subnetsClient.Get(ctx, req)
			if err != nil {
				ch <- nil
			} else {
				ch <- subnet
			}
		}(subnetSelfLink, ch)
	}

	// gather results and build return value
	for i := 0; i < n; i++ {
		subnet := <-ch
		if subnet == nil {
			return []Subnet{}, errors.New("could not retrieve subnets information")
		}
		subnets = append(subnets, NewSubnet(subnet)) // build each Subnet
	}

	return subnets, nil
}

// Given a slice of subnets, GetDistinctProjects returns a slice with no
// duplicates of projects containing those subnets.
func GetDistinctProjects(subnets []Subnet) []string {
	projects := []string{}
	presence := make(map[string]bool)

	for _, subnet := range subnets {
		if !presence[subnet.Project] {
			presence[subnet.Project] = true
			projects = append(projects, subnet.Project)
		}
	}

	return projects
}

func RetrieveVMs(projectID string) ([]VM, error) {
	vms := []VM{}

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return []VM{}, err
	}
	defer instancesClient.Close()

	// get information about VM instances
	req := &computepb.AggregatedListInstancesRequest{
		Project: projectID,
	}
	it := instancesClient.AggregatedList(ctx, req)

	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []VM{}, err
		}

		// build each VM (for each <zone,vm> pair)
		for _, vm := range pair.Value.Instances {
			vms = append(vms, NewVM(vm))
		}
	}

	return vms, nil
}
