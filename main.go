package main

import (
	"fmt"

	"projects/imgcollector/browser"
	"projects/imgcollector/webserver"
)

func main() {

	// // 测试写入docx
	// pth := "D:\\workdir\\一般性的xml数据处理\\20210923任务-韩晓玲\\test\\result.docx"
	// para1 := office.Para{Typ: "text", Text: "abcdergajgklajglgfjakglagfjlak"}
	// imgPth := "D:\\workdir\\一般性的xml数据处理\\20210923任务-韩晓玲\\test\\test.png"
	// para2 := office.Para{Typ: "image", Pth: imgPth}
	// imgPth3 := "D:\\workdir\\一般性的xml数据处理\\20210923任务-韩晓玲\\test\\resized.jpg"
	// para3 := office.Para{Typ: "image", Pth: imgPth3}
	// paras := []office.Para{para1, para2, para3}
	// office.WriteDocxFile(pth, paras)

	// imgoutpath := "D:\\workdir\\一般性的xml数据处理\\20210923任务-韩晓玲\\test\\resized.jpg"
	// imgresize.ResizeImageFile(imgPth, imgoutpath)

	// // 测试
	// urltest := "http://media1.zbk100.com/image/repository/api/objdata/10052002000000088073787168824/0"
	// header := new(downloader.Header)
	// header.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"
	// sleepMilliSec := 500

	// resp, _ := downloader.SlowGet(urltest, header, sleepMilliSec)
	// fmt.Print(resp.Header)

	// 打开浏览器
	port := 10005
	url := fmt.Sprintf("http://localhost:%v", port)
	go openBrowser(url)

	// 运行
	err := webserver.Run(port)
	if err != nil {
		panic(err)
	}
}

func openBrowser(url string) {
	err := browser.RetardOpen(url)
	if err != nil {
		panic(err)
	}
}
