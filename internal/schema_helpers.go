package internal

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Source: https://github.com/Doridian/terraform-provider-hexonet/blob/486274bbd561dda6bcfc5f1a1c808dd05a6d9053/hexonet/utils/schema_helpers.go#L10
// Related issue: https://github.com/hashicorp/terraform-plugin-framework/issues/568

// ResourceSchemaToDataSourceSchema converts a resource schema to a data source
// schema by copying all common fields between the two schemas, and setting
// required = true for the specified id field
func ResourceSchemaToDataSourceSchema(resourceSchema map[string]resource_schema.Attribute, idField *string) map[string]datasource_schema.Attribute {
	foundIdField := false

	datasourceSchema := make(map[string]datasource_schema.Attribute)
	for name, srcAttr := range resourceSchema {
		optional := false
		required := false
		computed := true

		if idField != nil && name == *idField {
			required = true
			computed = false
			foundIdField = true
		}

		switch srcAttrTyped := srcAttr.(type) {
		case resource_schema.StringAttribute:
			datasourceSchema[name] = datasource_schema.StringAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.BoolAttribute:
			datasourceSchema[name] = datasource_schema.BoolAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.Int64Attribute:
			datasourceSchema[name] = datasource_schema.Int64Attribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.ListAttribute:
			datasourceSchema[name] = datasource_schema.ListAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				ElementType:         srcAttrTyped.ElementType,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.MapAttribute:
			datasourceSchema[name] = datasource_schema.MapAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				ElementType:         srcAttrTyped.ElementType,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.MapNestedAttribute:
			datasourceSchema[name] = datasource_schema.MapNestedAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				NestedObject: datasource_schema.NestedAttributeObject{
					Attributes: ResourceSchemaToDataSourceSchema(srcAttrTyped.NestedObject.Attributes, nil),
					Validators: srcAttrTyped.NestedObject.Validators,
					CustomType: srcAttrTyped.NestedObject.CustomType,
				},
				Optional: optional,
				Required: required,
				Computed: computed,
			}
		case resource_schema.SetAttribute:
			datasourceSchema[name] = datasource_schema.SetAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				ElementType:         srcAttrTyped.ElementType,
				Optional:            optional,
				Required:            required,
				Computed:            computed,
			}
		case resource_schema.SetNestedAttribute:
			datasourceSchema[name] = datasource_schema.SetNestedAttribute{
				Validators:          srcAttrTyped.Validators,
				Description:         srcAttrTyped.Description,
				MarkdownDescription: srcAttrTyped.MarkdownDescription,
				CustomType:          srcAttrTyped.CustomType,
				Sensitive:           srcAttrTyped.Sensitive,
				NestedObject: datasource_schema.NestedAttributeObject{
					Attributes: ResourceSchemaToDataSourceSchema(srcAttrTyped.NestedObject.Attributes, nil),
					Validators: srcAttrTyped.NestedObject.Validators,
					CustomType: srcAttrTyped.NestedObject.CustomType,
				},
				Optional: optional,
				Required: required,
				Computed: computed,
			}
		default:
			log.Panicf("unknown attribute type: %v", srcAttr.GetType().String())
		}
	}

	if idField != nil && !foundIdField {
		log.Panicf("id field \"%s\" not found in resource schema", *idField)
	}

	return datasourceSchema
}

// Source: https://github.com/hashicorp/terraform-provider-aws/blob/0e19050852dadd4498d77467b8c2692b49881b22/internal/framework/attrtypes.go#L13

// AttributeTypes returns a map of attribute types for the specified type T.
// T must be a struct and reflection is used to find exported fields of T with the `tfsdk` tag.
func AttributeTypes[T any](ctx context.Context) (map[string]attr.Type, error) {
	var t T
	val := reflect.ValueOf(t)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%T has unsupported type: %s", t, typ)
	}

	attributeTypes := make(map[string]attr.Type)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields.
		}
		tag := field.Tag.Get(`tfsdk`)
		if tag == "-" {
			continue // Skip explicitly excluded fields.
		}
		if tag == "" {
			return nil, fmt.Errorf(`%T needs a struct tag for "tfsdk" on %s`, t, field.Name)
		}

		if v, ok := val.Field(i).Interface().(attr.Value); ok {
			attributeTypes[tag] = v.Type(ctx)
		}
	}

	return attributeTypes, nil
}
