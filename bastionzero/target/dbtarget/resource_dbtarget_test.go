package dbtarget_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/targetstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target/dbtarget"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDbTarget_Basic(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			// Verify create works for a config set with all required attributes
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env.ID,
						Name:               &rName,
						ProxyTargetID:      &bzeroTarget.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
						DatabaseAuthConfig: &dbauthconfig.DatabaseAuthenticationConfig{AuthenticationType: bastionzero.PtrTo(dbauthconfig.Default), Label: bastionzero.PtrTo("None")},
						LocalPort:          nil,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					// Check the state value we explicitly configured in this
					// test is correct
					resource.TestCheckResourceAttr(resourceName, "environment_id", env.ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "proxy_target_id", bzeroTarget.ID),
					resource.TestCheckResourceAttr(resourceName, "remote_host", "localhost"),
					resource.TestCheckResourceAttr(resourceName, "remote_port", "5432"),
					// Check default values are set in state
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.authentication_type", dbauthconfig.Default),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.label", "None"),
					resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.cloud_service_provider"),
					resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.database"),
					// Check that unspecified values remain null
					resource.TestCheckNoResourceAttr(resourceName, "local_port"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDbTarget_Disappears(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					acctest.CheckResourceDisappears(resourceName, func(c *bastionzero.Client, ctx context.Context, id string) (*http.Response, error) {
						return c.Targets.DeleteDatabaseTarget(ctx, id)
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDbTarget_DatabaseAuthConfig(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	dbAuthConfig1 := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType: bastionzero.PtrTo(dbauthconfig.Default),
		Database:           bastionzero.PtrTo(dbauthconfig.Postgres),
		Label:              bastionzero.PtrTo("Postgres"),
	}
	dbAuthConfig2 := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType: bastionzero.PtrTo(dbauthconfig.SplitCert),
		Database:           bastionzero.PtrTo(dbauthconfig.MongoDB),
		Label:              bastionzero.PtrTo("MongoDB"),
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigDbAuthConfig(rName, env.ID, bzeroTarget.ID, "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(ctx, dbAuthConfig1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env.ID,
						Name:               &rName,
						ProxyTargetID:      &bzeroTarget.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
						DatabaseAuthConfig: dbAuthConfig1,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.authentication_type", *dbAuthConfig1.AuthenticationType),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.database", *dbAuthConfig1.Database),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.label", *dbAuthConfig1.Label),
					resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.cloud_service_provider"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update db auth config
			{
				Config: testAccDbTargetConfigDbAuthConfig(rName, env.ID, bzeroTarget.ID, "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(ctx, dbAuthConfig2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env.ID,
						Name:               &rName,
						ProxyTargetID:      &bzeroTarget.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
						DatabaseAuthConfig: dbAuthConfig2,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.authentication_type", *dbAuthConfig2.AuthenticationType),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.database", *dbAuthConfig2.Database),
					resource.TestCheckResourceAttr(resourceName, "database_authentication_config.label", *dbAuthConfig2.Label),
					resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.cloud_service_provider"),
				),
			},
		},
	})
}

func testAccDbTargetConfigBasic(name string, envID string, proxyTargetID string, remoteHost string, remotePort string) string {
	return fmt.Sprintf(`
resource "bastionzero_db_target" "test" {
  environment_id = %[2]q
  name = %[1]q
  proxy_target_id = %[3]q
  remote_host = %[4]q
  remote_port = %[5]q
}
`, name, envID, proxyTargetID, remoteHost, remotePort)
}

func testAccDbTargetConfigDbAuthConfig(name string, envID string, proxyTargetID string, remoteHost string, remotePort string, dbAuthConfig types.Object) string {
	return fmt.Sprintf(`
resource "bastionzero_db_target" "test" {
  environment_id = %[2]q
  name = %[1]q
  proxy_target_id = %[3]q
  remote_host = %[4]q
  remote_port = %[5]q
  database_authentication_config = %[6]s
}
`, name, envID, proxyTargetID, remoteHost, remotePort, acctest.TerraformObjectToString(dbAuthConfig))
}

type expectedDbTarget struct {
	EnvironmentID      *string
	Name               *string
	ProxyTargetID      *string
	RemoteHost         *string
	RemotePort         *int
	DatabaseAuthConfig *dbauthconfig.DatabaseAuthenticationConfig
	LocalPort          *int
}

func testAccCheckDbTargetAttributes(t *testing.T, target *targets.DatabaseTarget, expected *expectedDbTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if expected.EnvironmentID != nil && *expected.EnvironmentID != target.EnvironmentID {
			return fmt.Errorf("Bad environment_id, expected \"%s\", got: %#v", *expected.EnvironmentID, target.EnvironmentID)
		}
		if expected.Name != nil && *expected.Name != target.Name {
			return fmt.Errorf("Bad name, expected \"%s\", got: %#v", *expected.Name, target.Name)
		}
		if expected.ProxyTargetID != nil && *expected.ProxyTargetID != target.ProxyTargetID {
			return fmt.Errorf("Bad proxy_target_id, expected \"%s\", got: %#v", *expected.ProxyTargetID, target.ProxyTargetID)
		}
		if expected.RemoteHost != nil && *expected.RemoteHost != target.RemoteHost {
			return fmt.Errorf("Bad remote_host, expected \"%s\", got: %#v", *expected.RemoteHost, target.RemoteHost)
		}
		if expected.RemotePort != nil && !assert.Equal(t, expected.RemotePort, target.RemotePort.Value) {
			return fmt.Errorf("Bad remote_port, expected \"%s\", got: %s", acctest.SafePrettyInt(expected.RemotePort), acctest.SafePrettyInt(target.RemotePort.Value))
		}
		if expected.DatabaseAuthConfig != nil && !assert.Equal(t, *expected.DatabaseAuthConfig, target.DatabaseAuthenticationConfig) {
			return fmt.Errorf("Bad database_authentication_config, expected \"%#v\", got: %#v", *expected.DatabaseAuthConfig, target.DatabaseAuthenticationConfig)
		}
		if !assert.Equal(t, expected.LocalPort, target.LocalPort.Value) {
			return fmt.Errorf("Bad local_port, expected \"%s\", got: %s", acctest.SafePrettyInt(expected.LocalPort), acctest.SafePrettyInt(target.LocalPort.Value))
		}

		return nil
	}
}

func testAccCheckDbTargetExists(namedTFResource string, target *targets.DatabaseTarget) resource.TestCheckFunc {
	return acctest.CheckExistsAtBastionZero(namedTFResource, target, func(c *bastionzero.Client, ctx context.Context, id string) (*targets.DatabaseTarget, *http.Response, error) {
		return c.Targets.GetDatabaseTarget(ctx, id)
	})
}

func testAccCheckResourceDbTargetComputedAttr(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "agent_public_key"),
		resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestMatchResourceAttr(resourceName, "last_agent_update", regexp.MustCompile(acctest.RFC3339RegexPattern)),
		resource.TestCheckResourceAttrSet(resourceName, "region"),
		acctest.CheckResourceAttrIsOneOf(resourceName, "status", bastionzero.ToStringSlice(targetstatus.TargetStatusValues())),
		resource.TestCheckResourceAttr(resourceName, "type", string(targettype.Db)),
	)
}

func testAccCheckDbTargetDestroy(s *terraform.State) error {
	return acctest.CheckAllResourcesWithTypeDestroyed(
		"bastionzero_db_target",
		func(client *bastionzero.Client, ctx context.Context, id string) (*targets.DatabaseTarget, *http.Response, error) {
			return client.Targets.GetDatabaseTarget(ctx, id)
		},
	)(s)
}