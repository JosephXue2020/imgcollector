package webserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"projects/imgcollector/docgenerator"
	"projects/imgcollector/office"

	"baliance.com/gooxml/document"
)

// temparary directory for downloaded images
func initPath() (string, string) {
	dir, _ := os.Getwd()

	tempDir := filepath.Join(dir, "temp")
	_, err := os.Stat(tempDir)
	if err != nil {
		os.Mkdir(tempDir, 0777)
	}

	docxPath := filepath.Join(dir, "result.docx")

	return tempDir, docxPath
}

// index page handle funcion
func index(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "GET" {
		fmt.Fprintf(w, getIndex())
	}
	if method == "POST" {
		r.ParseForm()
		direcSli := r.Form["direc"]
		var direc string

		var alert string
		if len(direcSli) == 0 {
			alert = "没有路径被输入"
		} else if len(direcSli) > 1 {
			alert = "输入路径有误"
		} else if len(direcSli) == 1 {
			direc = direcSli[0]
			fInfo, err := os.Stat(direc)
			if err != nil {
				alert = "输入路径不存在"
			} else if !fInfo.IsDir() {
				alert = "输入路径不是文件夹"
			}
		}

		if alert != "" {
			fmt.Fprintf(w, getAlert(alert))
		} else {
			tempdir, docxPath := initPath()
			Result = docgenerator.GenDocx(direc, docxPath, tempdir)
			fmt.Fprintf(w, getReply())
		}
	}

}

// result
var Result *document.Document

// download page handle function
func downloadResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Disposition", "attachment; filename=result.docx")
	office.WriteDocxToWriter(w, Result)
}
