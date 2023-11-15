package dbtarget_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/environments"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/dbauthconfig"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/targets/targetstatus"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/targettype"
	"github.com/bastionzero/terraform-provider-bastionzero/bastionzero/target/dbtarget"
	"github.com/bastionzero/terraform-provider-bastionzero/internal/acctest"
	"github.com/google/uuid"
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

func TestAccDbTarget_EnvironmentID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env1 := new(environments.Environment)
	env2 := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env1, env2)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigBasic(rName, env1.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env1.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environment_id", env1.ID),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update environment
			{
				Config: testAccDbTargetConfigBasic(rName, env2.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env2.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "environment_id", env2.ID),
				),
			},
		},
	})
}

func TestAccDbTarget_Name(t *testing.T) {
	ctx := context.Background()
	rName1 := acctest.RandomName()
	rName2 := acctest.RandomName()
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
				Config: testAccDbTargetConfigBasic(rName1, env.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName1,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update name
			{
				Config: testAccDbTargetConfigBasic(rName2, env.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName2,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccDbTarget_ProxyTargetID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget1 := new(targets.BzeroTarget)
	bzeroTarget2 := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget1, bzeroTarget2)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget1.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget1.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_target_id", bzeroTarget1.ID),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update proxy target ID
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget2.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget2.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_target_id", bzeroTarget2.ID),
				),
			},
		},
	})
}

func TestAccDbTarget_ProxyEnvironmentID(t *testing.T) {
	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"
	var target targets.DatabaseTarget

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env1 := new(environments.Environment)
	env2 := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env1, env2)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigProxyEnvID(rName, env1.ID, env1.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env1.ID,
						Name:               &rName,
						ProxyEnvironmentID: &env1.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttrProxyEnvID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_environment_id", env1.ID),
					resource.TestCheckNoResourceAttr(resourceName, "proxy_target_id"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update proxy environment ID
			{
				Config: testAccDbTargetConfigProxyEnvID(rName, env1.ID, env2.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env1.ID,
						Name:               &rName,
						ProxyEnvironmentID: &env2.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttrProxyEnvID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_environment_id", env2.ID),
					resource.TestCheckNoResourceAttr(resourceName, "proxy_target_id"),
				),
			},
			// Verify switch from proxy environment ID --> proxy target ID
			{
				Config: testAccDbTargetConfigBasic(rName, env1.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env1.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_target_id", bzeroTarget.ID),
					resource.TestCheckNoResourceAttr(resourceName, "proxy_environment_id"),
				),
			},
			// And then back proxy target ID --> proxy environment ID
			{
				Config: testAccDbTargetConfigProxyEnvID(rName, env1.ID, env2.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID:      &env1.ID,
						Name:               &rName,
						ProxyEnvironmentID: &env2.ID,
						RemoteHost:         bastionzero.PtrTo("localhost"),
						RemotePort:         bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttrProxyEnvID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "proxy_environment_id", env2.ID),
					resource.TestCheckNoResourceAttr(resourceName, "proxy_target_id"),
				),
			},
		},
	})
}

func TestAccDbTarget_RemoteHost(t *testing.T) {
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
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_host", "localhost"),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update remote host
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost2", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost2"),
						RemotePort:    bastionzero.PtrTo(5432),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_host", "localhost2"),
				),
			},
		},
	})
}

func TestAccDbTarget_RemotePort(t *testing.T) {
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

	remotePort1 := 3000
	remotePort2 := 4000

	dbAuthConfig := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType:   bastionzero.PtrTo(dbauthconfig.ServiceAccountInjection),
		Database:             bastionzero.PtrTo(dbauthconfig.Postgres),
		Label:                bastionzero.PtrTo("GCP Postgres"),
		CloudServiceProvider: bastionzero.PtrTo(dbauthconfig.GCP),
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", strconv.Itoa(remotePort1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    &remotePort1,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_port", strconv.Itoa(remotePort1)),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update remote port
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", strconv.Itoa(remotePort2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    &remotePort2,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_port", strconv.Itoa(remotePort2)),
				),
			},
			// Verify setting remote port even when target is GCP
			{
				Config: testAccDbTargetConfigDbAuthConfig(rName, env.ID, bzeroTarget.ID, "gcp://localhost", strconv.Itoa(remotePort2), dbtarget.FlattenDatabaseAuthenticationConfig(ctx, dbAuthConfig)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("gcp://localhost"),
						RemotePort:    &remotePort2,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_port", strconv.Itoa(remotePort2)),
				),
			},
			// Also check that it works with 0 (the recommended way of
			// configuring remote_port for GCP DB target)
			{
				Config: testAccDbTargetConfigDbAuthConfig(rName, env.ID, bzeroTarget.ID, "gcp://localhost", strconv.Itoa(0), dbtarget.FlattenDatabaseAuthenticationConfig(ctx, dbAuthConfig)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("gcp://localhost"),
						RemotePort:    bastionzero.PtrTo(0),
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "remote_port", "0"),
				),
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

