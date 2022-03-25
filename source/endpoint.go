package source

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/andriiyaremenko/mg/dto"
)

func GetFromEndpoint(client *http.Client, managerURL string) Source {
	return func() (*dto.Target, bool, error) {
		r, err := client.Get(managerURL)
		if err != nil {
			return nil, true, err
		}

		defer r.Body.Close()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, true, err
		}

		target := new(dto.Target)
		if err := json.Unmarshal(b, target); err != nil {
			return nil, true, err
		}

		target.URL = strings.TrimSpace(target.URL)
		return target, true, nil
	}
}
