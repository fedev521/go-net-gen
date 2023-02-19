package gcputils

import "strings"

// Extract VPC name from a full VPC self link.
func GetVPCFromVPC(selfLink string) string {
	return strings.Split(selfLink, "/")[9]
}

// Extract project name from a full VPC self link.
func GetProjectFromVPC(selfLink string) string {
	return strings.Split(selfLink, "/")[6]
}