func TestAccDbTarget_LocalPort(t *testing.T) {
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

	localPort1 := 3000
	localPort2 := 4000

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDbTargetConfigLocalPort(rName, env.ID, bzeroTarget.ID, "localhost", "5432", strconv.Itoa(localPort1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
						LocalPort:     &localPort1,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "local_port", strconv.Itoa(localPort1)),
				),
			},
			// Verify import works
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update local port
			{
				Config: testAccDbTargetConfigLocalPort(rName, env.ID, bzeroTarget.ID, "localhost", "5432", strconv.Itoa(localPort2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
						LocalPort:     &localPort2,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckResourceAttr(resourceName, "local_port", strconv.Itoa(localPort2)),
				),
			},
			// Verify setting to null clears it
			{
				Config: testAccDbTargetConfigBasic(rName, env.ID, bzeroTarget.ID, "localhost", "5432"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbTargetExists(resourceName, &target),
					testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
						EnvironmentID: &env.ID,
						Name:          &rName,
						ProxyTargetID: &bzeroTarget.ID,
						RemoteHost:    bastionzero.PtrTo("localhost"),
						RemotePort:    bastionzero.PtrTo(5432),
						LocalPort:     nil,
					}),
					testAccCheckResourceDbTargetComputedAttr(resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "local_port"),
				),
			},
		},
	})
}

