package cmd

import (
	"fmt"
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
		if len(args) < 1 {
			if _, err := tea.NewProgram(model.InitialModel(), tea.WithAltScreen()).Run(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
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
