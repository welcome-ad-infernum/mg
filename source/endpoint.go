package source

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/andriiyaremenko/mg/client"
	"github.com/andriiyaremenko/mg/dto"
)

func GetFromEndpoint(managerURL string) Source {
	return func() (*dto.Target, bool, error) {
		client := client.New(time.Second * 10)

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
