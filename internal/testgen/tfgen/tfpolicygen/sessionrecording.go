package tfpolicygen

import (
	"context"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies/policytype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/sessionrecording"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/tfgen"
	"pgregory.net/rapid"
)

func SessionRecordingPolicySchemaGen(ctx context.Context) *rapid.Generator[sessionrecording.SessionRecordingPolicyModel] {
	return rapid.Custom(func(t *rapid.T) sessionrecording.SessionRecordingPolicyModel {
		return sessionrecording.SessionRecordingPolicyModel{
			ID:          tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaIDGen()).Draw(t, "IDOrNull"),
			Name:        tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaNameGen()).Draw(t, "NameOrNull"),
			Type:        tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaTypeGen(policytype.Kubernetes)).Draw(t, "TypeOrNull"),
			Description: tfgen.StringWithValueOrNullOrEmptyGen(PolicySchemaDescriptionGen()).Draw(t, "DescriptionOrNull"),
			Subjects:    tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaSubjectsGen()).Draw(t, "SubjectsOrNull"),
			Groups:      tfgen.SetWithValueOrNullOrEmptyGen(ctx, PolicySchemaGroupsGen()).Draw(t, "GroupsOrNull"),
			RecordInput: tfgen.BoolWithValueOrNullGen(ctx).Draw(t, "RecordInputOrNull"),
		}
	})
}
