package bzpolicygen

import (
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"pgregory.net/rapid"
)

func PolicyClusterGen() *rapid.Generator[policies.Cluster] {
	return rapid.Custom(func(t *rapid.T) policies.Cluster {
		return policies.Cluster{
			ID: rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
		}
	})
}

func PolicyClusterUserGen() *rapid.Generator[policies.ClusterUser] {
	return rapid.Custom(func(t *rapid.T) policies.ClusterUser {
		return policies.ClusterUser{
			Name: rapid.String().Draw(t, "Name"),
		}
	})
}

func PolicyClusterGroupGen() *rapid.Generator[policies.ClusterGroup] {
	return rapid.Custom(func(t *rapid.T) policies.ClusterGroup {
		return policies.ClusterGroup{
			Name: rapid.String().Draw(t, "Name"),
		}
	})
}

func KubernetesPolicyGen() *rapid.Generator[policies.KubernetesPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.KubernetesPolicy {
		return policies.KubernetesPolicy{
			Policy:        PolicyGen().Draw(t, "BasePolicy"),
			Environments:  rapid.Ptr(rapid.SliceOf(PolicyEnvironmentGen()), false).Draw(t, "Environments"),
			Clusters:      rapid.Ptr(rapid.SliceOf(PolicyClusterGen()), false).Draw(t, "Clusters"),
			ClusterUsers:  rapid.Ptr(rapid.SliceOf(PolicyClusterUserGen()), false).Draw(t, "ClusterUsers"),
			ClusterGroups: rapid.Ptr(rapid.SliceOf(PolicyClusterGroupGen()), false).Draw(t, "ClusterGroups"),
		}
	})
}
