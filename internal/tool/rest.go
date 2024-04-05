package tool

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

var client = &http.Client{}

func callApi(method string, api string, param []byte) ([]byte, error) {

	addr := viper.GetString("rest.address")
	user := viper.GetString("rest.username")
	pass := viper.GetString("rest.password")
	timeout := viper.GetInt("rest.timeout")

	api, err := url.JoinPath(addr, api)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, api, bytes.NewReader(param))
	req.SetBasicAuth(user, pass)

	client.Timeout = time.Duration(timeout) * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rest: %d %s", resp.StatusCode, b)
	}
	return b, nil
}
