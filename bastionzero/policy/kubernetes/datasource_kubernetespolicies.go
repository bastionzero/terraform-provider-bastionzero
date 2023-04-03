package kubernetes

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

func NewKubernetesPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSourceWithPractitionerParameters(
		&bzdatasource.ListDataSourceWithPractitionerParametersConfig[KubernetesPolicyModel, policy.ListPolicyParametersModel, policies.KubernetesPolicy]{
			BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[KubernetesPolicyModel, policies.KubernetesPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeKubernetesPolicyResourceSchema(), nil),
				MetadataTypeName:    "kubernetes_policies",
				ResultAttributeName: "policies",
				PrettyAttributeName: "Kubernetes policies",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.KubernetesPolicy) (state *KubernetesPolicyModel, diags diag.Diagnostics) {
					state = new(KubernetesPolicyModel)
					SetKubernetesPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				Description: "Get a list of all Kubernetes policies in your BastionZero organization.",
			},
			PractitionerParamsRecordSchema: policy.ListPolicyParametersSchema(),
			ListAPIModels: func(ctx context.Context, listParameters policy.ListPolicyParametersModel, client *bastionzero.Client) ([]policies.KubernetesPolicy, error) {
				subjectsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Subjects), ",")
				groupsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Groups), ",")

				policies, _, err := client.Policies.ListKubernetesPolicies(ctx, &policies.ListPolicyOptions{Subjects: subjectsFilter, Groups: groupsFilter})
				return policies, err
			},
		})
}
