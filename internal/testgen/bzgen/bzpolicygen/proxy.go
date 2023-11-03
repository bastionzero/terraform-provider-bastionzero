package bzpolicygen

import (
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"pgregory.net/rapid"
)

func ProxyPolicyGen() *rapid.Generator[policies.ProxyPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.ProxyPolicy {
		return policies.ProxyPolicy{
			Policy:       PolicyGen().Draw(t, "BasePolicy"),
			Environments: rapid.Ptr(rapid.SliceOf(PolicyEnvironmentGen()), true).Draw(t, "Environments"),
			Targets:      rapid.Ptr(rapid.SliceOf(PolicyTargetGen()), true).Draw(t, "Targets"),
			TargetUsers:  rapid.Ptr(rapid.SliceOf(PolicyTargetUserGen()), true).Draw(t, "TargetUsers"),
		}
	})
}
