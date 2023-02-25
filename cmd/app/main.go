package main

import (
	"fmt"
	"io"
	"os"

	"context"
	"path/filepath"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/lib/textmeasure"

	"github.com/spf13/pflag"
	"gitlab.com/garzelli95/go-net-gen/internal/app"
	"gitlab.com/garzelli95/go-net-gen/internal/gcputils"
	"gitlab.com/garzelli95/go-net-gen/internal/log"
)

const (
	exitError      = 1
	exitUnexpected = 125
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(exitUnexpected)
		}
	}()
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitError)
	}
}

func run(args []string, _ io.Reader, _ io.Writer) error {
	// init viper and pflag
	configureDefaultSettings()

	// parse CLI arguments
	pflag.Parse()

	// load configuration
	config, err := loadConfiguration()
	if err != nil {
		return err
	}

	// now configuration is loaded, but not necessarily valid

	logger := log.NewLogger(config.Log) // create logger (log config is valid)
	log.SetStandardLogger(logger)       // override the global standard logger

	logger.Debug("Loaded configuration")

	// validate configuration
	err = config.Validate()
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("configuration is invalid: %w", err)
	}

	logger.Info("App started", map[string]interface{}{
		"name":           config.App.Name,
		"hub_project_id": config.App.HubProject,
	})

	logger.Info("Setup completed")

	logger.Info("Start gathering information from GCP")
	vpcs, peerings, err := app.RetrieveVPCsAndPeerings(config.App.HubProject)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	fmt.Println("VPCs:")
	for _, vpc := range vpcs {
		fmt.Printf("- %v in project %v\n", vpc.Name, vpc.Project)
	}
	fmt.Println("Peerings:")
	for _, peering := range peerings {
		fmt.Printf("- %v <-> %v\n",
			gcputils.GetVPCName(peering.VPC1SelfLink),
			gcputils.GetVPCName(peering.VPC2SelfLink))
	}
	fmt.Println("Subnets by VPC:")
	for _, vpc := range vpcs {
		fmt.Printf("- vpc: %v\n", vpc.Name)
		subnets, err := app.RetrieveVPCSubnets(vpc)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		for _, subnet := range subnets {
			fmt.Printf("  - %v %v\n", subnet.IPv4Range, subnet.Name)
		}
	}
	fmt.Println("Subnets:")
	allSubnets := []app.Subnet{}
	for _, vpc := range vpcs {
		subnets, err := app.RetrieveVPCSubnets(vpc)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		allSubnets = append(allSubnets, subnets...)
	}
	for _, subnet := range allSubnets {
		fmt.Printf("- %v %v\n", subnet.Name, subnet.IPv4Range)
	}
	fmt.Println("Projects:")
	projects := app.GetDistinctProjects(allSubnets)
	for _, p := range projects {
		fmt.Printf("- %v\n", p)
	}
	fmt.Println("VMs:")
	allVMs := []app.VM{}
	for _, project := range projects {
		vms, err := app.RetrieveVMs(project)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		allVMs = append(allVMs, vms...)
	}
	for _, vm := range allVMs {
		fmt.Printf("- %v %v\n", vm.Name, vm.InternalIP)
	}
	logger.Info("Finished gathering information from GCP")

	logger.Info("End")

	return nil
}

func ListAllInstances(projectID string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer instancesClient.Close()

	// Use the `MaxResults` parameter to limit the number of results that the
	// API returns per response page.
	req := &computepb.AggregatedListInstancesRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(3),
	}

	it := instancesClient.AggregatedList(ctx, req)

	// Despite using the `MaxResults` parameter, you don't need to handle the
	// pagination yourself. The returned iterator object handles pagination
	// automatically, returning separated pages as you iterate over the results.
	for {
		pair, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		instances := pair.Value.Instances
		for _, vm := range instances {
			nic := vm.GetNetworkInterfaces()[0]
			fmt.Printf("- %s in %s\n", vm.GetName(), pair.Key)
			fmt.Printf("  %s in %s (%s)\n", nic.GetNetworkIP(), nic.GetNetwork(), nic.GetSubnetwork())
		}
	}
	return nil
}

func ListNetworks(config app.Config) error {
	projectID := config.HubProject

	ctx := context.Background()
	networksClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return err
	}
	defer networksClient.Close()

	req := &computepb.ListNetworksRequest{
		Project:    projectID,
		MaxResults: proto.Uint32(3),
	}
	it := networksClient.List(ctx, req)
	for {
		network, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		fmt.Printf("Peerings of %v:\n", network.GetName())
		for _, peering := range network.Peerings {
			fmt.Printf("- %v\n", gcputils.GetVPCName(peering.GetNetwork()))
		}

		fmt.Printf("Subnets of %v:\n", network.GetSelfLink())
		for _, subnet := range network.GetSubnetworks() {
			fmt.Printf("- %v (%v)\n", subnet, network.GetIPv4Range())
		}

		fmt.Println("VMs:")
		ListAllInstances(projectID)
		for _, peering := range network.Peerings {
			ListAllInstances(gcputils.GetVPCProject(peering.GetNetwork()))
		}
	}

	return nil
}

func CreateDiagram(config app.Config) error {
	// Start with a new, empty graph g
	_, g, err := d2lib.Compile(context.Background(), "", nil)
	if err != nil {
		return err
	}

	// Create shapes
	g, _, _ = d2oracle.Create(g, "hub")
	g, _, _ = d2oracle.Create(g, "spk1")
	g, _, _ = d2oracle.Create(g, "spk2")

	g, _, _ = d2oracle.Create(g, "hub.prj")
	g, _, _ = d2oracle.Create(g, "spk1.prj")
	g, _, _ = d2oracle.Create(g, "spk2.prj")

	// Assign labels
	label1, label2 := "Spoke 1 VPC", "Spoke 2 VPC"
	g, _ = d2oracle.Set(g, "spk1.label", nil, &label1)
	g, _ = d2oracle.Set(g, "spk2.label", nil, &label2)

	// Create connections
	vpcPeering := "VPC Peering"
	g, k1, _ := d2oracle.Create(g, "hub <-> spk1")
	g, _ = d2oracle.Set(g, fmt.Sprintf("%s.label", k1), nil, &vpcPeering)
	g, k2, _ := d2oracle.Create(g, "hub <-> spk2")
	g, _ = d2oracle.Set(g, fmt.Sprintf("%s.label", k2), nil, &vpcPeering)

	// Beutify with icons
	image := "image"
	vpcIcon := "https://icons.terrastruct.com/gcp%2FProducts%20and%20services%2FNetworking%2FVirtual%20Private%20Cloud.svg"
	g, _ = d2oracle.Set(g, "hub.prj.icon", nil, &vpcIcon)
	g, _ = d2oracle.Set(g, "hub.prj.shape", nil, &image)

	// Turn the graph into a script
	script := d2format.Format(g.AST)

	// Initialize a ruler to measure font glyphs
	ruler, _ := textmeasure.NewRuler()

	// Compile the script into a diagram
	ctx := context.Background()
	diagram, _, _ := d2lib.Compile(ctx, script, &d2lib.CompileOptions{
		Layout: d2dagrelayout.DefaultLayout,
		Ruler:  ruler,
	})

	// Render to SVG
	diagramImage, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad: d2svg.DEFAULT_PADDING,
	})

	// Write to disk the script and the SVG image
	_ = os.WriteFile(filepath.Join("out.svg"), diagramImage, 0600)
	_ = os.WriteFile(filepath.Join("out.d2"), []byte(script), 0600)

	return nil
}
