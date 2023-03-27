package environment

import (
	"context"
	"testing"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestSetTFEnvironmentAttributes(t *testing.T) {
	now := time.Now()
	elementType := getEnvironmentTargetModelType(context.Background())
	attributeTypes := elementType.AttrTypes

	cases := []struct {
		env      *environments.Environment
		expected *environmentModel
	}{
		// simple, no targets
		{
			env: &environments.Environment{
				ID:                         "id",
				OrganizationID:             "orgId",
				IsDefault:                  false,
				Name:                       "name",
				Description:                "desc",
				TimeCreated:                types.Timestamp{Time: now},
				OfflineCleanupTimeoutHours: 30,
				Targets:                    make([]environments.TargetSummary, 0),
			},
			expected: &environmentModel{
				ID:                         basetypes.NewStringValue("id"),
				OrganizationID:             basetypes.NewStringValue("orgId"),
				IsDefault:                  basetypes.NewBoolValue(false),
				Name:                       basetypes.NewStringValue("name"),
				Description:                basetypes.NewStringValue("desc"),
				TimeCreated:                basetypes.NewStringValue(now.UTC().Format(time.RFC3339)),
				OfflineCleanupTimeoutHours: basetypes.NewInt64Value(30),
				Targets:                    basetypes.NewMapValueMust(getEnvironmentTargetModelType(context.Background()), make(map[string]attr.Value)),
			},
		},
		// with targets
		{
			env: &environments.Environment{
				ID:                         "id",
				OrganizationID:             "orgId",
				IsDefault:                  false,
				Name:                       "name",
				Description:                "desc",
				TimeCreated:                types.Timestamp{Time: now},
				OfflineCleanupTimeoutHours: 30,
				Targets:                    []environments.TargetSummary{{ID: "foo", Type: targettype.Bzero}},
			},
			expected: &environmentModel{
				ID:                         basetypes.NewStringValue("id"),
				OrganizationID:             basetypes.NewStringValue("orgId"),
				IsDefault:                  basetypes.NewBoolValue(false),
				Name:                       basetypes.NewStringValue("name"),
				Description:                basetypes.NewStringValue("desc"),
				TimeCreated:                basetypes.NewStringValue(now.UTC().Format(time.RFC3339)),
				OfflineCleanupTimeoutHours: basetypes.NewInt64Value(30),
				Targets: basetypes.NewMapValueMust(elementType, map[string]attr.Value{
					"foo": basetypes.NewObjectValueMust(attributeTypes, map[string]attr.Value{
						"id":   basetypes.NewStringValue("foo"),
						"type": basetypes.NewStringValue(string(targettype.Bzero)),
					}),
				}),
			},
		},
	}

	for _, c := range cases {
		got := new(environmentModel)
		setEnvironmentAttributes(context.Background(), got, c.env)
		require.EqualValues(t, c.expected, got)
	}
}
