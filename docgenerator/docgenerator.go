package docgenerator

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"projects/imgcollector/downloader"
	"projects/imgcollector/imgresize"
	"regexp"
	"strconv"
	"strings"
	"time"

	"baliance.com/gooxml/document"

	"projects/imgcollector/office"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// getPathInfo function collects all the files in root directory
func getPathInfo(direc string) ([][]string, error) {
	var fInfo [][]string

	walkFunc := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		var inSli []string
		if info.IsDir() {
			return nil
		} else {
			_, fname := filepath.Split(p)
			ext := path.Ext(fname)
			inSli = []string{fname, ext, p}
		}
		fInfo = append(fInfo, inSli)
		return nil
	}

	err := filepath.Walk(direc, walkFunc)
	return fInfo, err
}

func xmlFilter(fInfo [][]string) [][]string {
	var res [][]string
	for _, item := range fInfo {
		ext := item[1]
		if ext == ".xml" {
			res = append(res, item)
		}
	}
	return res
}

func readStr(p string) (string, error) {
	bytes, err := ioutil.ReadFile(p)
	return string(bytes[:]), err
}

// meta data to describ each entry
type Meta struct {
	title           string
	urls            []string
	successDownload []bool
	localPath       []string
}

func getEntryTitle(xml string) string {
	pat := regexp.MustCompile("<subclauses.*?>")
	finds := pat.FindAllString(xml, -1)
	if len(finds) == 0 {
		err := fmt.Errorf("No title appeared in the entry.")
		panic(err)
	} else if len(finds) > 1 {
		err := fmt.Errorf("Multiple titles appeared in the entry.")
		panic(err)
	}
	seg := finds[0]

	pat = regexp.MustCompile("name=\"(.*?)\"")
	title := pat.FindStringSubmatch(seg)[1]
	return title

}

// unescape
func unescape(s string) string {
	s = strings.Replace(s, "&amp;", "&", -1)
	s = strings.Replace(s, "&gt;", ">", -1)
	s = strings.Replace(s, "&lt;", "<", -1)
	s = strings.Replace(s, "&quot;", "\"", -1)
	s = strings.Replace(s, "&nbsp;", " ", -1)
	return s
}

func getEntryBody(xml string) string {
	xml = unescape(xml)

	pat := regexp.MustCompile("<templateRelates templateType=\"3\" orderId=\"3\" flag=\"1\" status=\"0\">.*?</templateRelates>")
	finds := pat.FindAllString(xml, -1)
	if len(finds) != 1 {
		// err := fmt.Errorf("can not find body segment in xml.")
		return ""
	}
	seg := finds[0]

	pat = regexp.MustCompile("<titles>(.*?)</titles>")
	finds = pat.FindAllString(seg, -1)
	if len(finds) != 1 {
		// err := fmt.Errorf("can not find body segment in xml.")
		return ""
	}
	body := finds[0]
	return body
}

func imgFilter(imgs []string) []string {
	drops := []string{"data-type=\"postil\"", "data-latex=\""}
	res := []string{}
	for _, item := range imgs {
		flag := false
		for _, substr := range drops {
			if strings.Contains(item, substr) {
				flag = true
			}
		}
		if !flag {
			res = append(res, item)
		}
	}
	return res
}

func priorMatch(s string) string {
	// s is img tag segment
	pat := regexp.MustCompile("url=\"(.*?)\"")
	url := pat.FindStringSubmatch(s)
	if len(url) != 0 {
		return url[1]
	}

	pat = regexp.MustCompile("src=\"(.*?)\"")
	url = pat.FindStringSubmatch(s)
	if len(url) != 0 {
		return url[1]
	}

	return ""
}

func collectURL(body string) []string {
	res := make([]string, 0)

	pat := regexp.MustCompile("<img.*?>")
	finds := pat.FindAllString(body, -1)
	if len(finds) == 0 {
		return res
	}

	finds = imgFilter(finds)
	if len(finds) == 0 {
		return res
	}

	for _, item := range finds {
		url := priorMatch(item)
		if url != "" {
			res = append(res, url)
		}
	}

	return res
}

func parseXML(xml string) (string, []string) {
	// 获取title
	title := getEntryTitle(xml)

	// 获取body
	body := getEntryBody(xml)

	// 获取urls
	urls := collectURL(body)

	return title, urls
}

func getRandomStr(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getImgPath(tempdir, imgName string) string {
	if imgName != "" {
		return filepath.Join(tempdir, imgName)
	}
	imgName = getRandomStr(15) + ".jpg"
	return filepath.Join(tempdir, imgName)
}

//
func download(url string, tempdir string) (string, bool) {
	sleepMilliSec := 500
	resp, err := downloader.SlowGet(url, nil, sleepMilliSec)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	seg := resp.Header.Get("Content-disposition")
	pat := regexp.MustCompile("filename=\"(.*?)\"")
	finds := pat.FindStringSubmatch(seg)
	var imgName string
	if len(finds) > 1 {
		imgName = finds[1]
	} else {
		imgName = ""
	}

	imgPath := getImgPath(tempdir, imgName)

	err = ioutil.WriteFile(imgPath, body, 0666)
	if err != nil {
		return "", false
	}

	return imgPath, true
}

func getParas(meta []Meta) []office.Para {
	var paras []office.Para
	for _, m := range meta {
		// 条头
		para := office.Para{Typ: "text", Text: "条目名称：" + m.title}
		paras = append(paras, para)

		// 图片
		for i, flag := range m.successDownload {
			if !flag {
				para := office.Para{Typ: "text", Text: "下载图片失败：" + m.urls[i]}
				paras = append(paras, para)
			} else {
				// 加入编号段
				para := office.Para{Typ: "text", Text: "图片序号：" + strconv.Itoa(i+1)}
				paras = append(paras, para)
				// 图片段
				// 调整尺寸
				inpath := m.localPath[i]
				dir, fname := filepath.Split(inpath)
				base := path.Base(fname)
				outpath := filepath.Join(dir, base+"_thumbnail"+".jpg")
				err := imgresize.ResizeImageFile(inpath, outpath)
				if err != nil {
					outpath = inpath
				}
				para = office.Para{Typ: "image", Pth: outpath}
				paras = append(paras, para)
			}
		}
	}
	return paras
}

func GenDocx(direc string, docxPath string, tempdir string) *document.Document {
	// 获取所有文件信息
	fInfo, _ := getPathInfo(direc)

	// 筛选出xml文件
	fInfo = xmlFilter(fInfo)

	metas := []Meta{}

	// 主循环
	for _, item := range fInfo {
		// 解析entry xml
		fp := item[2]
		xml, err := readStr(fp)
		if err != nil {
			panic(err)
		}
		t, urls := parseXML(xml)

		// 下载图片至本地目录
		var picPaths []string
		var flags []bool
		for _, url := range urls {
			picPath, flag := download(url, tempdir)
			picPaths = append(picPaths, picPath)
			flags = append(flags, flag)
		}
		mt := Meta{title: t, urls: urls, successDownload: flags, localPath: picPaths}
		metas = append(metas, mt)
	}

	// 生成docx
	paras := getParas(metas)
	// office.WriteDocxFile(docxPath, paras)
	// return document.New()

	return office.WriteDocxDocument(paras)
}
