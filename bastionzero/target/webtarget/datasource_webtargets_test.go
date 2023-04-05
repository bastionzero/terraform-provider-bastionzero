package webtarget_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(webTarget *targets.WebTarget) map[string]string {
	valuesCheckMap := map[string]string{
		"agent_public_key": webTarget.AgentPublicKey,
		"agent_version":    webTarget.AgentVersion,
		"environment_id":   webTarget.EnvironmentID,
		"region":           webTarget.Region,
		"status":           string(webTarget.Status),
		"type":             string(targettype.Web),
		"id":               webTarget.ID,
		"name":             webTarget.Name,
		"proxy_target_id":  webTarget.ProxyTargetID,
		"remote_host":      webTarget.RemoteHost,
	}

	if webTarget.LastAgentUpdate != nil {
		valuesCheckMap["last_agent_update"] = webTarget.LastAgentUpdate.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["last_agent_update"] = ""
	}

	if webTarget.LocalPort.Value != nil {
		valuesCheckMap["local_port"] = strconv.Itoa(*webTarget.LocalPort.Value)
	} else {
		valuesCheckMap["local_port"] = ""
	}
	if webTarget.RemotePort.Value != nil {
		valuesCheckMap["remote_port"] = strconv.Itoa(*webTarget.RemotePort.Value)
	} else {
		valuesCheckMap["remote_port"] = ""
	}

	return valuesCheckMap
}

func TestAccWebTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_web_targets.test"
	webTarget := new(targets.WebTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNWebTargetsOrSkip(t, webTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWebTargetsDataSourceConfig(),
				// Check that the Web target we queried for is returned in the
				// list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "web_targets.*", getValuesCheckMap(webTarget)),
			},
		},
	})
}

func testAccWebTargetsDataSourceConfig() string {
	return `
data "bastionzero_web_targets" "test" {
}
`
}
