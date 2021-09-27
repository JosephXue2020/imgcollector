package browser

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

var command = map[string]string{
	"windows": "cmd",
	"linux":   "xdg-open",
	"darwin":  "open",
}

func Open(url string) error {
	runningOS := runtime.GOOS
	run, ok := command[runningOS]
	if !ok {
		return fmt.Errorf("Don't know how to open browser on %s platform.", runningOS)
	}
	cmd := exec.Command(run, "/C", "start", url)
	err := cmd.Start()
	return err
}

func RetardOpen(url string) error {
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	if code := resp.StatusCode; code != 200 {
		return err
	}
	return Open(url)
}
