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

	err = ListNetworks(config.App)
	if err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("cannot connect to GCP: %w", err)
	}

	return nil
}

func ListAllInstances(config app.Config) error {
	projectID := config.HubProject

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
		if len(instances) > 0 {
			fmt.Printf("%s\n", pair.Key)
			for _, instance := range instances {
				fmt.Printf("- %s %s\n", instance.GetName(), instance.GetMachineType())
			}
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

		fmt.Printf("Peering: %v\n", network.Peerings)
	}

	return nil
}

func D2hello() {
	ctx := context.Background()
	// Start with a new, empty graph
	_, graph, _ := d2lib.Compile(ctx, "", nil)

	// Create a shape, "meow"
	graph, _, _ = d2oracle.Create(graph, "meow")

	// Turn the graph into a script (which would just be "meow")
	script := d2format.Format(graph.AST)

	// Initialize a ruler to measure font glyphs
	ruler, _ := textmeasure.NewRuler()

	// Compile the script into a diagram
	diagram, _, _ := d2lib.Compile(context.Background(), script, &d2lib.CompileOptions{
		Layout: d2dagrelayout.DefaultLayout,
		Ruler:  ruler,
	})

	// Render to SVG
	out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad: d2svg.DEFAULT_PADDING,
	})

	// Write to disk
	_ = os.WriteFile(filepath.Join("out.svg"), out, 0600)
	_ = os.WriteFile(filepath.Join("out.d2"), []byte(script), 0600)
}
