package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Authorization запрос к Zabbix с целью получения хеша авторизации необходимого для
// дальнейшей работы с API
func (zc *ZabbixConnectionJsonRPC) Authorization(ctx context.Context) error {
	data := strings.NewReader(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"user.login",
	  "params": {
	    "username":"%s",
		"password":"%s"
	  },
	  "id":1
	}`, zc.login, zc.passwd))

	result := ZabbixAuthorizationData{}
	res, err := zc.PostRequest(ctx, data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return err
	}

	if len(result.Error) > 0 {
		var shortMsg, fullMsg string
		for k, v := range result.Error {
			if k == "message" {
				shortMsg = fmt.Sprint(v)
			}
			if k == "data" {
				fullMsg = fmt.Sprint(v)
			}
		}

		return fmt.Errorf("error authorization, (%s %s)", shortMsg, fullMsg)
	}

	zc.authorizationHash = result.Result

	return nil
}

// GetAuthorizationData хеш авторизации
func (zc *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return zc.authorizationHash
}

// PostRequest HTTP запрос типа POST
func (zc *ZabbixConnectionJsonRPC) PostRequest(ctx context.Context, data *strings.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", zc.url, data)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", zc.authorizationHash))
	req.Header.Set("Content-Type", "application/json-rpc")

	res, err := zc.connClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error sending the request, response status is %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
