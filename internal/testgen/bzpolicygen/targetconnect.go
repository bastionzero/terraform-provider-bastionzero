package bzpolicygen

import (
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/verbtype"
	"pgregory.net/rapid"
)

func PolicyVerbGen() *rapid.Generator[policies.Verb] {
	return rapid.Custom(func(t *rapid.T) policies.Verb {
		return policies.Verb{
			Type: verbtype.VerbType(rapid.SampledFrom(verbtype.VerbTypeValues()).Draw(t, "Type")),
		}
	})
}

func TargetConnectPolicyGen() *rapid.Generator[policies.TargetConnectPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.TargetConnectPolicy {
		return policies.TargetConnectPolicy{
			Policy:       PolicyGen().Draw(t, "BasePolicy"),
			Environments: rapid.Ptr(rapid.SliceOf(PolicyEnvironmentGen()), true).Draw(t, "Environments"),
			Targets:      rapid.Ptr(rapid.SliceOf(PolicyTargetGen()), true).Draw(t, "Targets"),
			TargetUsers:  rapid.Ptr(rapid.SliceOf(PolicyTargetUserGen()), true).Draw(t, "TargetUsers"),
			Verbs:        rapid.Ptr(rapid.SliceOf(PolicyVerbGen()), true).Draw(t, "Verbs"),
		}
	})
}
