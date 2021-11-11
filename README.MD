
> Simple web project on golang with [gin](https://github.com/gin-gonic/gin)

[![title.png](title.png)](#title)

### Project setup
> copy project to GOPATH/src directory

> copy config file [config.empty.json](config.empty.json) to config.json

> set connection to your mysql server in config.json db section

> run database migrations "go-gin-web migrate up", migrations comes with 3 users "admin", "editor" and "user", all has password "qwe"

### Install dependencies
```bash
go get
```

### Build
```bash
go build
```

### Run migrations
```bash
go-gin-web migrate up
```

### Bootstrap 5 theme
Mazer: https://github.com/zuramai/mazer
