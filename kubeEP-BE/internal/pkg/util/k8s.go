package util

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseGCPNodePoolByNodeName(projectName, nodeName string) string {
	nodePoolName := strings.ReplaceAll(nodeName, fmt.Sprintf("gke-%s-", projectName), "")
	re := regexp.MustCompile("-[a-z0-9]{8}-[a-z0-9]{4}$")
	return re.ReplaceAllString(nodePoolName, "")
}
