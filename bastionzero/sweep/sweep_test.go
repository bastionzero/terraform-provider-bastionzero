package sweep_test

import (
	"testing"

	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/environment"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/jit"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/kubernetes"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/proxy"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/sessionrecording"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/policy/targetconnect"
	_ "github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target/dbtarget"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
