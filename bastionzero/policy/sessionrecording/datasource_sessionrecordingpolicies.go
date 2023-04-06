package sessionrecording

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

func NewSessionRecordingPoliciesDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSourceWithPractitionerParameters(
		&bzdatasource.ListDataSourceWithPractitionerParametersConfig[SessionRecordingPolicyModel, policy.ListPolicyParametersModel, policies.SessionRecordingPolicy]{
			BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[SessionRecordingPolicyModel, policies.SessionRecordingPolicy]{
				RecordSchema:        internal.ResourceSchemaToDataSourceSchema(makeSessionRecordingPolicyResourceSchema(), nil),
				MetadataTypeName:    "sessionrecording_policies",
				ResultAttributeName: "policies",
				PrettyAttributeName: "session recording policies",
				FlattenAPIModel: func(ctx context.Context, apiObject *policies.SessionRecordingPolicy) (state *SessionRecordingPolicyModel, diags diag.Diagnostics) {
					state = new(SessionRecordingPolicyModel)
					SetSessionRecordingPolicyAttributes(ctx, state, apiObject, true)
					return
				},
				MarkdownDescription: "Get a list of all session recording policies in your BastionZero organization. A session recording policy governs whether users' I/O during shell connections are recorded.",
			},
			PractitionerParamsRecordSchema: policy.ListPolicyParametersSchema(),
			ListAPIModels: func(ctx context.Context, listParameters policy.ListPolicyParametersModel, client *bastionzero.Client) ([]policies.SessionRecordingPolicy, error) {
				subjectsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Subjects), ",")
				groupsFilter := strings.Join(internal.ExpandFrameworkStringSet(ctx, listParameters.Groups), ",")

				policies, _, err := client.Policies.ListSessionRecordingPolicies(ctx, &policies.ListPolicyOptions{Subjects: subjectsFilter, Groups: groupsFilter})
				return policies, err
			},
		})
}