func TestAccDbTarget_AllSupportedDatabaseAuthConfig(t *testing.T) {
	// Test creating a db_target resource for each supported database
	// authentication config returned by BastionZero.
	//
	// This test should fail if we add a new supported auth config and forget to
	// update the Terraform provider's validation checks accordingly

	ctx := context.Background()
	rName := acctest.RandomName()
	resourceName := "bastionzero_db_target.test"

	acctest.SkipIfNotInAcceptanceTestMode(t)
	acctest.PreCheck(ctx, t)

	env := new(environments.Environment)
	bzeroTarget := new(targets.BzeroTarget)
	acctest.FindNEnvironmentsOrSkip(t, env)
	acctest.FindNBzeroTargetsOrSkip(t, bzeroTarget)

	// Get all currently supported db auth configs
	supportedDbConfigs, _, err := acctest.APIClient.Targets.ListDatabaseAuthenticationConfigs(ctx)
	if err != nil {
		t.Fatalf("failed to list supported database auth configs: %s", err)
	}

	// Get valid remote_host based on database_authentication_config
	validRemoteHost := func(authConfig *dbauthconfig.DatabaseAuthenticationConfig) string {
		if authConfig.CloudServiceProvider != nil && *authConfig.CloudServiceProvider == dbauthconfig.GCP {
			return "gcp://localhost"
		} else if authConfig.CloudServiceProvider != nil && *authConfig.CloudServiceProvider == dbauthconfig.AWS && authConfig.Database != nil && *authConfig.Database == dbauthconfig.MySQL {
			return "rdsmysql://localhost"
		} else if authConfig.CloudServiceProvider != nil && *authConfig.CloudServiceProvider == dbauthconfig.AWS && authConfig.Database != nil && *authConfig.Database == dbauthconfig.Postgres {
			return "rds://localhost"
		} else {
			return "localhost"
		}
	}

	// Create TF config per db_auth_config. First step creates the resource with
	// the first returned supported db_auth_config. Then each subsequent step
	// updates the TF resource to the next supported db_auth_config, changing
	// the remote_host if needed
	steps := []resource.TestStep{}
	for _, supportedConfigIter := range supportedDbConfigs {
		var target targets.DatabaseTarget

		// Why? See: https://go.dev/doc/faq#closures_and_goroutines
		//
		// Go language issue not fixed until Go 1.22: "This is because each
		// iteration of the loop uses the same instance of the variable v, so
		// each closure shares that single variable."
		supportedConfig := supportedConfigIter // create a new 'supportedConfig' to bind in closure below.
		steps = append(steps, func() resource.TestStep {
			remoteHost := validRemoteHost(&supportedConfig)

			// Build TF resource attr checks based on whether the value is set
			// or not per field
			tfResourceAttrChecks := []resource.TestCheckFunc{}
			if supportedConfig.AuthenticationType != nil {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckResourceAttr(resourceName, "database_authentication_config.authentication_type", *supportedConfig.AuthenticationType))
			} else {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.authentication_type"))
			}
			if supportedConfig.CloudServiceProvider != nil {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckResourceAttr(resourceName, "database_authentication_config.cloud_service_provider", *supportedConfig.CloudServiceProvider))
			} else {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.cloud_service_provider"))
			}
			if supportedConfig.Database != nil {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckResourceAttr(resourceName, "database_authentication_config.database", *supportedConfig.Database))
			} else {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.database"))
			}
			if supportedConfig.Label != nil {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckResourceAttr(resourceName, "database_authentication_config.label", *supportedConfig.Label))
			} else {
				tfResourceAttrChecks = append(tfResourceAttrChecks, resource.TestCheckNoResourceAttr(resourceName, "database_authentication_config.label"))
			}
			// json, _ := json.Marshal(supportedConfig)
			// t.Log(fmt.Sprintf("%v: %v", i, string(json)))

			return resource.TestStep{
				Config: testAccDbTargetConfigDbAuthConfig(rName, env.ID, bzeroTarget.ID, remoteHost, "5432", dbtarget.FlattenDatabaseAuthenticationConfig(ctx, &supportedConfig)),
				Check: resource.ComposeTestCheckFunc(
					append([]resource.TestCheckFunc{
						testAccCheckDbTargetExists(resourceName, &target),
						testAccCheckDbTargetAttributes(t, &target, &expectedDbTarget{
							EnvironmentID:      &env.ID,
							Name:               &rName,
							ProxyTargetID:      &bzeroTarget.ID,
							RemoteHost:         &remoteHost,
							RemotePort:         bastionzero.PtrTo(5432),
							DatabaseAuthConfig: &supportedConfig,
						}),
						testAccCheckResourceDbTargetComputedAttr(resourceName)},
						tfResourceAttrChecks...)...,
				),
			}
		}())
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDbTargetDestroy,
		Steps:                    steps,
	})
}

func TestDbTarget_MutualExclProxyTargetEnv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Cannot specify both proxy target id and proxy environment id
				Config: `
				resource "bastionzero_db_target" "test" {
			      name = "foo"
				  proxy_target_id = "b6d841ca-39ae-414f-ab5b-14be29e5573a"
				  proxy_environment_id = "b6d841ca-39ae-414f-ab5b-14be29e5573a"
				  environment_id = "b6d841ca-39ae-414f-ab5b-14be29e5573a"
				  remote_host = "localhost"
				  remote_port = 5432
				}
				`,
				ExpectError: regexp.MustCompile(`cannot be configured together`),
			},
		},
	})
}

func TestDbTarget_AtLeastOneProxyTargetEnv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// At least one of proxy_target_id or proxy_environment_id must
				// be specified
				Config: `
				resource "bastionzero_db_target" "test" {
			      name = "foo"
				  environment_id = "b6d841ca-39ae-414f-ab5b-14be29e5573a"
				  remote_host = "localhost"
				  remote_port = 5432
				}
				`,
				ExpectError: regexp.MustCompile(`At least one of these attributes must be configured`),
			},
		},
	})
}

