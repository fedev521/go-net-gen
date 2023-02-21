package gcputils

import "strings"

// Extract VPC name from a full VPC self link.
func GetVPCName(selfLink string) string {
	return strings.Split(selfLink, "/")[9]
}

// Extract project name from a full VPC self link.
func GetVPCProject(selfLink string) string {
	return strings.Split(selfLink, "/")[6]
}

// -----------------------------------------------------------------------------

// Extract subnet name from a full subnet self link.
func GetSubnetName(selfLink string) string {
	return strings.Split(selfLink, "/")[10]
}

// Extract subnet project from a full subnet self link.
func GetSubnetProject(selfLink string) string {
	return strings.Split(selfLink, "/")[6]
}

// Extract subnet region from a full subnet self link.
func GetSubnetRegion(selfLink string) string {
	return strings.Split(selfLink, "/")[8]
}
