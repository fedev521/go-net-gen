package app

import (
	"context"

	"gitlab.com/garzelli95/go-net-gen/internal/d2utils"
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

	vpcs []VPC
}

func NewVMDiagramDrawer(vpcs []VPC) (*VMDiagramDrawer, error) {
	_, graph, err := d2lib.Compile(context.Background(), "", nil)
	if err != nil {
		return &VMDiagramDrawer{}, err
	}

	d := VMDiagramDrawer{
		g:    graph,
		keys: make(map[string]string),
		vpcs: vpcs,
	}

	return &d, nil
}

// add a configuration argument
func (d *VMDiagramDrawer) Draw() error {
	// draw VPC shapes and associate self link to shape key
	for _, vpc := range d.vpcs {
		g, k, err := d2oracle.Create(d.g, vpc.Name)
		d.keys[vpc.SelfLink] = k
		if err != nil {
			return err
		}
		d.g = g
	}

	d.beautify()

	return nil
}

func (d *VMDiagramDrawer) beautify() {

}

// add a configuration argument
func (d *VMDiagramDrawer) Render() error {
	return d2utils.RenderSVG(d.g)
}
