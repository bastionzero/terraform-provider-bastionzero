package clustertarget_test

import (
	"context"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(clusterTarget *targets.ClusterTarget) map[string]string {
	valuesCheckMap := map[string]string{
		"agent_public_key":                   clusterTarget.AgentPublicKey,
		"agent_version":                      clusterTarget.AgentVersion,
		"control_channel.connection_node_id": clusterTarget.ControlChannel.ConnectionNodeID,
		"control_channel.control_channel_id": clusterTarget.ControlChannel.ControlChannelID,
		"control_channel.start_time":         clusterTarget.ControlChannel.StartTime.UTC().Format(time.RFC3339),
		"environment_id":                     clusterTarget.EnvironmentID,
		"region":                             clusterTarget.Region,
		"status":                             string(clusterTarget.Status),
		"type":                               string(targettype.Bzero),
		"id":                                 clusterTarget.ID,
		"name":                               clusterTarget.Name,
	}

	if clusterTarget.LastAgentUpdate != nil {
		valuesCheckMap["last_agent_update"] = clusterTarget.LastAgentUpdate.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["last_agent_update"] = ""
	}
	if clusterTarget.ControlChannel.EndTime != nil {
		valuesCheckMap["control_channel.end_time"] = clusterTarget.ControlChannel.EndTime.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["control_channel.end_time"] = ""
	}

	// TODO-Yuval: Find a way to test valid_cluster_users without depending on
	// index ordering. TestCheckTypeSetElemNestedAttrs() doesn't seem to be able
	// to assert a nested object/sets/lists within a nested object

	return valuesCheckMap
}

func TestAccClusterTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_cluster_targets.test"
	clusterTarget := new(targets.ClusterTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNClusterTargetsOrSkip(t, clusterTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTargetsDataSourceConfig(),
				// Check that the Cluster target we queried for is returned in
				// the list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "cluster_targets.*", getValuesCheckMap(clusterTarget)),
			},
		},
	})
}

func testAccClusterTargetsDataSourceConfig() string {
	return `
data "bastionzero_cluster_targets" "test" {
}
`
}
