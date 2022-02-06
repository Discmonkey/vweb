package vars

import (
	"os"
	"strconv"
)

func HttpStaticDir() string {
	value, ok := os.LookupEnv("STATIC_HTTP_DIR")
	if !ok {
		// TODO(should be relative path to root of project)
		value = "client/dist"
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
