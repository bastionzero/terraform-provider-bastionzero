package proxy_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/proxy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen/bzpolicygen"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen/tfpolicygen"
	"github.com/stretchr/testify/require"

	"pgregory.net/rapid"
)

func TestFlatExpandProxy_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genProxyPolicy := bzpolicygen.ProxyPolicyGen().Draw(t, "api")
		genTFSchemaProxyPolicy := tfpolicygen.ProxyPolicySchemaGen(context.Background()).Draw(t, "schema")

		// Flatten generated BastionZero API policy into Terraform schema type
		proxy.SetProxyPolicyAttributes(context.Background(), &genTFSchemaProxyPolicy, &genProxyPolicy, false)

		// Then expand the flattened value back to a BastionZero API policy type
		expanded := proxy.ExpandProxyPolicy(context.Background(), &genTFSchemaProxyPolicy)

		// And assert no data loss occurred when converting by asserting the
		// expanded value matches the original, generated policy
		require.EqualValues(t, genProxyPolicy.Name, expanded.Name)
		require.EqualValues(t, genProxyPolicy.GetPolicyType(), expanded.GetPolicyType())
		require.EqualValues(t, genProxyPolicy.GetDescription(), expanded.GetDescription())
		require.EqualValues(t, genProxyPolicy.GetSubjects(), expanded.GetSubjects())
		require.EqualValues(t, genProxyPolicy.GetGroups(), expanded.GetGroups())
		require.EqualValues(t, genProxyPolicy.GetEnvironments(), expanded.GetEnvironments())
		require.EqualValues(t, genProxyPolicy.GetTargets(), expanded.GetTargets())
		require.EqualValues(t, genProxyPolicy.GetTargetUsers(), expanded.GetTargetUsers())
	})
}
