package main

import (
	yaml "gopkg.in/yaml.v2"
	"flag"
	"os"
	"fmt"
	"io/ioutil"

	"kubelogin/kubeconfig"
	"kubelogin/client"
	"kubelogin/autoupdate"

	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"runtime"
	"syscall"
)

const (
	DefaultKubeConfig = ".kube/config"
	loginPath = "/ldapAuth"
	validatePath = "/authenticate"
	initialConfigUrl = "https://updateserver.local/initialconfig"
	autoUpdateVersionUrl = "https://raw.githubusercontent.com/skippie81/kubelogin/master/VERSION"
	autoUpdateGetUrl = "https://updatserver.local/kubelogin"
	myVersion = "2.0.0"
)

func main() {

	// get homedir
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}

	// input flags
	kubeconfigfileFlag := flag.String("kubeconfig",home + "/" + DefaultKubeConfig,"kube config file location")
	loginServerFlag := flag.String("auth-server","auto-detect","Authentication webhook servers")
	versionFlag := flag.Bool("version",false,"display version and exit")
	flag.Parse()

	// print version and exit
	if *versionFlag {
		fmt.Printf("Kubelogin version: v%s ( %s )\n",myVersion,runtime.Version())
		os.Exit(0)
	}

	// autoupdater
	upd := autoupdate.CreateNew(autoUpdateVersionUrl)
	update,version,err := upd.Check(myVersion)
	if err == nil {
		if update {
			fmt.Printf("Auto Udating: %s => %s\n", myVersion, version)
			goos := runtime.GOOS
			myfilename := os.Args[0]
			url := autoUpdateGetUrl + "/v" + version + "/bin/" + goos + "/kubelogin"
			fmt.Printf("Downloading: %s ... ",url)
			err = upd.Update(url, myfilename)
			if err != nil {
				fmt.Printf("Error while updating: %s\n",err)
			}
			fmt.Printf("updated\n")
		}
	}

	initializeCmd := flag.NewFlagSet("init",flag.ExitOnError)
	loginCmd := flag.NewFlagSet("login",flag.ExitOnError)
	whoamiCmd := flag.NewFlagSet("whoami",flag.ExitOnError)
	tokenCmd := flag.NewFlagSet("token",flag.ExitOnError)

	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "init":
			initializeCmd.Parse(flag.Args()[1:])
		case "whoami":
			whoamiCmd.Parse(flag.Args()[1:])
		case "token":
			tokenCmd.Parse(flag.Args()[1:])
		case "help":
			doHelp(os.Args[0])
			flag.PrintDefaults()
			os.Exit(0)
		default:
			loginCmd.Parse(flag.Args()[1:])
		}
	}

	k := kubeconfig.KubeConfig{}
	// read the kubeconfig file
	data,err := ioutil.ReadFile(*kubeconfigfileFlag)
	if err != nil && ! initializeCmd.Parsed()  {
		fmt.Printf("WARN: Kubeconfig file %s not found, creating initial config with %s init\n",*kubeconfigfileFlag,os.Args[0])
		os.Exit(0)
	} else {
		err = yaml.Unmarshal(data, &k)
		if err != nil {
			fmt.Println("Error parsing data in kubeconfig")
			os.Exit(1)
		}
	}

	// create inital kubeconfig with cmd init
	if initializeCmd.Parsed() {
		err = k.CreateInitalConfig(initialConfigUrl)
		if err != nil {
			fmt.Printf("Error creating inital config: %s\n",err)
			os.Exit(1)
		}
		err = k.Save(*kubeconfigfileFlag)
		if err != nil {
			fmt.Printf("Error saving %s: %s\n",kubeconfigfileFlag,err)
		}
		os.Exit(0)
	}

	// print current token
	if tokenCmd.Parsed() {
		currentToken,err := k.GetCurrentUserToken()
		if err != nil {
			fmt.Println("Error getting current token")
			os.Exit(1)
		}
		fmt.Println("Your current token:")
		fmt.Printf("%s\n",currentToken)
		os.Exit(0)
	}


	// search current apiserver in kubecconfig
	apiserver,err := k.GetCurrentAPIServer()
	if err != nil {
		fmt.Printf("%s\n",err)
	}

	// search loginserver
	var loginserver string
	if *loginServerFlag != "auto-detect" {
		loginserver = *loginServerFlag
	} else {
		loginserver, err = kubeconfig.GetLoginServer(apiserver)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}

	// create client
	c := client.Client{
		Server: loginserver,
		Path:   loginPath,
		ValidatePath: validatePath,
	}

	if whoamiCmd.Parsed() {
		currentToken,err := k.GetCurrentUserToken()
		if err != nil {
			os.Exit(1)
		}
		trs,err := c.Validate(currentToken)
		if err != nil {
			fmt.Printf("Error validating token: %s\n",err)
			os.Exit(1)
		}
		if trs.Status.Authenticated {
			fmt.Printf("Username: %s\n", trs.Status.User.Username)
			fmt.Printf("Groups: ")
			for _, group := range trs.Status.User.Groups {
				fmt.Printf("[%s]  ",group)
			}
			fmt.Printf("\n")
			//fmt.Printf("Token expire time: %s\n",trs.Status.User.Extra["Exp"])
		}
		os.Exit(0)
	}

	fmt.Printf("Current login server: %s\n",loginserver)

	// read username and password
	var username,password string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ = reader.ReadString('\n')
	//username = strings.Replace(username, "\n", "", -1)
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	p,_ := terminal.ReadPassword(int(syscall.Stdin))
	password = string(p)
	password = strings.TrimSpace(password)

	// autenticate and get token
	token,err := c.Authenticate(string(username),string(password))

	if err != nil {
		fmt.Printf("ERROR while authenticating: %s\n",err)
		os.Exit(1)
	}

	// testing new token
	trr,err := c.Validate(token)
	if err != nil {
		fmt.Printf("Error wile validating new token: %s\n",err)
		os.Exit(1)
	}

	if ! trr.Status.Authenticated {
		fmt.Println("Error: new token validation returend NOT AUTHENTICATED")
		os.Exit(1)
	}

	user,err := k.GetCurrentUser()
	if err != nil {
		fmt.Println("Error getting current user in kubecofnig to store token")
		os.Exit(1)
	}

	user.User.Token = token
	fmt.Printf("Saving new token ... ")
	err = k.Save(*kubeconfigfileFlag)
	if err != nil {
		fmt.Printf("FAIL (%s)",err)
		os.Exit(1)
	}
	fmt.Printf("ok\n")
}

func doHelp(cmdline string) {
	fmt.Printf("Usage: %s [--kube-config] [--auth-server] COMMAND\n",cmdline)
	fmt.Println("    Command:")
	fmt.Println("                 login    :   Default, do a login to kube webhook server")
	fmt.Println("                 init     :   Create an initial kubectl config file")
	fmt.Println("                 whoami   :   Check your current credentials")
	fmt.Println("                 token    :   Display you current token data")
	fmt.Println("                 help     :   Display this help")
	fmt.Println("")
}