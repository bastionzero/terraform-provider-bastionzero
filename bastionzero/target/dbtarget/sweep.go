package dbtarget

import (
	"context"
	"log"
	"strings"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("bastionzero_db_target", &resource.Sweeper{
		Name: "bastionzero_db_target",
		F:    sweepDbTarget,
	})

}

func sweepDbTarget(region string) error {
	client, err := sweep.SweeperClient()
	if err != nil {
		return err
	}

	dbTargets, _, err := client.Targets.ListDatabaseTargets(context.Background())
	if err != nil {
		return err
	}

	for _, dbTarget := range dbTargets {
		if strings.HasPrefix(dbTarget.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Db target %s (%s)", dbTarget.Name, dbTarget.ID)

			if _, err := client.Targets.DeleteDatabaseTarget(context.Background(), dbTarget.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
