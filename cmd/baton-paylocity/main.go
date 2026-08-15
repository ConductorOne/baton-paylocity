package main

import (
	"context"

	cfg "github.com/conductorone/baton-paylocity/pkg/config"
	"github.com/conductorone/baton-paylocity/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(
		ctx,
		"baton-paylocity",
		version,
		cfg.Configuration,
		connector.New,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilderV2(&connector.Connector{}),
	)
}
