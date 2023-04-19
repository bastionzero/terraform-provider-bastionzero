package dactarget_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(dacTarget *targets.DynamicAccessConfiguration) map[string]string {
	valuesCheckMap := map[string]string{
		"start_webhook":  dacTarget.StartWebhook,
		"stop_webhook":   dacTarget.StopWebhook,
		"health_webhook": dacTarget.HealthWebhook,
		"environment_id": dacTarget.EnvironmentId,
		"status":         string(dacTarget.Status),
		"type":           string(targettype.DynamicAccessConfig),
		"id":             dacTarget.ID,
		"name":           dacTarget.Name,
	}

	return valuesCheckMap
}

func TestAccDACTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_dac_targets.test"
	dacTarget := new(targets.DynamicAccessConfiguration)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNDACTargetsOrSkip(t, dacTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDACTargetsDataSourceConfig(),
				// Check that the DAC target we queried for is returned in the
				// list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "targets.*", getValuesCheckMap(dacTarget)),
			},
		},
	})
}

func testAccDACTargetsDataSourceConfig() string {
	return `
data "bastionzero_dac_targets" "test" {
}
`
}
