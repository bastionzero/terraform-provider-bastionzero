// bzdatasource provides abstractions that make it easier to write data sources
// for the BastionZero API.
package bzdatasource

// APIModel is a BastionZero API object struct.
type APIModel = interface{}

// TODO-Yuval: Consider using less reflection in SingleDataSource and ListDataSource
// if https://github.com/hashicorp/terraform-plugin-framework/issues/242 is
// resolved at some point.
//
// The main use of reflection is so we don't have to make duplicate fields and
// structs
