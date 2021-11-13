
> Simple web project on golang with [gin](https://github.com/gin-gonic/gin)

[![title.png](title.png)](#title)

### About
The project is created to study the web capabilities of the golang gin web framework.

### The project used:
- bootstrap 5
- font awesome 4
- sweetalert2

### Implemented in the project:
- authorization and registration pages,
- role-based system with three roles: user, editor and administrator,
- updating user data, uploading an avatar
- for the role administrator, management (editing, blocking, deleting and registering) users is available
- blog section: manage (edit, create and delete) categories and posts
- editing categories and posts is available only for editor and administrator
- blog categories can be hidden for unauthorized users
- post editing is available using summernote editor 

### Project setup
- copy project to GOPATH/src directory
- copy config file [config.empty.json](config.empty.json) to config.json
- set connection to your mysql server in config.json db section
- run database migrations "go-gin-web migrate up", migrations comes with 3 users "admin", "editor" and "user", all has password "qwe"

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

### Database migrations
```bash
// apply all migrations, apply new migrations
go-gin-web migrate up
// roll back one last migration
go-gin-web migrate down
```

### Bootstrap 5 theme
Mazer: https://github.com/zuramai/mazer
