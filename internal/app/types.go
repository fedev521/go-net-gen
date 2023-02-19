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
}

func NewVPC(pb *computepb.Network) VPC {
	return VPC{
		SelfLink: pb.GetSelfLink(),
		Name:     pb.GetName(),
		Project:  gcputils.GetProjectFromVPC(pb.GetSelfLink()),
	}
}

// -----------------------------------------------------------------------------
type VPCPeering struct {
	VPC1SelfLink, VPC2SelfLink string
}

// -----------------------------------------------------------------------------
type Subnet struct {
	SelfLink  string
	Name      string
	Project   string
	IPv4Range string
}

func NewSubnet(pb *computepb.Subnetwork) Subnet {
	return Subnet{} // TODO implement
}

// -----------------------------------------------------------------------------
type VM struct {
	SelfLink   string
	Name       string
	Project    string
	Zone       string
	InternalIP string
	ExternalIP string
}

func NewVM(pb *computepb.Instance) VM {
	return VM{} // TODO implement
}
