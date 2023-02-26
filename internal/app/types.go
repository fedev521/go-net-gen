package app

import (
	"cloud.google.com/go/compute/apiv1/computepb"
	"gitlab.com/garzelli95/go-net-gen/internal/gcputils"
)

// -----------------------------------------------------------------------------

type VPC struct {
	SelfLink string
	Name     string
	Project  string
	Subnets  []string
}

func NewVPC(pb *computepb.Network) VPC {
	return VPC{
		SelfLink: pb.GetSelfLink(),
		Name:     pb.GetName(),
		Project:  gcputils.GetVPCProject(pb.GetSelfLink()),
		Subnets:  pb.GetSubnetworks(),
	}
}

// -----------------------------------------------------------------------------

type VPCPeering struct {
	VPC1SelfLink, VPC2SelfLink string
}

// -----------------------------------------------------------------------------

type Subnet struct {
	SelfLink    string
	Name        string
	Project     string
	IPv4Range   string
	VPCSelfLink string
}

func NewSubnet(pb *computepb.Subnetwork) Subnet {
	return Subnet{
		SelfLink:    pb.GetSelfLink(),
		Name:        pb.GetName(),
		Project:     gcputils.GetSubnetProject(pb.GetSelfLink()),
		IPv4Range:   pb.GetIpCidrRange(),
		VPCSelfLink: pb.GetNetwork(),
	}
}

// -----------------------------------------------------------------------------

type VM struct {
	SelfLink       string
	Name           string
	Project        string
	Zone           string
	InternalIP     string
	ExternalIP     string
	VPCSelfLink    string
	SubnetSelfLink string
}

func NewVM(pb *computepb.Instance) VM {
	nic := pb.GetNetworkInterfaces()[0]

	extIP := ""
	ac := nic.GetAccessConfigs()
	if len(ac) > 0 {
		extIP = ac[0].GetNatIP()
	}

	return VM{
		SelfLink:       pb.GetSelfLink(),
		Name:           pb.GetName(),
		Project:        gcputils.GetVMProject(pb.GetSelfLink()),
		Zone:           gcputils.GetVMZone(pb.GetSelfLink()),
		ExternalIP:     extIP,
		InternalIP:     nic.GetNetworkIP(),
		VPCSelfLink:    nic.GetNetwork(),
		SubnetSelfLink: nic.GetSubnetwork(),
	}
}
