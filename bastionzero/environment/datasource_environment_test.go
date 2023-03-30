package environment_test

// func TestAccDataSourceEnvironment_BasicById(t *testing.T) {
// 	var env environments.Environment
// 	// Create random env name
// 	name := acctest.RandomName()

// 	resourceConfig := fmt.Sprintf(`
// 	resource "bastionzero_environment" "env" {
// 		name = %v
// 	}`, name)
// 	dataSourceConfig := `
// 	data "bastionzero_environment" "env" {
// 		id = bastionzero_environment.env.id
// 	}`

// 	dataSourceRefName := "data.bastionzero_environment.env"
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { acctest.PreCheck(context.Background(), t) },
// 		ProtoV6ProviderFactories: acctest.TestProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// First ensure we can create the resource
// 			{
// 				Config: resourceConfig,
// 			},
// 			// Read testing
// 			{
// 				Config: resourceConfig + dataSourceConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					acctest.CheckEnvironmentExists(dataSourceRefName, &env),
// 					resource.TestCheckResourceAttr(dataSourceRefName, "name", name),
// 					resource.TestCheckResourceAttr(dataSourceRefName, "description", ""),
// 					resource.TestCheckResourceAttr(dataSourceRefName, "offline_cleanup_timeout_hours", strconv.Itoa(environment.DefaultOfflineCleanupTimeoutHours)),
// 					resource.TestCheckResourceAttr(dataSourceRefName, "is_default", "false"),
// 					resource.TestCheckResourceAttr(dataSourceRefName, "targets.%", "0"),
// 					resource.TestMatchResourceAttr(dataSourceRefName, "id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
// 					resource.TestMatchResourceAttr(dataSourceRefName, "organization_id", regexp.MustCompile(acctest.UUIDV4RegexPattern)),
// 					resource.TestMatchResourceAttr(dataSourceRefName, "time_created", regexp.MustCompile(acctest.RFC3339RegexPattern)),
// 				),
// 			},
// 		},
// 	})
// }
