package clustertarget

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
	_ datasource.DataSource                     = &clusterTargetDataSource{}
	_ datasource.DataSourceWithConfigure        = &clusterTargetDataSource{}
	_ datasource.DataSourceWithConfigValidators = &clusterTargetDataSource{}
)

type clusterTargetDataSource struct {
	datasource.DataSourceWithConfigure
}

func (*clusterTargetDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		// Validate only one of the schema defined attributes named id and name
		// has a known, non-null value.
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func NewClusterTargetDataSource() datasource.DataSource {
	baseDesc := "Get information about a specific Cluster target in your BastionZero organization."
	return &clusterTargetDataSource{
		DataSourceWithConfigure: bzdatasource.NewSingleDataSourceWithTimeout(
			&bzdatasource.SingleDataSourceWithTimeoutConfig[clusterTargetModel, targets.ClusterTarget]{
				BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[clusterTargetModel, targets.ClusterTarget]{
					RecordSchema: makeClusterTargetDataSourceSchema(
						&target.BaseTargetDataSourceAttributeOptions{
							IsIDComputed:   true,
							IsNameComputed: true,
							IsIDOptional:   true,
							IsNameOptional: true,
						}),
					MetadataTypeName:    "cluster_target",
					PrettyAttributeName: "Cluster target",
					FlattenAPIModel: func(ctx context.Context, apiObject *targets.ClusterTarget, state *clusterTargetModel) (diags diag.Diagnostics) {
						setClusterTargetAttributes(ctx, state, apiObject)
						return
					},
					GetAPIModel: func(ctx context.Context, tfModel clusterTargetModel, client *bastionzero.Client) (*targets.ClusterTarget, error) {
						if !tfModel.ID.IsNull() {
							// ID provided. Use GET API for single target with ID.
							target, _, err := client.Targets.GetClusterTarget(ctx, tfModel.ID.ValueString())
							return target, err
						} else if !tfModel.Name.IsNull() {
							// Name provided. List targets and find target with
							// specified name.
							targets, _, err := client.Targets.ListClusterTargets(ctx)
							if err != nil {
								return nil, err
							}

							return findClusterTargetByName(targets, tfModel.Name.ValueString())
						}

						// This should never happen due to
						// ConfigValidator.ExactlyOneOf
						panic("Expected one of \"id\" or \"name\" to be set. Please report this issue to the provider developers.")
					},
					Description:         baseDesc,
					MarkdownDescription: target.TargetDataSourceWithTimeoutMarkdownDescription(baseDesc, targettype.Cluster),
				},
				DefaultTimeout: 15 * time.Minute,
			},
		),
	}
}

func findClusterTargetByName(targetList []targets.ClusterTarget, name string) (*targets.ClusterTarget, error) {
	results := make([]targets.ClusterTarget, 0)
	for _, target := range targetList {
		if target.Name == name {
			results = append(results, target)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("No cluster target found with name %s", name)
	}
	return nil, &backoff.PermanentError{Err: fmt.Errorf("Too many cluster targets found with name %s (found %d, expected 1)", name, len(results))}
}
