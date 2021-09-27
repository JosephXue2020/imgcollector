package webserver

import (
	"fmt"
	"net/http"
	"strconv"
)

// Run function runs a web server.
func Run(port int) error {

	// regist handle funcion
	routerInit()

	// run service
	portStr := strconv.Itoa(port)
	host := fmt.Sprintf(":%s", portStr)
	return http.ListenAndServe(host, nil)
}

func routerInit() {
	// index page
	http.HandleFunc("/", index)

	// download page
	http.HandleFunc("/download", downloadResult)
}
