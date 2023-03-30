package environment_test

// func TestAccEnvironment_CreateAndUpdate(t *testing.T) {
// 	var env environments.Environment
// 	name := acctest.RandomName()
// 	resourceName := "bastionzero_environment.env"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { acctest.PreCheck(context.Background(), t) },
// 		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
// 		CheckDestroy:             testAccCheckEnvironmentDestroy,
// 		Steps: []resource.TestStep{
// 			// Create environment
// 			{
// 				Config: environmentResourceTFConfig("env", &environmentResourceOptions{
// 					Name:        &name,
// 					Description: bastionzero.PtrTo("test")}),
// 				Check: resource.ComposeTestCheckFunc(
// 					acctest.CheckEnvironmentExists(resourceName, &env),
// 					testAccCheckEnvironmentAttributes(&env, &environmentResourceOptions{
// 						Name:                       &name,
// 						Description:                bastionzero.PtrTo("test"),
// 						OfflineCleanupTimeoutHours: bastionzero.PtrTo(environment.DefaultOfflineCleanupTimeoutHours)}),
// 					resource.TestCheckResourceAttr(resourceName, "name", name),
// 					resource.TestCheckResourceAttr(resourceName, "description", "test"),
// 					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
// 					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
// 					resource.TestCheckResourceAttr(resourceName, "targets.%", "0"),
// 					resource.TestMatchResourceAttr(resourceName, "id", acctest.ExpectedIDRegEx()),
// 					resource.TestMatchResourceAttr(resourceName, "organization_id", acctest.ExpectedIDRegEx()),
// 					resource.TestMatchResourceAttr(resourceName, "time_created", acctest.ExpectedTimestampRegEx()),
// 				),
// 			},
// 			// Modify description and cleanup timeout
// 			{
// 				Config: environmentResourceTFConfig("env", &environmentResourceOptions{
// 					Name:                       &name,
// 					Description:                bastionzero.PtrTo("new-desc"),
// 					OfflineCleanupTimeoutHours: bastionzero.PtrTo(3000)}),
// 				Check: resource.ComposeTestCheckFunc(
// 					acctest.CheckEnvironmentExists(resourceName, &env),
// 					testAccCheckEnvironmentAttributes(&env, &environmentResourceOptions{
// 						Name:                       &name,
// 						Description:                bastionzero.PtrTo("new-desc"),
// 						OfflineCleanupTimeoutHours: bastionzero.PtrTo(3000)}),
// 					resource.TestCheckResourceAttr(resourceName, "name", name),
// 					resource.TestCheckResourceAttr(resourceName, "description", "new-desc"),
// 					resource.TestCheckResourceAttr(resourceName, "offline_cleanup_timeout_hours", strconv.Itoa(3000)),
// 					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
// 					resource.TestCheckResourceAttr(resourceName, "targets.%", "0"),
// 					resource.TestMatchResourceAttr(resourceName, "id", acctest.ExpectedIDRegEx()),
// 					resource.TestMatchResourceAttr(resourceName, "organization_id", acctest.ExpectedIDRegEx()),
// 					resource.TestMatchResourceAttr(resourceName, "time_created", acctest.ExpectedTimestampRegEx()),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccEnvironment_UpdateName(t *testing.T) {
// 	var afterCreate, afterUpdate environments.Environment
// 	name := acctest.RandomName()
// 	name2 := acctest.RandomName()
// 	resourceName := "bastionzero_environment.env"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { acctest.PreCheck(context.Background(), t) },
// 		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
// 		CheckDestroy:             testAccCheckEnvironmentDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: environmentResourceTFConfig("env", &environmentResourceOptions{Name: &name}),
// 				Check: resource.ComposeTestCheckFunc(
// 					acctest.CheckEnvironmentExists(resourceName, &afterCreate),
// 					testAccCheckEnvironmentAttributes(&afterCreate, &environmentResourceOptions{Name: &name}),
// 					resource.TestCheckResourceAttr(resourceName, "name", name),
// 				),
// 			},
// 			// Change name should force re-create the resource
// 			{
// 				Config: environmentResourceTFConfig("env", &environmentResourceOptions{Name: &name2}),
// 				Check: resource.ComposeTestCheckFunc(
// 					acctest.CheckEnvironmentExists(resourceName, &afterUpdate),
// 					testAccCheckEnvironmentAttributes(&afterUpdate, &environmentResourceOptions{Name: &name2}),
// 					resource.TestCheckResourceAttr(resourceName, "name", name2),
// 					testAccCheckEnvironmentRecreated(t, &afterCreate, &afterUpdate),
// 				),
// 			},
// 		},
// 	})
// }

// type environmentResourceOptions struct {
// 	Name                       *string
// 	Description                *string
// 	OfflineCleanupTimeoutHours *int
// }

// func environmentResourceTFConfig(resourceName string, opts *environmentResourceOptions) string {
// 	var name, description, cleanupTimeout string
// 	if opts.Name != nil {
// 		name = acctest.SurroundDoubleQuote(*opts.Name)
// 	} else {
// 		name = "null"
// 	}
// 	if opts.Description != nil {
// 		description = acctest.SurroundDoubleQuote(*opts.Description)
// 	} else {
// 		description = "null"
// 	}
// 	if opts.OfflineCleanupTimeoutHours != nil {
// 		cleanupTimeout = strconv.Itoa(*opts.OfflineCleanupTimeoutHours)
// 	} else {
// 		cleanupTimeout = "null"
// 	}

// 	return fmt.Sprintf(`
// resource "bastionzero_environment" "%s" {
//   name   = %s
//   description = %s
//   offline_cleanup_timeout_hours = %s
// }
// `, resourceName, name, description, cleanupTimeout)
// }

// func testAccCheckEnvironmentDestroy(s *terraform.State) error {
// 	client := acctest.GetBastionZeroClient()

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "bastionzero_environment" {
// 			continue
// 		}

// 		// Try to find the environment
// 		_, _, err := client.Environments.GetEnvironment(context.Background(), rs.Primary.ID)
// 		if err != nil && !apierror.IsAPIErrorStatusCode(err, http.StatusNotFound) {
// 			return fmt.Errorf("Error waiting for environment (%s) to be destroyed: %s", rs.Primary.ID, err)
// 		}
// 	}

// 	return nil
// }

// func testAccCheckEnvironmentRecreated(t *testing.T, before, after *environments.Environment) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if before.ID == after.ID {
// 			t.Fatalf("Expected change of environment IDs, but both were %v", before.ID)
// 		}
// 		return nil
// 	}
// }

// func testAccCheckEnvironmentAttributes(env *environments.Environment, opts *environmentResourceOptions) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {

// 		if opts.Name != nil && *opts.Name != env.Name {
// 			return fmt.Errorf("Bad name, expected \"%s\", got: %#v", *opts.Name, env.Name)
// 		}
// 		if opts.Description != nil && *opts.Description != env.Description {
// 			return fmt.Errorf("Bad description, expected \"%s\", got: %#v", *opts.Description, env.Description)
// 		}
// 		if opts.OfflineCleanupTimeoutHours != nil && *opts.OfflineCleanupTimeoutHours != int(env.OfflineCleanupTimeoutHours) {
// 			return fmt.Errorf("Bad offline_cleanup_timeout_hours, expected \"%v\", got: %#v", *opts.OfflineCleanupTimeoutHours, env.OfflineCleanupTimeoutHours)
// 		}

// 		return nil
// 	}
// }
