# Auth-gin
##### _Register and login on golang gin_


|  ## _Current version have many bugs!!!_   |
---------------------------------------------

## Features

- routes: Register, Login, Dashboard, Profile
- change profile info
- on Bootstrap 5
- use framefork Gin on Golang
- connect to MySQL

 
## Get started

1. First, make migrate for your database
```sh
$ migrate -path db/migration -database "mysql://<db_user>:<user_password>@tcp(localhost:3306)/<db_name>?multiStatements=true" -verbose up
```
2. Or delete all migration in ```db/migration``` and create new:
```sh
$ migrate create -ext sql -dir db/migration -seq <name_migration> 
```
3. Set environment variable DSN for connect to database:
```sh
export DSN=<db_user>:<user_password>@/<db_name>?parseTime=true
```

Next launch server and enjoy:
- enter in root derictory of project
- ```go run ./cmd```

