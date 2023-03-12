package main

import (
	"fmt"
	"io"
	"os"

	"context"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/proto"

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
	fmt.Println("Host Projects:")
	hostProjects := app.GetDistinctProjects(allSubnets)
	for _, p := range hostProjects {
		fmt.Printf("- %v\n", p)
	}
	fmt.Println("Service Projects:")
	serviceProjects, _ := app.RetrieveAllServiceProjects(hostProjects)
	for _, p := range serviceProjects {
		fmt.Printf("- %v\n", p)
	}
	fmt.Println("VMs:")
	projects := append(hostProjects, serviceProjects...)
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

	logger.Info("Start diagram creation")
	var drawer app.Drawer
	drawer, err = app.NewVMDiagramDrawer(vpcs, peerings, allSubnets, allVMs)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = drawer.Draw()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = drawer.Render()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info("Diagram created successfully")

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
