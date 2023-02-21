package gcputils

import (
	"testing"
)

func TestGetVPCProject(t *testing.T) {
	type args struct {
		selfLink string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetVPCProject SelfLink",
			args: args{"https://www.googleapis.com/compute/v1/projects/prj-gonetgen-infra-hub/global/networks/vpc-gonetgen-hub-01"},
			want: "prj-gonetgen-infra-hub",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVPCProject(tt.args.selfLink); got != tt.want {
				t.Errorf("GetVPCProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubnetName(t *testing.T) {
	type args struct {
		selfLink string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetSubnetName SelfLink",
			args: args{"https://www.googleapis.com/compute/v1/projects/prj-gonetgen-infra-hub/regions/us-central1/subnetworks/subnet-gonetgen-hub-usc1-01"},
			want: "subnet-gonetgen-hub-usc1-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSubnetName(tt.args.selfLink); got != tt.want {
				t.Errorf("GetSubnetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubnetProject(t *testing.T) {
	type args struct {
		selfLink string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetSubnetProject SelfLink",
			args: args{"https://www.googleapis.com/compute/v1/projects/prj-gonetgen-infra-hub/regions/us-central1/subnetworks/subnet-gonetgen-hub-usc1-01"},
			want: "prj-gonetgen-infra-hub",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSubnetProject(tt.args.selfLink); got != tt.want {
				t.Errorf("GetSubnetProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubnetRegion(t *testing.T) {
	type args struct {
		selfLink string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetSubnetRegion SelfLink",
			args: args{"https://www.googleapis.com/compute/v1/projects/prj-gonetgen-infra-hub/regions/us-central1/subnetworks/subnet-gonetgen-hub-usc1-01"},
			want: "us-central1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSubnetRegion(tt.args.selfLink); got != tt.want {
				t.Errorf("GetSubnetRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
