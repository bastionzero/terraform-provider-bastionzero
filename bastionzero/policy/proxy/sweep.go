package proxy

import (
	"context"
	"log"
	"strings"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("bastionzero_proxy_policy", &resource.Sweeper{
		Name: "bastionzero_proxy_policy",
		F:    sweepProxyPolicy,
	})

}

func sweepProxyPolicy(region string) error {
	client, err := sweep.SweeperClient()
	if err != nil {
		return err
	}

	policies, _, err := client.Policies.ListProxyPolicies(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if strings.HasPrefix(policy.Name, sweep.TestNamePrefix) {
			log.Printf("Destroying Proxy policy %s (%s)", policy.Name, policy.ID)

			if _, err := client.Policies.DeleteProxyPolicy(context.Background(), policy.ID); err != nil {
				return err
			}
		}
	}

	return nil
}
