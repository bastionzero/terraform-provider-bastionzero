package bzerotarget_test

import (
	"context"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(bzeroTarget *targets.BzeroTarget) map[string]string {
	valuesCheckMap := map[string]string{
		"agent_public_key":                   bzeroTarget.AgentPublicKey,
		"agent_version":                      bzeroTarget.AgentVersion,
		"control_channel.connection_node_id": bzeroTarget.ControlChannel.ConnectionNodeID,
		"control_channel.control_channel_id": bzeroTarget.ControlChannel.ControlChannelID,
		"control_channel.start_time":         bzeroTarget.ControlChannel.StartTime.UTC().Format(time.RFC3339),
		"environment_id":                     bzeroTarget.EnvironmentID,
		"region":                             bzeroTarget.Region,
		"status":                             string(bzeroTarget.Status),
		"type":                               string(targettype.Bzero),
		"id":                                 bzeroTarget.ID,
		"name":                               bzeroTarget.Name,
	}

	if bzeroTarget.LastAgentUpdate != nil {
		valuesCheckMap["last_agent_update"] = bzeroTarget.LastAgentUpdate.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["last_agent_update"] = ""
	}
	if bzeroTarget.ControlChannel.EndTime != nil {
		valuesCheckMap["control_channel.end_time"] = bzeroTarget.ControlChannel.EndTime.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["control_channel.end_time"] = ""
	}

	return valuesCheckMap
}

func TestAccBzeroTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_bzero_targets.test"
	bzeroTarget := new(targets.BzeroTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBzeroTargetsDataSourceConfig(),
				// Check that the Bzero target we queried for is returned in the
				// list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "targets.*", getValuesCheckMap(bzeroTarget)),
			},
		},
	})
}

func testAccBzeroTargetsDataSourceConfig() string {
	return `
data "bastionzero_bzero_targets" "test" {
}
`
}
