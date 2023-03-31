package policy_test

import (
	"context"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestFlatExpandPolicySubjects(t *testing.T) {
	elementType := policy.GetPolicySubjectModelType(context.Background())
	attributeTypes := elementType.AttrTypes

	cases := []struct {
		setSubjects basetypes.SetValue
		expected    []policies.Subject
	}{
		// simple
		{
			setSubjects: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue("type"),
				}),
			}),
			expected: []policies.Subject{
				{ID: "id", Type: "type"},
			},
		},
		// many
		{
			setSubjects: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue("type"),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"type": basetypes.NewStringValue("type2"),
				}),
			}),
			expected: []policies.Subject{
				{ID: "id", Type: "type"},
				{ID: "id2", Type: "type2"},
			},
		},
		// missing field value
		{
			setSubjects: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue(""),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"type": basetypes.NewStringValue(""),
				}),
			}),
			expected: []policies.Subject{
				{ID: "id", Type: ""},
				{ID: "id2", Type: ""},
			},
		},
	}

	for _, c := range cases {
		// Expand
		gotExpand := policy.ExpandPolicySubjects(context.Background(), c.setSubjects)
		require.EqualValues(t, c.expected, gotExpand)

		// Flatten back
		gotFlatten := policy.FlattenPolicySubjects(context.Background(), gotExpand)
		require.EqualValues(t, c.setSubjects, gotFlatten)
	}
}

func TestFlatExpandPolicyGroups(t *testing.T) {
	elementType := policy.GetPolicyGroupModelType(context.Background())
	attributeTypes := elementType.AttrTypes

	cases := []struct {
		setGroups basetypes.SetValue
		expected  []policies.Group
	}{
		// simple
		{
			setGroups: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"name": basetypes.NewStringValue("name"),
				}),
			}),
			expected: []policies.Group{
				{ID: "id", Name: "name"},
			},
		},
		// many
		{
			setGroups: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"name": basetypes.NewStringValue("name"),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"name": basetypes.NewStringValue("name2"),
				}),
			}),
			expected: []policies.Group{
				{ID: "id", Name: "name"},
				{ID: "id2", Name: "name2"},
			},
		},
		// missing field value
		{
			setGroups: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"name": basetypes.NewStringValue(""),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"name": basetypes.NewStringValue(""),
				}),
			}),
			expected: []policies.Group{
				{ID: "id", Name: ""},
				{ID: "id2", Name: ""},
			},
		},
	}

	for _, c := range cases {
		// Expand
		gotExpand := policy.ExpandPolicyGroups(context.Background(), c.setGroups)
		require.EqualValues(t, c.expected, gotExpand)

		// Flatten back
		gotFlatten := policy.FlattenPolicyGroups(context.Background(), gotExpand)
		require.EqualValues(t, c.setGroups, gotFlatten)
	}
}

func TestFlatExpandEnvironments(t *testing.T) {
	elementType := basetypes.StringType{}

	cases := []struct {
		setEnvs  basetypes.SetValue
		expected []policies.Environment
	}{
		// simple
		{
			setEnvs: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("env"),
			}),
			expected: []policies.Environment{
				{ID: "env"},
			},
		},
		// many
		{
			setEnvs: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("env1"),
				types.StringValue("env2"),
			}),
			expected: []policies.Environment{
				{ID: "env1"},
				{ID: "env2"},
			},
		},
		// missing field value
		{
			setEnvs: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("env1"),
				types.StringValue(""),
			}),
			expected: []policies.Environment{
				{ID: "env1"},
				{ID: ""},
			},
		},
	}

	for _, c := range cases {
		// Expand
		gotExpand := policy.ExpandPolicyEnvironments(context.Background(), c.setEnvs)
		require.EqualValues(t, c.expected, gotExpand)

		// Flatten back
		gotFlatten := policy.FlattenPolicyEnvironments(context.Background(), gotExpand)
		require.EqualValues(t, c.setEnvs, gotFlatten)
	}
}

func TestFlatExpandPolicyTargets(t *testing.T) {
	elementType := policy.GetPolicyTargetModelType(context.Background())
	attributeTypes := elementType.AttrTypes

	cases := []struct {
		setTargets basetypes.SetValue
		expected   []policies.Target
	}{
		// simple
		{
			setTargets: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue("type"),
				}),
			}),
			expected: []policies.Target{
				{ID: "id", Type: "type"},
			},
		},
		// many
		{
			setTargets: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue("type"),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"type": basetypes.NewStringValue("type2"),
				}),
			}),
			expected: []policies.Target{
				{ID: "id", Type: "type"},
				{ID: "id2", Type: "type2"},
			},
		},
		// missing field value
		{
			setTargets: basetypes.NewSetValueMust(elementType, []attr.Value{
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id"),
					"type": basetypes.NewStringValue(""),
				}),
				basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
					"id":   basetypes.NewStringValue("id2"),
					"type": basetypes.NewStringValue(""),
				}),
			}),
			expected: []policies.Target{
				{ID: "id", Type: ""},
				{ID: "id2", Type: ""},
			},
		},
	}

	for _, c := range cases {
		// Expand
		gotExpand := policy.ExpandPolicyTargets(context.Background(), c.setTargets)
		require.EqualValues(t, c.expected, gotExpand)

		// Flatten back
		gotFlatten := policy.FlattenPolicyTargets(context.Background(), gotExpand)
		require.EqualValues(t, c.setTargets, gotFlatten)
	}
}

func TestFlatExpandTargetUsers(t *testing.T) {
	elementType := basetypes.StringType{}

	cases := []struct {
		setTargetUsers basetypes.SetValue
		expected       []policies.TargetUser
	}{
		// simple
		{
			setTargetUsers: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("user"),
			}),
			expected: []policies.TargetUser{
				{Username: "user"},
			},
		},
		// many
		{
			setTargetUsers: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("user1"),
				types.StringValue("user2"),
			}),
			expected: []policies.TargetUser{
				{Username: "user1"},
				{Username: "user2"},
			},
		},
		// missing field value
		{
			setTargetUsers: basetypes.NewSetValueMust(elementType, []attr.Value{
				types.StringValue("user1"),
				types.StringValue(""),
			}),
			expected: []policies.TargetUser{
				{Username: "user1"},
				{Username: ""},
			},
		},
	}

	for _, c := range cases {
		// Expand
		gotExpand := policy.ExpandPolicyTargetUsers(context.Background(), c.setTargetUsers)
		require.EqualValues(t, c.expected, gotExpand)

		// Flatten back
		gotFlatten := policy.FlattenPolicyTargetUsers(context.Background(), gotExpand)
		require.EqualValues(t, c.setTargetUsers, gotFlatten)
	}
}
