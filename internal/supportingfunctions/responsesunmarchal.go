package supportingfunctions

import "encoding/json"

// ResponseUnmarchal декодирование JSON ответа
func ResponseUnmarchal[T, E any](b []byte, res T, resErr E) (T, E, error) {
	err := json.Unmarshal(b, res)
	if err != nil {
		err = json.Unmarshal(b, resErr)

		return res, resErr, err
	}

	return res, resErr, nil
}
