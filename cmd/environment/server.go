package environment

import (
	"os"
	"strconv"
)

var PORT string = os.Getenv("PORT")
var JWT_SECRET string = os.Getenv("JWT_SECRET")
var JWT_EXPIRATION, _ = strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
