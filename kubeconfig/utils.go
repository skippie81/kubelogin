package kubeconfig

import (
	"regexp"
	"errors"
	"fmt"
)

const (
	authServerName = "https://auth."
)

func GetLoginServer(apiserver string) (server string, err error) {
	re := regexp.MustCompile(`https:\/\/`)
	srv := re.ReplaceAllString(apiserver,"")

	re = regexp.MustCompile(`^k8s\-api\.`)
	if re.Match([]byte(srv)) == true {
		server = re.ReplaceAllString(srv, authServerName)
		return
	}
	err = errors.New(fmt.Sprintf("Unable to determine login server for: %s",apiserver))
	return
}
