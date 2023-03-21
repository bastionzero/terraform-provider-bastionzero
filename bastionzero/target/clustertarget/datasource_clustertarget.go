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
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			&bzdatasource.SingleDataSourceConfigWithTimeout[clusterTargetModel, targets.ClusterTarget]{
				BaseSingleDataSourceConfig: &bzdatasource.BaseSingleDataSourceConfig[clusterTargetModel, targets.ClusterTarget]{
					RecordSchema: makeClusterTargetDataSourceSchema(
						&target.BaseTargetDataSourceAttributeOptions{
							IsIDComputed:   true,
							IsNameComputed: true,
							IsIDOptional:   true,
							IsNameOptional: true,
						}),
					ResultAttributeName: "cluster_target",
					PrettyAttributeName: "Cluster target",
					FlattenAPIModel: func(ctx context.Context, apiObject *targets.ClusterTarget, _ clusterTargetModel) (state *clusterTargetModel, diags diag.Diagnostics) {
						state = new(clusterTargetModel)
						setClusterTargetAttributes(ctx, state, apiObject)
						return
					},
					Description:         baseDesc,
					MarkdownDescription: target.TargetDataSourceWithTimeoutMarkdownDescription(baseDesc, targettype.Cluster),
				},
				DefaultTimeout: 15 * time.Minute,
				GetAPIModelWithTimeout: func(ctx context.Context, tfModel clusterTargetModel, client *bastionzero.Client, timeout time.Duration) (*targets.ClusterTarget, error) {
					// An operation that may fail.
					operation := func() (*targets.ClusterTarget, error) {
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
					}

					// Attempt to find target with backoff
					backOffConfig := backoff.NewExponentialBackOff()
					backOffConfig.MaxElapsedTime = timeout
					target, err := backoff.RetryNotifyWithData(operation, backoff.WithContext(backOffConfig, ctx), func(err error, dur time.Duration) {
						tflog.Info(ctx, fmt.Sprintf("%v. Retrying in %s...", err, dur))
					})

					// We timed out, or a backoff.PermanentError was returned
					if err != nil {
						return nil, err
					}

					// Operation is successful.
					return target, nil
				},
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
