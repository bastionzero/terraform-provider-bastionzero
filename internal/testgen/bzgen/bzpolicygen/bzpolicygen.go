// Package bzpolicygen provides property based test (PBT) generators for
// BastionZero policy API types
package bzpolicygen

import (
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"pgregory.net/rapid"
)

func PolicySubjectGen() *rapid.Generator[policies.Subject] {
	return rapid.Custom(func(t *rapid.T) policies.Subject {
		return policies.Subject{
			ID:   rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			Type: subjecttype.SubjectType(rapid.SampledFrom(subjecttype.SubjectTypeValues()).Draw(t, "Type")),
		}
	})
}

func PolicyGroupGen() *rapid.Generator[policies.Group] {
	return rapid.Custom(func(t *rapid.T) policies.Group {
		return policies.Group{
			ID:   rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			Name: rapid.String().Draw(t, "Name"),
		}
	})
}

func PolicyEnvironmentGen() *rapid.Generator[policies.Environment] {
	return rapid.Custom(func(t *rapid.T) policies.Environment {
		return policies.Environment{
			ID: rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
		}
	})
}

func PolicyTargetGen() *rapid.Generator[policies.Target] {
	return rapid.Custom(func(t *rapid.T) policies.Target {
		return policies.Target{
			ID:   rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			Type: targettype.TargetType(rapid.SampledFrom(targettype.TargetTypeValues()).Draw(t, "Type")),
		}
	})
}

func PolicyTargetUserGen() *rapid.Generator[policies.TargetUser] {
	return rapid.Custom(func(t *rapid.T) policies.TargetUser {
		return policies.TargetUser{
			Username: rapid.String().Draw(t, "Username"),
		}
	})
}

func PolicyGen() *rapid.Generator[policies.Policy] {
	return rapid.Custom(func(t *rapid.T) policies.Policy {
		return policies.Policy{
			ID:          rapid.StringMatching(acctest.UUIDV4RegexPattern).Draw(t, "ID"),
			Description: rapid.Ptr(rapid.String(), true).Draw(t, "Description"),
			Name:        rapid.String().Draw(t, "Name"),
			TimeExpires: &types.Timestamp{Time: time.Now().Add(time.Duration(rapid.Int64().Draw(t, "TimeExpires")))},
			Subjects:    rapid.Ptr(rapid.SliceOf(PolicySubjectGen()), true).Draw(t, "Subjects"),
			Groups:      rapid.Ptr(rapid.SliceOf(PolicyGroupGen()), true).Draw(t, "Groups"),
		}
	})
}
