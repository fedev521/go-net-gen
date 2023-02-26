package app

import (
	"context"
	"fmt"

	"gitlab.com/garzelli95/go-net-gen/internal/d2utils"
	"gitlab.com/garzelli95/go-net-gen/internal/gcputils"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
)

// Start from a project (hub). Get all of its VPCs and all of their peered VPCs.
// For each VPC consider the corresponding project (spokes). For each project,
// get all Compute Engine instances and load balancers. Create a diagram with
// VMs, highlighting subnets and projects. Consider Shared VPCs as well.

// VPC -> Subnet -> Project -> VM

type VMDiagramDrawer struct {
	// Graph representing the final diagram that will be rendered to an image.
	g *d2graph.Graph
	// Map that associates a resource unique identifier to the shape key in the
	// D2 diagram. For example, it may associate a VPC self link to the key of
	// the rectangle representing the VPC in the diagram.
	keys map[string]string

	vpcs     []VPC
	peerings []VPCPeering
	subnets  []Subnet
}

func NewVMDiagramDrawer(vpcs []VPC, peerings []VPCPeering, subnets []Subnet) (*VMDiagramDrawer, error) {
	_, graph, err := d2lib.Compile(context.Background(), "", nil)
	if err != nil {
		return &VMDiagramDrawer{}, err
	}

	d := VMDiagramDrawer{
		g:        graph,
		keys:     make(map[string]string),
		vpcs:     vpcs,
		peerings: peerings,
		subnets:  subnets,
	}

	return &d, nil
}

// add a configuration argument
func (d *VMDiagramDrawer) Draw() error {
	// draw VPC shapes and associate id to shape key
	for _, vpc := range d.vpcs {
		g, k, err := d2oracle.Create(d.g, d.vpcTmpKey(vpc))
		if err != nil {
			return err
		}
		d.keys[d.vpcId(vpc)] = k
		d.g = g
	}

	// draw peering connections and associate id to shape key
	for _, peering := range d.peerings {
		g, k, err := d2oracle.Create(d.g, d.peeringTmpKey(peering))
		if err != nil {
			return err
		}
		d.keys[d.peeringId(peering)] = k
		d.g = g
	}

	// draw subnet shapes under VPCs and associate id to shape key
	for _, subnet := range d.subnets {
		g, k, err := d2oracle.Create(d.g, d.subnetTmpKey(subnet))
		if err != nil {
			return err
		}
		d.keys[d.subnetId(subnet)] = k
		d.g = g
	}

	err := d.beautify()
	if err != nil {
		return err
	}

	return nil
}

func (d *VMDiagramDrawer) beautify() error {
	// set VPC labels
	for _, vpc := range d.vpcs {
		key := d.keys[d.vpcId(vpc)]
		label := d.vpcLabel(vpc)
		g, err := d2oracle.Set(d.g, fmt.Sprintf("%s.label", key), nil, &label)
		if err != nil {
			return err
		}
		d.g = g
	}

	// set peering labels
	for _, peering := range d.peerings {
		key := d.keys[d.peeringId(peering)]
		label := d.peeringLabel(peering)
		g, err := d2oracle.Set(d.g, fmt.Sprintf("%s.label", key), nil, &label)
		if err != nil {
			return err
		}
		d.g = g
	}

	// set VPC icons
	for _, vpc := range d.vpcs {
		key := d.keys[d.vpcId(vpc)]
		icon := "https://icons.terrastruct.com/gcp%2FProducts%20and%20services%2FNetworking%2FVirtual%20Private%20Cloud.svg"
		g, err := d2oracle.Set(d.g, fmt.Sprintf("%s.icon", key), nil, &icon)
		if err != nil {
			return err
		}
		d.g = g
	}

	// set subnet labels
	for _, subnet := range d.subnets {
		key := d.keys[d.subnetId(subnet)]
		label := d.subnetLabel(subnet)
		g, err := d2oracle.Set(d.g, fmt.Sprintf("%s.label", key), nil, &label)
		if err != nil {
			return err
		}
		d.g = g
	}

	return nil
}

// add a configuration argument
func (d *VMDiagramDrawer) Render() error {
	return d2utils.RenderSVG(d.g)
}

// -----------------------------------------------------------------------------

// A resource id is a string that identifies a GCP resource (e.g., a VPC, a VM,
// or a peering). The shape key is the key found in the d2 script, i.e., the one
// that Create() method returns. The temporary/tentative key is a string used as
// parameter in the Create() method; as such, it should not contain slashes or
// special characters. The label is the resource/shape displayed name (i.e., the
// value of the label property associated to the shape key).

// Resource ids mapped to (tentative) keys:
// - VPC: self link => name
// - VPCPeering: <vpc1sl> <-> <vpc2sl> => <vpc1name> <-> <vpc2name>
// - Subnet: self link => <vpcname>.<name>

func (d *VMDiagramDrawer) vpcId(vpc VPC) string {
	return vpc.SelfLink
}

func (d *VMDiagramDrawer) vpcTmpKey(vpc VPC) string {
	return vpc.Name
}

func (d *VMDiagramDrawer) vpcLabel(vpc VPC) string {
	return fmt.Sprintf("%s VPC", vpc.Name)
}

// ---

func (d *VMDiagramDrawer) peeringId(peering VPCPeering) string {
	sl1, sl2 := peering.VPC1SelfLink, peering.VPC2SelfLink
	return fmt.Sprintf("%s <-> %s", sl1, sl2)
}

func (d *VMDiagramDrawer) peeringTmpKey(peering VPCPeering) string {
	sl1, sl2 := peering.VPC1SelfLink, peering.VPC2SelfLink
	// NOTE: works as long as VPC tentative key and final one match. Should use
	// self link to get the VPC structure, compute the id with vpcId(), use the
	// id to find the true key in the keys map.
	return fmt.Sprintf("%s <-> %s", gcputils.GetVPCName(sl1), gcputils.GetVPCName(sl2))
}

func (d *VMDiagramDrawer) peeringLabel(peering VPCPeering) string {
	return "VPC Peering"
}

// ---

func (d *VMDiagramDrawer) subnetId(subnet Subnet) string {
	return subnet.SelfLink
}

func (d *VMDiagramDrawer) subnetTmpKey(subnet Subnet) string {
	vpcName := gcputils.GetVPCName(subnet.VPCSelfLink)
	// NOTE: works as long as VPC tentative key and final one match. Should use
	// self link to get the VPC structure, compute the id with vpcId(), use the
	// id to find the true key in the keys map.
	return fmt.Sprintf("%s.%s", vpcName, subnet.Name)
}

func (d *VMDiagramDrawer) subnetLabel(subnet Subnet) string {
	return fmt.Sprintf("Range %s", subnet.IPv4Range)
}
