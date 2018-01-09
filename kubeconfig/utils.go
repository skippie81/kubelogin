package kubeconfig

import (
	"regexp"
	"errors"
)

const (
	authServerName = "https://auth"
)

func GetLoginServer(apiserver string) (server string, err error) {

	re := regexp.MustCompile(`http[s]{1}:\/\/[a-zA-z0-9_\-]+`)

	if ! re.Match([]byte(apiserver)) {
		err = errors.New("unable to auto-detect authentication server")
	}

	server = re.ReplaceAllString(apiserver,authServerName)

	return
}
