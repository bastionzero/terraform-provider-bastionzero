package bzerotarget

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource                     = &bzeroTargetDataSource{}
	_ datasource.DataSourceWithConfigure        = &bzeroTargetDataSource{}
	_ datasource.DataSourceWithConfigValidators = &bzeroTargetDataSource{}
)

type bzeroTargetDataSource struct {
	datasource.DataSourceWithConfigure
}

func (*bzeroTargetDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		// Validate only one of the schema defined attributes named id and name
		// has a known, non-null value.
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func NewBzeroTargetDataSource() datasource.DataSource {
	baseDesc := "Get information about a specific Bzero target in your BastionZero organization."
	return &bzeroTargetDataSource{
		DataSourceWithConfigure: bzdatasource.NewSingleDataSourceWithTimeout(
			&bzdatasource.SingleDataSourceWithTimeoutConfig[bzeroTargetModel, targets.BzeroTarget]{
				BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[bzeroTargetModel, targets.BzeroTarget]{
					RecordSchema: makeBzeroTargetDataSourceSchema(
						&target.BaseTargetDataSourceAttributeOptions{
							IsIDComputed:   true,
							IsNameComputed: true,
							IsIDOptional:   true,
							IsNameOptional: true,
						}),
					MetadataTypeName:    "bzero_target",
					PrettyAttributeName: "Bzero target",
					FlattenAPIModel: func(ctx context.Context, apiObject *targets.BzeroTarget, state *bzeroTargetModel) (diags diag.Diagnostics) {
						setBzeroTargetAttributes(ctx, state, apiObject)
						return
					},
					GetAPIModel: func(ctx context.Context, tfModel bzeroTargetModel, client *bastionzero.Client) (*targets.BzeroTarget, error) {
						if !tfModel.ID.IsNull() {
							// ID provided. Use GET API for single target with
							// ID.
							target, _, err := client.Targets.GetBzeroTarget(ctx, tfModel.ID.ValueString())
							return target, err
						} else if !tfModel.Name.IsNull() {
							// Name provided. List targets and find target with
							// specified name.
							targets, _, err := client.Targets.ListBzeroTargets(ctx)
							if err != nil {
								return nil, err
							}

							return findBzeroTargetByName(targets, tfModel.Name.ValueString())
						}

						// This should never happen due to
						// ConfigValidator.ExactlyOneOf
						panic("Expected one of \"id\" or \"name\" to be set. Please report this issue to the provider developers.")
					},
					Description:         baseDesc,
					MarkdownDescription: target.TargetDataSourceWithTimeoutMarkdownDescription(baseDesc, targettype.Bzero),
				},
				DefaultTimeout: 15 * time.Minute,
			},
		),
	}
}

func findBzeroTargetByName(targetList []targets.BzeroTarget, name string) (*targets.BzeroTarget, error) {
	results := make([]targets.BzeroTarget, 0)
	for _, target := range targetList {
		if target.Name == name {
			results = append(results, target)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("No bzero target found with name %s", name)
	}
	return nil, &backoff.PermanentError{Err: fmt.Errorf("Too many bzero targets found with name %s (found %d, expected 1)", name, len(results))}
}
