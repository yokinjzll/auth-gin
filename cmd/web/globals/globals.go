package globals

import (
	"os"
)

var Secret = []byte("secret")

const Userkey = "user"
const UserID = "user_id"
const UserDetail = "user_info"

var Dsn = os.Getenv("DSN")
