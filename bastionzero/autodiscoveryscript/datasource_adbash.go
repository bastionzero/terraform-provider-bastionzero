package autodiscoveryscript

import (
	"context"
	"fmt"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/autodiscoveryscripts"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/autodiscoveryscripts/targetnameoption"

	"github.com/bastionzero/terraform-provider-bastionzero/internal"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// adBashModel maps the bash autodiscovery schema data.
type adBashModel struct {
	TargetNameOption types.String `tfsdk:"target_name_option"`
	EnvironmentID    types.String `tfsdk:"environment_id"`
	Script           types.String `tfsdk:"script"`
	ID               types.String `tfsdk:"id"`
}

func makeAdBashModelDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"target_name_option": schema.StringAttribute{
			Required:    true,
			Description: fmt.Sprintf("The target name schema option to use during autodiscovery %s.", internal.PrettyOneOf(targetnameoption.TargetNameOptionValues())),
			Validators: []validator.String{
				stringvalidator.OneOf(bastionzero.ToStringSlice(targetnameoption.TargetNameOptionValues())...),
			},
		},
		"environment_id": schema.StringAttribute{
			Required:    true,
			Description: "The unique environment ID the target should associate with.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"script": schema.StringAttribute{
			Computed:    true,
			Description: "Bash script that can be used to autodiscover a target.",
		},
		// Dummy "id" attribute. Required in order to test this data source.
		//
		// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84
		// Source: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
		// Source: https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
		"id": schema.StringAttribute{
			Computed:           true,
			Description:        "Deprecated. Do not depend on this attribute. This attribute will be removed in the future.",
			DeprecationMessage: "Do not depend on this attribute. This attribute will be removed in the future.",
		},
	}
}

func NewAdBashDataSource() datasource.DataSource {
	baseDesc := "Get a bash script that can be used to install the latest production BastionZero agent (bzero) on your targets."
	return bzdatasource.NewSingleDataSource(
		&bzdatasource.SingleDataSourceConfig[adBashModel, autodiscoveryscripts.BzeroBashAutodiscoveryScript]{
			BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[adBashModel, autodiscoveryscripts.BzeroBashAutodiscoveryScript]{
				RecordSchema:        makeAdBashModelDataSourceSchema(),
				MetadataTypeName:    "ad_bash",
				PrettyAttributeName: "autodiscovery script (bash)",
				FlattenAPIModel: func(ctx context.Context, apiObject *autodiscoveryscripts.BzeroBashAutodiscoveryScript, state *adBashModel) (diags diag.Diagnostics) {
					state.Script = types.StringValue(apiObject.Script)
					return
				},
				GetAPIModel: func(ctx context.Context, tfModel adBashModel, client *bastionzero.Client) (*autodiscoveryscripts.BzeroBashAutodiscoveryScript, error) {
					script, _, err := client.AutodiscoveryScripts.GetBzeroBashAutodiscoveryScript(ctx, &autodiscoveryscripts.BzeroBashAutodiscoveryOptions{
						TargetNameOption: targetnameoption.TargetNameOption(tfModel.TargetNameOption.ValueString()),
						EnvironmentID:    tfModel.EnvironmentID.ValueString(),
					})
					return script, err
				},
				Description: baseDesc,
				MarkdownDescription: baseDesc +
					"\n\nThe data source's `script` does not contain the registration secret that is required to register your targets with BastionZero. " +
					"You must replace `<REGISTRATION-SECRET-GOES-HERE>` with " +
					"a valid registration secret before attempting to execute the script.",
			},
		},
	)
}
