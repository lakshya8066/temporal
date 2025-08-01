package cassandra

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/log"
	dbschemas "go.temporal.io/server/schema"
	"go.temporal.io/server/temporal/environment"
	"go.temporal.io/server/tools/common/schema"
)

// RunTool runs the temporal-cassandra-tool command line tool
func RunTool(args []string) error {
	app := buildCLIOptions()
	return app.Run(args)
}

var osExit = os.Exit

// root handler for all cli commands
func cliHandler(c *cli.Context, handler func(c *cli.Context, logger log.Logger) error, logger log.Logger) {
	quiet := c.GlobalBool(schema.CLIOptQuiet)
	err := handler(c, logger)
	if err != nil && !quiet {
		osExit(1)
	}
}

func buildCLIOptions() *cli.App {

	app := cli.NewApp()
	app.Name = "temporal-cassandra-tool"
	app.Usage = "Command line tool for temporal cassandra operations"
	app.Version = headers.ServerVersion
	logger := log.NewCLILogger()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   schema.CLIFlagEndpoint,
			Value:  environment.GetCassandraAddress(),
			Usage:  "hostname or ip address of cassandra host to connect to",
			EnvVar: "CASSANDRA_HOST",
		},
		cli.IntFlag{
			Name:   schema.CLIFlagPort,
			Value:  environment.GetCassandraPort(),
			Usage:  "Port of cassandra host to connect to",
			EnvVar: "CASSANDRA_PORT",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagUser,
			Value:  "",
			Usage:  "User name used for authentication for connecting to cassandra host",
			EnvVar: "CASSANDRA_USER",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagPassword,
			Value:  "",
			Usage:  "Password used for authentication for connecting to cassandra host",
			EnvVar: "CASSANDRA_PASSWORD",
		},
		cli.StringSliceFlag{
			Name:   schema.CLIFlagAllowedAuthenticators,
			Value:  nil,
			Usage:  "List of authenticators allowed to be used by the gocql client while connecting to the server.",
			EnvVar: "CASSANDRA_ALLOWED_AUTHENTICATORS",
		},
		cli.IntFlag{
			Name:   schema.CLIFlagTimeout,
			Value:  defaultTimeout,
			Usage:  "request Timeout in seconds used for cql client",
			EnvVar: "CASSANDRA_TIMEOUT",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagKeyspace,
			Value:  "temporal",
			Usage:  "name of the cassandra Keyspace",
			EnvVar: "CASSANDRA_KEYSPACE",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagDatacenter,
			Value:  "",
			Usage:  "enable NetworkTopologyStrategy by providing datacenter name",
			EnvVar: "CASSANDRA_DATACENTER",
		},
		cli.StringFlag{
			Name:   schema.CLIOptAddressTranslator,
			Value:  "",
			Usage:  "name of address translator for cassandra hosts",
			EnvVar: "CASSANDRA_ADDRESS_TRANSLATOR",
		},
		cli.StringFlag{
			Name:   schema.CLIOptAddressTranslatorOptions,
			Value:  "",
			Usage:  "colon-separated list of key=value pairs as options for address translator",
			EnvVar: "CASSANDRA_ADDRESS_TRANSLATOR_OPTIONS_CLI",
		},
		cli.BoolFlag{
			Name:  schema.CLIFlagQuiet,
			Usage: "Don't set exit status to 1 on error",
		},
		cli.BoolFlag{
			Name:   schema.CLIFlagDisableInitialHostLookup,
			Usage:  "instructs gocql driver to only connect to the supplied hosts vs. attempting to lookup additional hosts via the system.peers table",
			EnvVar: "CASSANDRA_DISABLE_INITIAL_HOST_LOOKUP",
		},
		cli.BoolFlag{
			Name:   schema.CLIFlagEnableTLS,
			Usage:  "enable TLS",
			EnvVar: "CASSANDRA_ENABLE_TLS",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagTLSCertFile,
			Usage:  "TLS cert file",
			EnvVar: "CASSANDRA_TLS_CERT",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagTLSKeyFile,
			Usage:  "TLS key file",
			EnvVar: "CASSANDRA_TLS_KEY",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagTLSCaFile,
			Usage:  "TLS CA file",
			EnvVar: "CASSANDRA_TLS_CA",
		},
		cli.StringFlag{
			Name:   schema.CLIFlagTLSHostName,
			Value:  "",
			Usage:  "override for target server name",
			EnvVar: "CASSANDRA_TLS_SERVER_NAME",
		},
		cli.BoolFlag{
			Name:   schema.CLIFlagTLSDisableHostVerification,
			Usage:  "disable tls host name verification (tls must be enabled)",
			EnvVar: "CASSANDRA_TLS_DISABLE_HOST_VERIFICATION",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "setup-schema",
			Aliases: []string{"setup"},
			Usage:   "setup initial version of cassandra schema",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  schema.CLIFlagVersion,
					Usage: "initial version of the schema, cannot be used with disable-versioning",
				},
				cli.StringFlag{
					Name:  schema.CLIFlagSchemaFile,
					Usage: "path to the .cql schema file; if un-specified, will just setup versioning tables",
				},
				cli.StringFlag{
					Name: schema.CLIFlagSchemaName,
					Usage: fmt.Sprintf("name of embedded schema directory with .cql file, one of: %v",
						dbschemas.PathsByDB("cassandra")),
				},
				cli.BoolFlag{
					Name:  schema.CLIFlagDisableVersioning,
					Usage: "disable setup of schema versioning",
				},
				cli.BoolFlag{
					Name:  schema.CLIFlagOverwrite,
					Usage: "drop all existing tables before setting up new schema",
				},
			},
			Action: func(c *cli.Context) {
				cliHandler(c, setupSchema, logger)
			},
		},
		{
			Name:    "update-schema",
			Aliases: []string{"update"},
			Usage:   "update cassandra schema to a specific version",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  schema.CLIFlagTargetVersion,
					Usage: "target version for the schema update, defaults to latest",
				},
				cli.StringFlag{
					Name:  schema.CLIFlagSchemaDir,
					Usage: "path to directory containing versioned schema",
				},
				cli.StringFlag{
					Name: schema.CLIFlagSchemaName,
					Usage: fmt.Sprintf("name of embedded versioned schema, one of: %v",
						dbschemas.PathsByDB("cassandra")),
				},
			},
			Action: func(c *cli.Context) {
				cliHandler(c, updateSchema, logger)
			},
		},
		{
			Name:    "create-keyspace",
			Aliases: []string{"create", "create-Keyspace"},
			Usage:   "creates a keyspace with simple strategy or network topology if datacenter name is provided",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  schema.CLIFlagKeyspace,
					Usage: "name of the keyspace",
				},
				cli.IntFlag{
					Name:  schema.CLIFlagReplicationFactor,
					Value: 1,
					Usage: "replication factor for the keyspace",
				},
				cli.StringFlag{
					Name:  schema.CLIFlagDatacenter,
					Value: "",
					Usage: "enable NetworkTopologyStrategy by providing datacenter name",
				},
			},
			Action: func(c *cli.Context) {
				cliHandler(c, createKeyspace, logger)
			},
		},
		{
			Name:    "drop-keyspace",
			Aliases: []string{"drop"},
			Usage:   "drops a keyspace with simple strategy or network topology if datacenter name is provided",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  schema.CLIFlagKeyspace,
					Usage: "name of the keyspace",
				},
				cli.IntFlag{
					Name:  schema.CLIFlagReplicationFactor,
					Value: 1,
					Usage: "replication factor for the keyspace",
				},
				cli.StringFlag{
					Name:  schema.CLIFlagDatacenter,
					Value: "",
					Usage: "enable NetworkTopologyStrategy by providing datacenter name",
				},
				cli.BoolFlag{
					Name:  schema.CLIFlagForce,
					Usage: "don't prompt for confirmation",
				},
			},
			Action: func(c *cli.Context) {
				drop := c.Bool(schema.CLIOptForce)
				if !drop {
					keyspace := c.String(schema.CLIOptKeyspace)
					fmt.Printf("Are you sure you want to drop keyspace %q (y/N)? ", keyspace)
					y := ""
					_, _ = fmt.Scanln(&y)
					if y == "y" || y == "Y" {
						drop = true
					}
				}
				if drop {
					cliHandler(c, dropKeyspace, logger)
				}
			},
		},
		{
			Name:    "validate-health",
			Aliases: []string{"vh"},
			Usage:   "validates health of cassandra by attempting to establish CQL session to system keyspace",
			Action: func(c *cli.Context) {
				cliHandler(c, validateHealth, logger)
			},
		},
	}

	return app
}
