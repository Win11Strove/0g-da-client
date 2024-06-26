package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/0glabs/0g-da-client/inabox/deploy"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName        = "testname"
	rootPathFlagName        = "root-path"
	localstackFlagName      = "localstack-port"
	deployResourcesFlagName = "deploy-resources"

	metadataTableName = "test-BlobMetadata"
	bucketTableName   = "test-zgda-blobstore"

	chainCmdName      = "chain"
	localstackCmdName = "localstack"
	expCmdName        = "exp"
	allCmdName        = "all"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    testNameFlagName,
				Usage:   "name of the test to run (in `inabox/testdata`)",
				EnvVars: []string{"ZGDA_TESTDATA_PATH"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:  rootPathFlagName,
				Usage: "path to the root of repo",
				Value: "../",
			},
			&cli.StringFlag{
				Name:  localstackFlagName,
				Value: "",
				Usage: "path to the config file",
			},
			&cli.BoolFlag{
				Name:  deployResourcesFlagName,
				Value: true,
				Usage: "whether to deploy localstack resources",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   chainCmdName,
				Usage:  "deploy the chain infrastructure (anvil, graph) for the inabox test",
				Action: getRunner(chainCmdName),
			},
			{
				Name:   localstackCmdName,
				Usage:  "deploy localstack and create the AWS resources needed for the inabox test",
				Action: getRunner(localstackCmdName),
			},
			{
				Name:   expCmdName,
				Usage:  "deploy the contracts and create configurations for all ZGDA components",
				Action: getRunner(expCmdName),
			},
			{
				Name:   allCmdName,
				Usage:  "deploy all infra, resources, contracts",
				Action: getRunner(allCmdName),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getRunner(command string) func(ctx *cli.Context) error {

	return func(ctx *cli.Context) error {

		if command != localstackCmdName {
			rootPath, err := filepath.Abs(ctx.String(rootPathFlagName))
			if err != nil {
				return err
			}
			testname := ctx.String(testNameFlagName)
			if testname == "" {
				testname, err = deploy.GetLatestTestDirectory(rootPath)
				if err != nil {
					return err
				}
			}
		}

		switch command {
		case localstackCmdName:
			return localstack(ctx)
		}

		return nil

	}

}

func localstack(ctx *cli.Context) error {

	pool, _, err := deploy.StartDockertestWithLocalstackContainer(ctx.String(localstackFlagName))
	if err != nil {
		return err
	}

	if ctx.Bool(deployResourcesFlagName) {
		return deploy.DeployResources(pool, ctx.String(localstackFlagName), metadataTableName, bucketTableName)
		//return deploy.DeployResources(nil, ctx.String(localstackFlagName), metadataTableName, bucketTableName)
	}

	return nil
}
