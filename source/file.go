package source

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/andriiyaremenko/mg/dto"
)

func GetFromFile(filePath string) Source {
	urls, err := readFile(filePath)
	if err != nil {
		return func() (*dto.Target, bool, error) { return nil, false, err }
	}

	if len(urls) == 0 {
		return func() (*dto.Target, bool, error) { return nil, false, nil }
	}

	i := 0

	return func() (*dto.Target, bool, error) {
		defer func() { i++ }()

		target := &dto.Target{
			ID:      "",
			URL:     urls[i],
			Method:  "",
			Data:    []byte{},
			Headers: [][2]string{},
			Proxy:   "",
		}

		return target, i < len(urls)-1, nil
	}
}

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	var line string

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())

		if line != "" {
			u, err := url.Parse(line)

			if err != nil || len(u.Host) == 0 {
				return nil, fmt.Errorf("Undefined host or error = %v", err)
			}

			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
