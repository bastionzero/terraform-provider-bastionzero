package proxy

import (
	"context"
	"strings"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewProxyPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSourceWithPractitionerParameters(
		&bzdatasource.ListDataSourceWithPractitionerParametersConfig[proxyPolicyModel, policy.ListPolicyParametersModel, policies.ProxyPolicy]{
			BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[proxyPolicyModel, policies.ProxyPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeProxyPolicyResourceSchema(), nil),
				MetadataTypeName:    "proxy_policies",
				ResultAttributeName: "policies",
				PrettyAttributeName: "proxy policies",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.ProxyPolicy) (state *proxyPolicyModel, diags diag.Diagnostics) {
					state = new(proxyPolicyModel)
					setProxyPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				Description: "Get a list of all proxy policies in your BastionZero organization.",
			},
			PractitionerParamsRecordSchema: policy.ListPolicyParametersSchema(),
			ListAPIModels: func(ctx context.Context, listParameters policy.ListPolicyParametersModel, client *bastionzero.Client) ([]policies.ProxyPolicy, error) {
				subjectsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Subjects), ",")
				groupsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Groups), ",")

				policies, _, err := client.Policies.ListProxyPolicies(ctx, &policies.ListPolicyOptions{Subjects: subjectsFilter, Groups: groupsFilter})
				return policies, err
			},
		})
}