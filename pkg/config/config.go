package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	// Add the SchemaFields for the Config.

	ClientIDField = field.StringField(
		"paylocity-client-id",
		field.WithDisplayName("Client ID"),
		field.WithDescription("Client ID for the Paylocity API"),
		field.WithRequired(true),
	)

	ClientSecretField = field.StringField(
		"paylocity-client-secret",
		field.WithDisplayName("Client Secret"),
		field.WithDescription("Client Secret for the Paylocity API"),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	BaseURLField = field.StringField(
		"paylocity-base-url",
		field.WithDisplayName("API Gateway URL"),
		field.WithDescription("The base API Gateway URL for the Paylocity environment. e.g., Testing: https://dc1demogwext.paylocity.com or Production: https://dc1prodgwext.paylocity.com"),
		field.WithRequired(true),
		field.WithDefaultValue("https://dc1prodgwext.paylocity.com"),
	)

	CompanyIDField = field.StringField(
		"paylocity-company-id",
		field.WithDisplayName("Company ID"),
		field.WithDescription("The ID of the specific company in Paylocity to fetch data from."),
		field.WithRequired(true),
	)

	ConfigurationFields = []field.SchemaField{
		ClientIDField,
		ClientSecretField,
		CompanyIDField,
		BaseURLField,
	}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	// For example, a username and password can be required together, or an access token can be
	// marked as mutually exclusive from the username password pair.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Configuration = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Paylocity"),
	field.WithHelpUrl("/docs/baton/paylocity"),
	field.WithIconUrl("/static/app-icons/paylocity.svg"),
)
