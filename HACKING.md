# Hacking

## Template
If you change a Template or want to contribute one of your own templates to the upstream repo you will need to regenerate the *default-templates-bindata.go* file.

For this we use [go-bindata](https://github.com/jteeuwen/go-bindata):

```
go get -u github.com/jteeuwen/go-bindata/...
go-bindata -o default-templates-bindata.go templates
```
