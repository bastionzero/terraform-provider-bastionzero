package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/sessionrecording"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tftypesgen"
	"pgregory.net/rapid"
)

func SessionRecordingPolicySchemaGen(ctx context.Context) *rapid.Generator[sessionrecording.SessionRecordingPolicyModel] {
	return rapid.Custom(func(t *rapid.T) sessionrecording.SessionRecordingPolicyModel {
		return sessionrecording.SessionRecordingPolicyModel{
			ID:          tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:        tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:        tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description: tftypesgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:    tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:      tftypesgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			RecordInput: tftypesgen.BoolWithValueOrNullGen(ctx).Draw(t, "RecordInputOrNull"),
		}
	})
}
