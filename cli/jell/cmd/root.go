package cmd

import (
	"os"

	"github.com/Ygg-Drasill/Jelling/cli/jell/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jell",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tea.NewProgram(model.InitialModel(), tea.WithAltScreen()).Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
