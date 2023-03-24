package bztftest

import "regexp"

// ExpectedTimestampRegEx returns the expected regex used to validate timestamps
// stored in the Terraform state
func ExpectedTimestampRegEx() *regexp.Regexp {
	// Source: https://www.regextester.com/115563
	expected, _ := regexp.Compile(`^[1-9]\d{3}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)
	return expected
}

// ExpectedIDRegEx returns the expected regex used to validate BastionZero API
// objects with unique IDs stored in the Terraform state
func ExpectedIDRegEx() *regexp.Regexp {
	// Source: https://gist.github.com/johnelliott/cf77003f72f889abbc3f32785fa3df8d?permalink_comment_id=4318295#gistcomment-4318295
	expected, _ := regexp.Compile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return expected
}
