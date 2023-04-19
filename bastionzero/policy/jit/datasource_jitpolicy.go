package jit

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewJITPolicyDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[JITPolicyModel, policies.JITPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[JITPolicyModel, policies.JITPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeJITPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "jit_policy",
				PrettyAttributeName: "JIT policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.JITPolicy, state *JITPolicyModel) (diags diag.Diagnostics) {
					SetJITPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel JITPolicyModel, client *bastionzero.Client) (*policies.JITPolicy, error) {
					policy, _, err := client.Policies.GetJITPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				MarkdownDescription: "Get information on a BastionZero JIT policy. A JIT policy provides just in time access to targets.",
			},
		},
	)
}
