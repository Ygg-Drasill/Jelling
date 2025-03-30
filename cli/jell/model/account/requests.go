package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/Jelling/common/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

type FetchCompleteMsg struct {
	statusCode int
	err        error
}

var (
	client          = http.DefaultClient
	registerUrl     string
	authenticateUrl string
)

func init() {
	var err error
	baseUrl := viper.GetString("baseUrl")
	registerUrl = baseUrl + "/account/register"
	authenticateUrl = baseUrl + "/account/auth"

	client.Timeout = time.Second * 5
	client.Jar, err = cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func register(username, password string) tea.Cmd {
	return func() tea.Msg {
		var err error
		account := api.AccountRequest{
			Username: username,
			Password: password,
		}
		data, err := json.Marshal(account)
		if err != nil {
			tea.Println(err)
		}
		body := bytes.NewBuffer(data)
		request, err := http.NewRequest("POST", "http://0.0.0.0:30420/api/v1/account/register", body)
		if err != nil {
			tea.Println(err)
		}
		request.Header.Set("Content-Type", "application/json")
		defer request.Body.Close()
		response, err := client.Do(request)
		if err != nil {
			tea.Println(err)
		}
		if response.StatusCode >= 400 {
			err = fmt.Errorf("help")
		}

		var res api.SessionTokenResponse
		err = json.NewDecoder(response.Body).Decode(&res)

		tea.Println(res.Token)

		return FetchCompleteMsg{
			statusCode: response.StatusCode,
			err:        err,
		}
	}
}

func authenticate(username, password string) tea.Cmd {
	return func() tea.Msg {
		return 200
	}
}
