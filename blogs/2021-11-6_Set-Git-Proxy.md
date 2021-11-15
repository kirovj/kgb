# Correct way to set git proxy

## 1. When use https://github.com/user/respository.git

```ssh
git config –global http.proxy protocol://127.0.0.1:port
```
> protocol & port is your proxy

## 2. For the specified domain

```ssh
> git config –global http.url.proxy protocol://127.0.0.1:port
```
> example: http.https://github.com.proxy

## 3. When use ssh

```ssh
vim ~/.ssh/config
```
```
Host github.com
    User git
    ProxyCommand connect -H 127.0.0.1:port %h %p
```

> thanks to [ericclose](https://ericclose.github.io/git-proxy-config.html)