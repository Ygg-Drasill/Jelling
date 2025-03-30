package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/Jelling/cli/jell/model/account"
	"github.com/Ygg-Drasill/Jelling/common/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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

	client = http.DefaultClient
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

	account := &api.AccountRequest{
		Username: "alexref",
		Password: "1234",
	}

	data, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", "http://0.0.0.0:30420/api/v1/account/register", body)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	defer request.Body.Close()
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)

}
