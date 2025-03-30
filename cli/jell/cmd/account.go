package cmd

import (
	"fmt"
	"github.com/Ygg-Drasill/Jelling/cli/jell/model/account"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"os"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login to Jelling",
		Run:   login,
	}

	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "logout from Jelling",
		Run:   logout,
	}

	registerCmd = &cobra.Command{
		Use:   "register",
		Short: "register an account on Jelling",
		Run:   register,
	}
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(registerCmd)
}

func login(cmd *cobra.Command, args []string) {
	if _, err := tea.NewProgram(account.InitAccountModel(account.ModeLogin)).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func logout(cmd *cobra.Command, args []string) {

}

func register(cmd *cobra.Command, args []string) {
	if _, err := tea.NewProgram(account.InitAccountModel(account.ModeRegister)).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
