package cli

import (
	"context"
	"os"

	"github.com/beamlit/toolkit/sdk"
	"github.com/spf13/cobra"
)

var BASE_URL = "https://api.beamlit.dev/v0"
var RUN_URL = "https://run.beamlit.dev"
var workspace string
var outputFormat string
var environment string
var client *sdk.ClientWithResponses
var reg *Operations
var verbose bool
var version string
var commit string
var date string

var rootCmd = &cobra.Command{
	Use:   "bl",
	Short: "Beamlit CLI is a command line tool to interact with Beamlit APIs.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if url := os.Getenv("BEAMLIT_API_URL"); url != "" {
			BASE_URL = url
		}
		if runUrl := os.Getenv("BEAMLIT_RUN_URL"); runUrl != "" {
			RUN_URL = runUrl
		}

		reg = &Operations{
			BaseURL: BASE_URL,
			RunURL:  RUN_URL,
		}
		credentials := sdk.LoadCredentials(workspace)
		var err error
		c, err := sdk.NewClientWithCredentials(
			sdk.RunClientWithCredentials{
				ApiURL:      BASE_URL,
				RunURL:      RUN_URL,
				Credentials: credentials,
				Workspace:   workspace,
			},
		)
		if err != nil {
			return err
		}
		client = c

		ctx := context.Background()
		c.RegisterCliCommands(reg, ctx)
		return nil
	},
}

func Execute(releaseVersion string, releaseCommit string, releaseDate string) error {
	reg = &Operations{
		BaseURL: BASE_URL,
	}

	rootCmd.AddCommand(reg.ListOrSetWorkspacesCmd())
	rootCmd.AddCommand(reg.LoginCmd())
	rootCmd.AddCommand(reg.LogoutCmd())
	rootCmd.AddCommand(reg.GetCmd())
	rootCmd.AddCommand(reg.ApplyCmd())
	rootCmd.AddCommand(reg.DeleteCmd())
	rootCmd.AddCommand(reg.RunCmd())
	rootCmd.AddCommand(reg.DocCmd())
	rootCmd.AddCommand(reg.MetricsModelDeploymentCmd())
	rootCmd.AddCommand(reg.ServeCmd())
	rootCmd.AddCommand(reg.CreateAgentAppCmd())
	rootCmd.AddCommand(reg.VersionCmd())

	rootCmd.PersistentFlags().StringVarP(&workspace, "workspace", "w", "", "Specify the workspace name")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: pretty,yaml,json,table")
	rootCmd.PersistentFlags().StringVarP(&environment, "env", "e", "", "Environment. One of: development,production")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	if workspace == "" {
		workspace = sdk.CurrentContext().Workspace
	}
	if environment == "" {
		environment = sdk.CurrentContext().Environment
	}
	if version == "" {
		version = releaseVersion
	}
	if commit == "" {
		commit = releaseCommit
	}
	if date == "" {
		date = releaseDate
	}
	return rootCmd.Execute()
}
