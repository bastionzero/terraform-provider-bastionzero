package clustertarget

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/bzdatasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewClusterTargetsDataSource() datasource.DataSource {
	return bzdatasource.NewListDataSource(&bzdatasource.ListDataSourceConfig[clusterTargetModel, targets.ClusterTarget]{
		BaseListDataSourceConfig: &bzdatasource.BaseListDataSourceConfig[clusterTargetModel, targets.ClusterTarget]{
			RecordSchema: makeClusterTargetDataSourceSchema(
				&target.BaseTargetDataSourceAttributeOptions{
					IsIDComputed:   true,
					IsNameComputed: true,
				}),
			ResultAttributeName: "cluster_targets",
			PrettyAttributeName: "Cluster targets",
			FlattenAPIModel: func(ctx context.Context, apiObject *targets.ClusterTarget) (state *clusterTargetModel, diags diag.Diagnostics) {
				state = new(clusterTargetModel)
				setClusterTargetAttributes(ctx, state, apiObject)
				return
			},
			Description: "Get a list of all Cluster targets in your BastionZero organization.",
		},
		ListAPIModels: func(ctx context.Context, client *bastionzero.Client) ([]targets.ClusterTarget, error) {
			targets, _, err := client.Targets.ListClusterTargets(ctx)
			return targets, err
		},
	})
}
