# kubelogin

Get a kubectl token for a kubernetes api when using a webhook as authorization method.

## Build

```
make
```

## Usage

```
$ ./kubectllogin help
Usage: ./kubectllogin [--kube-config] [--auth-server] COMMAND
    Command:
                 login    :   Default, do a login to kube webhook server
                 init     :   Create an initial kubectl config file
                 whoami   :   Check your current credentials
                 token    :   Display you current token data
                 help     :   Display this help

  -auth-server string
    	Authentication webhook servers (default "auto-detect")
  -kubeconfig string
    	kube config file location (default "/Users/skippie/.kube/config")
```

## License

https://github.com/skippie81/kubelogin/blob/master/LICENSE