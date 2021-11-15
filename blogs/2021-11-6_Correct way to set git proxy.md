# Correct way to set git proxy

### _When use https://github.com/user/respository.git_

```ssh
git config –global http.proxy protocol://127.0.0.1:port
```
> protocol & port is your proxy

## _For the specified domain_

```ssh
> git config –global http.url.proxy protocol://127.0.0.1:port
```
> example: http.https://github.com.proxy

## _When use ssh_

```ssh
vim ~/.ssh/config
```
```
Host github.com
    User git
    ProxyCommand connect -H 127.0.0.1:port %h %p
```

> thanks to [ericclose](https://ericclose.github.io/git-proxy-config.html)