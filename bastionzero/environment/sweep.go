package environment

import (
	"context"
	"log"
	"strings"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("bastionzero_environment", &resource.Sweeper{
		Name: "bastionzero_environment",
		F:    sweepEnvironment,
	})

}

func sweepEnvironment(region string) error {
	client, err := sweep.SweeperClient()
	if err != nil {
		return err
	}

	envs, _, err := client.Environments.ListEnvironments(context.Background())
	if err != nil {
		return err
	}

	for _, env := range envs {
		if strings.HasPrefix(env.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying environment %s (%s)", env.Name, env.ID)

			if _, err := client.Environments.DeleteEnvironment(context.Background(), env.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
