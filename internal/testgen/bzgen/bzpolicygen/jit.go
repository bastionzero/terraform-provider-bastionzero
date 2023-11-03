package bzpolicygen

import (
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"pgregory.net/rapid"
)

func ChildPolicyGen() *rapid.Generator[policies.ChildPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.ChildPolicy {
		return policies.ChildPolicy{
			ID:   rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			Type: policytype.PolicyType(rapid.SampledFrom([]policytype.PolicyType{policytype.TargetConnect, policytype.Kubernetes, policytype.Proxy}).Draw(t, "Type")),
			Name: rapid.String().Draw(t, "Name"),
		}
	})
}

func JITPolicyGen() *rapid.Generator[policies.JITPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.JITPolicy {
		return policies.JITPolicy{
			ID:                    rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			TimeExpires:           &types.Timestamp{Time: time.Now().Add(time.Duration(rapid.Int64().Draw(t, "TimeExpires")))},
			Name:                  rapid.String().Draw(t, "Name"),
			Description:           rapid.String().Draw(t, "Description"),
			Subjects:              rapid.SliceOf(PolicySubjectGen()).Draw(t, "Subjects"),
			Groups:                rapid.SliceOf(PolicyGroupGen()).Draw(t, "Groups"),
			ChildPolicies:         rapid.SliceOf(ChildPolicyGen()).Draw(t, "ChildPolicies"),
			AutomaticallyApproved: rapid.Bool().Draw(t, "AutomaticallyApproved"),
			Duration:              rapid.Uint().Draw(t, "Duration"),
		}
	})
}
