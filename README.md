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
  -version
    	display version and exit
```

## Remarks

before building you might want to change in cmd/kubelogin.go:

```
initialConfigUrl = "<url where your default kubectl config can be found>"
autoUpdateVersionUrl = "<url that returns the latest version nr>"
autoUpdateGetUrl = "<url that holds bindaries>"
```

The autoUpdateGetUrl is the base url that holds kubelogin binaries for autoupdate at path \<base url\>/v\<version\>/\<operating system\>/kubelogin. Exmple https://example.com/kubelogin as base url would hold https://example.com/kubelogin/v2.0.0/darwin/kubelogin

If autoUpdateVersion url is not set or http get would return error, autoupdate is disabled.

## License

https://github.com/skippie81/kubelogin/blob/master/LICENSE