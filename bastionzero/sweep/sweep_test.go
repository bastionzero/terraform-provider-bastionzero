package sweep_test

import (
	"testing"

	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
