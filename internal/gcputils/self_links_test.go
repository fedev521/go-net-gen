package gcputils

import "testing"

func TestGetProjectFromVPC(t *testing.T) {
	type args struct {
		selfLink string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetProjectFromVPC SelfLink",
			args: args{"https://www.googleapis.com/compute/v1/projects/prj-gonetgen-infra-hub/global/networks/vpc-gonetgen-hub-01"},
			want: "prj-gonetgen-infra-hub",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectFromVPC(tt.args.selfLink); got != tt.want {
				t.Errorf("GetProjectFromVPC() = %v, want %v", got, tt.want)
			}
		})
	}
}
