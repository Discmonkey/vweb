package vars

import (
	"os"
	"strconv"
)

func HttpStaticDir() string {
	value, ok := os.LookupEnv("REWINDER_HTTP_DIR")
	if !ok {
		// TODO(should be relative path to root of project)
		value = "/home/max/go/src/vweb/client/dist"
	}

	return value
}

func HttpServerPort() int {
	strValue, ok := os.LookupEnv("REWINDER_HTTP_PORT")
	if !ok {
		return 3000
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 3000
	}

	return intValue
}
