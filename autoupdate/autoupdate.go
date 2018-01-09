package autoupdate

import (
	"net/http"
	"io/ioutil"
	"os"
	"errors"
	"strings"
)

type AutoUpdater struct {
	CheckUrl		string
}

func (a *AutoUpdater) Update(url,file string) (err error) {
	req,_ := http.NewRequest("GET",url,nil)
	client := &http.Client{}
	resp,err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	if resp.Status != "200 OK" {
		err = errors.New("unable to get new version")
		return
	}

	data,err := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile(file,data,os.FileMode(0777))
	return
}

func (a *AutoUpdater) Check(v string) (update bool,version string, err error) {
	update = false
	req,_ := http.NewRequest("GET",a.CheckUrl,nil)
	client := &http.Client{}
	resp,err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		err = errors.New("Unable to get required version")
		return
	}

	vers, err := ioutil.ReadAll(resp.Body)
	version = strings.TrimSpace(string(vers))

	if v != version {
		update = true
	}

	return
}

func CreateNew(url string) AutoUpdater {
	a := AutoUpdater{
		CheckUrl: url,
	}

	return a
}