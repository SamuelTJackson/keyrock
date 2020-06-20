package keyrock

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c Client) getURL(suffix string) string {
	return fmt.Sprintf("%s%s",c.options.BaseURL, suffix)
}

func checkError(response *http.Response) error {
	if response.StatusCode >=300 {
		var keyrockError = struct {
			Error KeyrockError `json:"error"`
		}{}
		if err := json.NewDecoder(response.Body).Decode(&keyrockError); err != nil {
			return err
		}
		return keyrockError.Error
	}
	return nil
}
