package sessionrecording

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewSessionRecordingPolicyDataSource() datasource.DataSource {
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[sessionRecordingPolicyModel, policies.SessionRecordingPolicy]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[sessionRecordingPolicyModel, policies.SessionRecordingPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeSessionRecordingPolicyResourceSchema(), bastionzero.PtrTo("id")),
				MetadataTypeName:    "sessionrecording_policy",
				PrettyAttributeName: "session recording policy",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.SessionRecordingPolicy, state *sessionRecordingPolicyModel) (diags diag.Diagnostics) {
					setSessionRecordingPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel sessionRecordingPolicyModel, client *bastionzero.Client) (*policies.SessionRecordingPolicy, error) {
					policy, _, err := client.Policies.GetSessionRecordingPolicy(ctx, tfModel.ID.ValueString())
					return policy, err
				},
				Description: "Get information on a BastionZero session recording policy.",
			},
		},
	)
}
