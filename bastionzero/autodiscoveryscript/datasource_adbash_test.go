package autodiscoveryscript_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/autodiscoveryscripts/targetnameoption"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccADBashDataSource_Basic(t *testing.T) {
	ctx := context.Background()
	dataSourceName := "data.bastionzero_ad_bash.test"
	env := new(policies.Environment)

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)
	acctest.FindNEnvironmentsOrSkipAsPolicyEnvironment(t, env)

	makeTestStep := func(targetNameOption targetnameoption.TargetNameOption) resource.TestStep {
		return resource.TestStep{
			Config: testAccADBashDataSourceConfigBasic(env.ID, string(targetNameOption)),
			Check:  resource.TestCheckResourceAttrSet(dataSourceName, "script"),
		}
	}
	// Make a step per valid targetNameOption
	var steps []resource.TestStep
	for _, opt := range targetnameoption.TargetNameOptionValues() {
		steps = append(steps, makeTestStep(opt))
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestADBashDataSource_InvalidEnvironmentID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty environment_id not permitted
				Config:      testAccADBashDataSourceConfigBasic("", string(targetnameoption.BashHostName)),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestADBashDataSource_InvalidTargetNameOption(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid target_name_option not permitted
				Config:      testAccADBashDataSourceConfigBasic("test", "foo"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func testAccADBashDataSourceConfigBasic(environmentID string, targetNameOption string) string {
	return fmt.Sprintf(`
data "bastionzero_ad_bash" "test" {
  environment_id = %[1]q
  target_name_option = %[2]q
}
`, environmentID, targetNameOption)
}
