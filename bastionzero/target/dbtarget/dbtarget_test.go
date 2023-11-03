package dbtarget_test

import (
	"context"
	"testing"

	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target/dbtarget"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/testgen/bzgen"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestFlatExpandDatabaseAuthenticationConfig_NoDataLoss(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		genAPI := bzgen.DatabaseAuthenticationConfigGen().Draw(t, "DatabaseAuthenticationConfig")

		// Flatten the generated BastionZero type into a TF type
		flattened := dbtarget.FlattenDatabaseAuthenticationConfig(context.Background(), &genAPI)

		// Then expand the value back into a BastionZero API type
		expanded := dbtarget.ExpandDatabaseAuthenticationConfig(context.Background(), flattened)

		// And assert no data loss occurred
		require.EqualValues(t, &genAPI, expanded)
	})
}
