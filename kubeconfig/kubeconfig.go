package kubeconfig

import (
	yaml "gopkg.in/yaml.v2"
	"fmt"
	"errors"
	"net/http"
	"io/ioutil"
	"os"
)

type KubeConfig struct {
	Clusters	[]*Cluster	`yaml:"clusters"`
	Contexts	[]*Context  	`yaml:"contexts"`
	CurrentContext 	string		`yaml:"current-context"`
	Preferences	struct{}	`yaml:"preferences"`
	Users 		[]*User		`yaml:"users"`
}

type Cluster struct {
	Cluster		struct{
		CertificateAuthorityData	string		`yaml:"certificate-authority-data,omitempty"`
		Server				string		`yaml:"server"`
		LoginServer			string		`yaml:"login-server,omitempty"`
			       }	`yaml:"cluster"`
	Name		string		`yaml:"name"`
}

type Context struct {
	Context		struct {
		Cluster		string		`yaml:"cluster"`
		User 		string		`yaml:"user"`
		Namespace	string		`yaml:"namespace,omitempty"`
			       }	`yaml:"context"`
	Name		string		`yaml:"name"`
}

type User struct {
	User		struct {
		Token		string		`yaml:"token"`
		LoginUser	string		`yaml:"login-user,omitempty"`
			    }		`yaml:"user"`
	Name 		string		`yaml:"name"`
}

func (k KubeConfig) String() string {
	data,err := yaml.Marshal(&k)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s",data)
}

func (k *KubeConfig) GetCurrentContext() (context string, err error) {
	if context = k.CurrentContext; context != "" {
		return
	} else {
		err = errors.New("no current context")
	}
	return
}

func (k *KubeConfig) GetCurrentUsername() (username string, err error) {
	context,err := k.GetCurrentContext()

	if err != nil {
		return
	}

	for _,c := range k.Contexts {
		if c.Name == context {
			username = c.Context.User
			return
		}
	}

	err = errors.New("Current context not found")
	return
}

func (k *KubeConfig) GetCurrentClusterName() (clustername string, err error) {
	context,err := k.GetCurrentContext()

	if err != nil {
		return
	}

	for _,c := range k.Contexts {
		if c.Name == context {
			clustername = c.Context.Cluster
			return
		}
	}

	err = errors.New("Current cluster not found")
	return
}

func (k *KubeConfig) GetCurrentAPIServer() (apiserver string,err error) {
	cluster,err := k.GetCurrentClusterName()

	if err != nil {
		return
	}

	for _,c := range k.Clusters {
		if c.Name == cluster {
			apiserver = c.Cluster.Server
			return
		}
	}

	err = errors.New("Apiserver not found")
	return
}

func (k *KubeConfig) GetCurrentUser() (user *User, err error) {
	username,err := k.GetCurrentUsername()

	if err != nil {
		return
	}

	for _,user = range k.Users {
		if user.Name == username {
			return
		}
	}

	err = errors.New("User not found")
	return
}

func (k *KubeConfig) GetCurrentUserToken() ( token string, err error) {
	user,err := k.GetCurrentUser()
	if err != nil {
		return
	}

	token = user.User.Token
	return
}

func (k *KubeConfig) CreateInitalConfig(url string) (err error){
	req,_ := http.NewRequest("GET",url,nil)
	client := &http.Client{}
	resp,err := client.Do(req)

	if err != nil {
		return
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	if resp.Status != "200 OK" {
		err = errors.New("Unable to get initial config")
		return
	}

	initialConfig := KubeConfig{}
	yaml.Unmarshal(bodyText,&initialConfig)

	var usernames map[string]string
	var clusternames map[string]string
	var contextnames map[string]string

	for _,u := range k.Users {
		usernames[u.Name] = u.Name
	}
	for _,cl := range k.Clusters {
		clusternames[cl.Name] = cl.Name
	}
	for _,c := range k.Contexts {
		contextnames[c.Name] = c.Name
	}

	for _,user := range initialConfig.Users {
		if _,ok := usernames[user.Name]; ! ok {
			k.Users = append(k.Users, user)
		}
	}
	for _,cluster := range initialConfig.Clusters {
		if _,ok := clusternames[cluster.Name]; ! ok {
			k.Clusters = append(k.Clusters, cluster)
		}
	}
	for _,context := range initialConfig.Contexts {
		if _,ok := contextnames[context.Name]; ! ok {
			k.Contexts = append(k.Contexts, context)
		}
	}
	k.CurrentContext = initialConfig.CurrentContext

	return nil
}

func (k *KubeConfig) Save(filename string) (err error) {
  _,err = os.Stat(filename)
  if err != nil {
    _,err = os.Create(filename)
    if err != nil {
      return
    }
  }
	d := []byte(fmt.Sprintf("%s",k))
	err = ioutil.WriteFile(filename,d,0644)
	return
}