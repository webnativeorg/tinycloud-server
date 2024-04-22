package environment

import (
	"os"
)

var MONGO_URL string = os.Getenv("MONGO_URL")
var MONGO_DATABASE string = os.Getenv("MONGO_DATABASE")
var MONGO_COLLECTION_PREFIX string = os.Getenv("MONGO_COLLECTION_PREFIX")
