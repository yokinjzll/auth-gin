# Auth-gin
##### _Reginster and login on golang gin_

### _Current version have many bugs!!!_

## Features

- routes: Register, Login, Dashboard, Profile
- change profile info
- on Bootstrap 5
- use framefork Gin on Golang

 
## Get started

1. First, make migrate for your database
```sh
$ migrate -path db/migration -database "mysql://<db_user>:<user_password>@tcp(localhost:3306)/<db_name>?multiStatements=true" -verbose up
```
2. Or delete all migration in ```db/migration``` and create new:
```sh
$ migrate create -ext sql -dir db/migration -seq <name_migration> 
```
3. Set environment variable DNS for connect to database:
```sh
export DNS=<db_user>:<user_password>@/<db_name>?Time=true
```

Next launch server and enjoy:
- enter in root derictory of project
- ```go run ./cmd```

