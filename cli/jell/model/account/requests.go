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
	"net/url"
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
	if baseUrl == "" {
		panic("no base url found")
	}
	registerUrl = baseUrl + "/account/register"
	authenticateUrl = baseUrl + "/account/auth"

	client.Timeout = time.Second * 5
	client.Jar, err = cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cookies := make([]*http.Cookie, 1)
	cookies[0] = api.NewSessionCookie(
		viper.GetString("sessionToken"),
		time.UnixMilli(viper.GetInt64("sessionTokenExpiry")))
	client.Jar.SetCookies(baseUrlParsed, cookies)
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
			return FetchCompleteMsg{err: err}
		}
		body := bytes.NewBuffer(data)
		request, err := http.NewRequest("POST", registerUrl, body)
		if err != nil {
			tea.Println(err)
			return FetchCompleteMsg{err: err}
		}
		request.Header.Set("Content-Type", "application/json")
		defer request.Body.Close()
		response, err := client.Do(request)
		if err != nil {
			tea.Println(err)
			return FetchCompleteMsg{err: err}
		}
		if response.StatusCode >= 400 {
			err = fmt.Errorf("help")
		}

		var tokenResponse api.SessionTokenResponse
		err = json.NewDecoder(response.Body).Decode(&tokenResponse)

		err = saveSession(response.Cookies())

		return FetchCompleteMsg{
			statusCode: response.StatusCode,
			err:        err,
		}
	}
}

func authenticate(username, password string) tea.Cmd {
	return func() tea.Msg {
		account := api.AccountRequest{
			Username: username,
			Password: password,
		}

		requestBody, err := json.Marshal(account)
		requestBodyReader := bytes.NewBuffer(requestBody)
		request, err := http.NewRequest("POST", authenticateUrl, requestBodyReader)
		response, err := client.Do(request)
		if response.StatusCode < 400 {
			err = saveSession(response.Cookies())
		}

		return FetchCompleteMsg{
			statusCode: response.StatusCode,
			err:        err,
		}
	}
}

func saveSession(cookies []*http.Cookie) error {
	for _, c := range cookies {
		if c.Name == "session" {
			viper.Set("sessionToken", c.Value)
			viper.Set("sessionTokenExpiry", c.Expires.UnixMilli())
			err := viper.WriteConfig()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
