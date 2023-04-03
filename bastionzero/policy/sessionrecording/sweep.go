package sessionrecording

import (
	"context"
	"log"
	"strings"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("bastionzero_sessionrecording_policy", &resource.Sweeper{
		Name: "bastionzero_sessionrecording_policy",
		F:    sweepSessionRecordingPolicy,
	})

}

func sweepSessionRecordingPolicy(region string) error {
	client, err := sweep.SweeperClient()
	if err != nil {
		return err
	}

	policies, _, err := client.Policies.ListSessionRecordingPolicies(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if strings.HasPrefix(policy.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying session recording policy %s (%s)", policy.Name, policy.ID)

			if _, err := client.Policies.DeleteSessionRecordingPolicy(context.Background(), policy.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