func TestDbTarget_InvalidName(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty name not permitted
				Config:      testAccDbTargetConfigBasic("", uuid.New().String(), uuid.New().String(), "localhost", "5432"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestDbTarget_InvalidEnvironmentID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid ID not permitted
				Config:      testAccDbTargetConfigBasic("foo", "bad-id", uuid.New().String(), "localhost", "5432"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidProxyTargetID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid ID not permitted
				Config:      testAccDbTargetConfigBasic("foo", uuid.New().String(), "", "localhost", "5432"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidProxyEnvironmentID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Invalid ID not permitted
				Config:      testAccDbTargetConfigProxyEnvID("foo", uuid.New().String(), "foobar-bad-env-id", "localhost", "5432"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidRemoteHost(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty remote host not permitted
				Config:      testAccDbTargetConfigBasic("foo", uuid.New().String(), uuid.New().String(), "", "5432"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestDbTarget_InvalidAuthConfigAuthType(t *testing.T) {
	dbAuthConfig := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType: bastionzero.PtrTo("foobar"),
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Bad auth type not permitted
				Config:      testAccDbTargetConfigDbAuthConfig("foo", uuid.New().String(), uuid.New().String(), "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(context.Background(), dbAuthConfig)),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidAuthConfigCloudServiceProvider(t *testing.T) {
	dbAuthConfig := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType:   bastionzero.PtrTo(dbauthconfig.Default),
		CloudServiceProvider: bastionzero.PtrTo("foobar"),
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Bad cloud service provider not permitted
				Config:      testAccDbTargetConfigDbAuthConfig("foo", uuid.New().String(), uuid.New().String(), "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(context.Background(), dbAuthConfig)),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidAuthConfigDatabase(t *testing.T) {
	dbAuthConfig := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType: bastionzero.PtrTo(dbauthconfig.Default),
		Database:           bastionzero.PtrTo("foobar"),
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Bad database not permitted
				Config:      testAccDbTargetConfigDbAuthConfig("foo", uuid.New().String(), uuid.New().String(), "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(context.Background(), dbAuthConfig)),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDbTarget_InvalidAuthConfigLabel(t *testing.T) {
	dbAuthConfig := &dbauthconfig.DatabaseAuthenticationConfig{
		AuthenticationType: bastionzero.PtrTo(dbauthconfig.Default),
		Label:              bastionzero.PtrTo(""),
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Empty label not permitted
				Config:      testAccDbTargetConfigDbAuthConfig("foo", uuid.New().String(), uuid.New().String(), "localhost", "5432", dbtarget.FlattenDatabaseAuthenticationConfig(context.Background(), dbAuthConfig)),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
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

func testAccDbTargetConfigLocalPort(name string, envID string, proxyTargetID string, remoteHost string, remotePort string, localPort string) string {
	return fmt.Sprintf(`
resource "bastionzero_db_target" "test" {
  environment_id = %[2]q
  name = %[1]q
  proxy_target_id = %[3]q
  remote_host = %[4]q
  remote_port = %[5]q
  local_port = %[6]q
}
`, name, envID, proxyTargetID, remoteHost, remotePort, localPort)
}

func testAccDbTargetConfigProxyEnvID(name string, envID string, proxyEnvID string, remoteHost string, remotePort string) string {
	return fmt.Sprintf(`
resource "bastionzero_db_target" "test" {
  environment_id = %[2]q
  name = %[1]q
  proxy_environment_id = %[3]q
  remote_host = %[4]q
  remote_port = %[5]q
}
`, name, envID, proxyEnvID, remoteHost, remotePort)
}

type expectedDbTarget struct {
	EnvironmentID      *string
	Name               *string
	ProxyTargetID      *string
	ProxyEnvironmentID *string
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
		if expected.ProxyTargetID == nil && target.ProxyTargetID != "" {
			return fmt.Errorf("Bad proxy_target_id, expected \"%s\", got: %#v", "", target.ProxyTargetID)
		}
		if expected.ProxyEnvironmentID != nil && *expected.ProxyEnvironmentID != target.ProxyEnvironmentID {
			return fmt.Errorf("Bad proxy_environment_id, expected \"%s\", got: %#v", *expected.ProxyEnvironmentID, target.ProxyEnvironmentID)
		}
		if expected.ProxyEnvironmentID == nil && target.ProxyEnvironmentID != "" {
			return fmt.Errorf("Bad proxy_environment_id, expected \"%s\", got: %#v", "", target.ProxyEnvironmentID)
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

func testAccCheckResourceDbTargetComputedAttrProxyEnvID(resourceName string) resource.TestCheckFunc {
	// When using proxy_environment_id, these computed attributes have slightly
	// different values
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "agent_public_key", "n/a"),
		resource.TestCheckResourceAttr(resourceName, "agent_version", "n/a"),
		resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
		resource.TestCheckNoResourceAttr(resourceName, "last_agent_update"),
		resource.TestCheckResourceAttr(resourceName, "region", "n/a"),
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
