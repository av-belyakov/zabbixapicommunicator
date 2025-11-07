package connectionjsonrpc

import "encoding/json"

// ResponseDecode декодирует ответ полученный от Zabbix API
func ResponseDecode(data []byte) (*ResponseMessage, error) {
	res := &ResponseMessage{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}
