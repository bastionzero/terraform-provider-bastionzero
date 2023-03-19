package serviceaccount

import (
	"context"
	"fmt"
	"time"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/serviceaccounts"
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/types/subjecttype"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// serviceAccountModel maps service account schema data.
type serviceAccountModel struct {
	ID             types.String `tfsdk:"id"`
	Type           types.String `tfsdk:"type"`
	OrganizationID types.String `tfsdk:"organization_id"`
	Email          types.String `tfsdk:"email"`
	ExternalID     types.String `tfsdk:"external_id"`
	JwksURL        types.String `tfsdk:"jwks_url"`
	JwksURLPattern types.String `tfsdk:"jwks_url_pattern"`
	IsAdmin        types.Bool   `tfsdk:"is_admin"`
	TimeCreated    types.String `tfsdk:"time_created"`
	LastLogin      types.String `tfsdk:"last_login"`
	CreatedBy      types.String `tfsdk:"created_by"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

func (m serviceAccountModel) GetID() types.String { return m.ID }

// setServiceAccountAttributes populates the TF schema data from a service
// account API object.
func setServiceAccountAttributes(ctx context.Context, schema *serviceAccountModel, serviceAccount *serviceaccounts.ServiceAccount) {
	schema.ID = types.StringValue(serviceAccount.ID)
	schema.Type = types.StringValue(string(subjecttype.ServiceAccount))
	schema.OrganizationID = types.StringValue(serviceAccount.OrganizationID)
	schema.Email = types.StringValue(serviceAccount.Email)
	schema.ExternalID = types.StringValue(serviceAccount.ExternalID)
	schema.JwksURL = types.StringValue(serviceAccount.JwksURL)
	schema.JwksURLPattern = types.StringValue(serviceAccount.JwksURLPattern)
	schema.IsAdmin = types.BoolValue(serviceAccount.IsAdmin)
	schema.TimeCreated = types.StringValue(serviceAccount.TimeCreated.UTC().Format(time.RFC3339))
	schema.CreatedBy = types.StringValue(serviceAccount.CreatedBy)
	schema.Enabled = types.BoolValue(serviceAccount.Enabled)

	if serviceAccount.LastLogin != nil {
		schema.LastLogin = types.StringValue(serviceAccount.LastLogin.UTC().Format(time.RFC3339))
	} else {
		schema.LastLogin = types.StringNull()
	}
}

func makeServiceAccountDataSourceSchema(withRequiredID bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    !withRequiredID,
			Required:    withRequiredID,
			Description: "The service account's unique ID.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: fmt.Sprintf("The subject's type (constant value \"%s\").", subjecttype.ServiceAccount),
		},
		"organization_id": schema.StringAttribute{
			Computed:    true,
			Description: "The service account's organization's ID.",
		},
		"email": schema.StringAttribute{
			Computed:    true,
			Description: "The service account's email address.",
		},
		"external_id": schema.StringAttribute{
			Computed:    true,
			Description: "The service account's unique per service provider identifier provided by the user during creation.",
		},
		"jwks_url": schema.StringAttribute{
			Computed:    true,
			Description: "The service account's publicly available JWKS URL that provides the public key that can be used to verify the tokens signed by the private key of this service account.",
		},
		"jwks_url_pattern": schema.StringAttribute{
			Computed:    true,
			Description: " A URL pattern that all service accounts of the same service account provider follow in their JWKS URL.",
		},
		"is_admin": schema.BoolAttribute{
			Computed:    true,
			Description: "If true, the service account is an administrator. False otherwise.",
		},
		"time_created": schema.StringAttribute{
			Computed:    true,
			Description: "The time this service account was created in BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
		},
		"last_login": schema.StringAttribute{
			Computed:    true,
			Description: "The time this service account last logged into BastionZero formatted as a UTC timestamp string in RFC 3339 format.",
			Optional:    true,
		},
		"created_by": schema.StringAttribute{
			Computed:    true,
			Description: "Unique identifier for the subject that created this service account.",
		},
		"enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "If true, the service account is currently enabled. False otherwise.",
		},
	}
}
