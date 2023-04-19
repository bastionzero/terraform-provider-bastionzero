package dbtarget_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func getValuesCheckMap(dbTarget *targets.DatabaseTarget) map[string]string {
	valuesCheckMap := map[string]string{
		"agent_public_key": dbTarget.AgentPublicKey,
		"agent_version":    dbTarget.AgentVersion,
		"environment_id":   dbTarget.EnvironmentID,
		"is_split_cert":    fmt.Sprintf("%t", dbTarget.IsSplitCert),
		"region":           dbTarget.Region,
		"status":           string(dbTarget.Status),
		"type":             string(targettype.Db),
		"id":               dbTarget.ID,
		"name":             dbTarget.Name,
		"proxy_target_id":  dbTarget.ProxyTargetID,
		"remote_host":      dbTarget.RemoteHost,
	}

	if dbTarget.LastAgentUpdate != nil {
		valuesCheckMap["last_agent_update"] = dbTarget.LastAgentUpdate.UTC().Format(time.RFC3339)
	} else {
		valuesCheckMap["last_agent_update"] = ""
	}

	if dbTarget.LocalPort.Value != nil {
		valuesCheckMap["local_port"] = strconv.Itoa(*dbTarget.LocalPort.Value)
	} else {
		valuesCheckMap["local_port"] = ""
	}
	if dbTarget.RemotePort.Value != nil {
		valuesCheckMap["remote_port"] = strconv.Itoa(*dbTarget.RemotePort.Value)
	} else {
		valuesCheckMap["remote_port"] = ""
	}

	if dbTarget.DatabaseType != nil {
		valuesCheckMap["database_type"] = *dbTarget.DatabaseType
	} else {
		valuesCheckMap["database_type"] = ""
	}

	return valuesCheckMap
}

func TestAccDbTargetsDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_db_targets.test"
	dbTarget := new(targets.DatabaseTarget)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNDbTargetsOrSkip(t, dbTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetsDataSourceConfig(),
				// Check that the Db target we queried for is returned in the
				// list
				Check: resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "targets.*", getValuesCheckMap(dbTarget)),
			},
		},
	})
}

func testAccDbTargetsDataSourceConfig() string {
	return `
data "bastionzero_db_targets" "test" {
}
`
}
