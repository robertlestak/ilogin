# ilogin - Istio Login CLI

A small web app and cli utility to enable retrieving auth cookies from the browser for services backed by authorizing middleware such as [istio-oauth2-sso](https://github.com/robertlestak/istio-oauth2-sso).

## Usage

```bash
ilogin
  -auth string
        auth url
  -cookie string
        cookie name
  -f string
        output file
  -server
        run server mode
  -service string
        service url
```

## rc file

You can create a `~/.iloginrc` file with the following content:

```
# iloginrc
#
# iloginrc is a config file for ilogin
#
# The following variables are available:
#
auth https://login.example.com/oauth2/12345
service https://ilogin.example.com
cookie oauth2_sso
out_file ~/.ilogin-token
```