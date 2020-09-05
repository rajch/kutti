package cmd

import (
	"fmt"
	"os"

	"github.com/rajch/kutti/internal/pkg/kuttilog"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "kutti",
	Short:            "Manage multi-node kubernetes clusters in a local environment",
	Long:             `Manage multi-node kubernetes clusters in a local environment.`,
	PersistentPreRun: setverbosity,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil && err != cobra.ErrSubCommandRequired {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "produce minimum output")
	rootCmd.PersistentFlags().Bool("debug", false, "produce maximum output")
}

func setverbosity(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		kuttilog.Setloglevel(4)
	} else {
		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			kuttilog.Setloglevel(0)
		}
	}

	kuttilog.SetPrefix("")
}
