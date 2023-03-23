package kubernetes

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewKubernetesPolicyDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[kubernetesPolicyModel, policies.KubernetesPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[kubernetesPolicyModel, policies.KubernetesPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeKubernetesPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "kubernetes_policy",
				PrettyAttributeName: "Kubernetes policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.KubernetesPolicy, state *kubernetesPolicyModel) (diags diag.Diagnostics) {
					setKubernetesPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel kubernetesPolicyModel, client *bastionzero.Client) (*policies.KubernetesPolicy, error) {
					policy, _, err := client.Policies.GetKubernetesPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				Description: "Get information on a BastionZero Kubernetes policy.",
			},
		},
	)
}
