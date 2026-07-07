package github

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchPublicSSHKeys(login string) ([]string, error) {
	endpoint := fmt.Sprintf(PublicKeysURL, login)

	response, responseError := http.Get(endpoint)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(PublicKeysFetchFailed, response.StatusCode)
	}

	body, readError := io.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	keys := make([]string, 0)
	for _, line := range strings.Split(string(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			keys = append(keys, trimmed)
		}
	}

	return keys, nil
}
