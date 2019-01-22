package holder

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
	"runtime"
	"strings"
	"github.com/labstack/echo"
)

var (
	dir      = flag.String("d", ".", "-d dirname ")
	port     = flag.Int("p", 9090, "-p port")
	username = flag.String("u", "", "username")
	password = flag.String("m", "", "password")
	host     = flag.String("h", "127.0.0.1", "host")
	token    string
)

func init() {
	flag.Parse()

	fmt.Println("使用方法：")
	fmt.Println("进入pc-center目录")
	fmt.Println("./pc-center -p 9090 -d .")
	fmt.Println("-u 前台用户名")
	fmt.Println("-m 前台密码")
	fmt.Println("-h 前台地址  默认开发环境 aaaa.dev.front.me")
	fmt.Println("./pc-center -p 9090 -d . -u username -m password -h 127.0.0.1")

	if *username != "" && *password != "" {
		login()
	}

}

func main() {

	Open(fmt.Sprintf("http://localhost:%d/member/index.html", *port))

	e := echo.New()
	e.Static("/static", "./static")
	e.Any("/member/*", func(context echo.Context) error {
		Handle(context.Response(), context.Request())
		return nil
	})

	e.Any("/api/*", func(context echo.Context) error {
		reverseProxy(context.Response(), context.Request())
		return nil
	})

	e.POST("/api/member/nav", func(context echo.Context) error {
		b, _ := ioutil.ReadAll(context.Request().Body)
		var mp = make(map[string]string)
		json.Unmarshal(b, &mp)

		d, err := filepath.Abs(*dir)
		uri := mp["position"]
		b, err = ioutil.ReadFile(filepath.Join(d, "member", uri+".html"))

		fmt.Println(mp)
		if err != nil {
			context.Response().Write([]byte(err.Error()))
			return nil
		}

		tpl := template.New(uri)
		tpl = tpl.Funcs(map[string]interface{}{
			"include": include,
		})
		tpl = tpl.Delims("<{", "}>")
		tpl, err = tpl.Parse(string(b))
		tpl.Execute(context.Response(), map[string]interface{}{
			"MemUrl":    fmt.Sprintf("http://localhost:%d/static", *port),
			"PublicUrl": "http://192.168.11.200:9091/public",
		})

		return nil
	})

	e.Start(fmt.Sprintf(":%d", *port))
}

func Handle(writer http.ResponseWriter, request *http.Request) {
	d, err := filepath.Abs(*dir)
	uri := request.RequestURI
	b, err := ioutil.ReadFile(filepath.Join(d, uri))
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	tpl := template.New(uri)
	tpl = tpl.Funcs(map[string]interface{}{
		"include": include,
	})
	tpl = tpl.Delims("<{", "}>")
	tpl, err = tpl.Parse(string(b))
	tpl.Execute(writer, map[string]interface{}{
		"MemUrl":    fmt.Sprintf("http://localhost:%d/static", *port),
		"PublicUrl": "http://192.168.11.200:9091/public",
	})
}

var (
	commands = map[string]string{
		"windows": "cmd /c start",
		"darwin":  "open",
		"linux":   "xdg-open",
	}
)

func Open(uri string) {
	/*run, ok := commands[runtime.GOOS]
	if !ok {
		fmt.Errorf("未知的操作系统%s，请手动打开 http://localhost:%d/", runtime.GOOS, port)
		return
	}
	cmd := exec.Command(run, uri)
	err := cmd.Start()
	if err != nil {
		fmt.Errorf("打开浏览器失败,请手动打开 http://localhost:%d/", port)
	}*/
}

func include(partials string, path string) (h template.HTML, err error) {
	d, _ := filepath.Abs(*dir)
	fmt.Println(d)
	file := filepath.Join(d, path) + ".tpl"

	fmt.Println(file)
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}
	tpl := template.New(file)
	tpl.Funcs(map[string]interface{}{
		"include": include,
	})
	tpl = tpl.Delims("<{", "}>")
	tpl, err = tpl.Parse(string(b))

	var buf = bytes.NewBuffer(nil)
	err = tpl.Execute(buf, map[string]interface{}{
		"MemUrl":    fmt.Sprintf("http://localhost:%d/static", *port),
		"PublicUrl": "http://192.168.11.200:9091/public",
	})

	if err != nil {
		log.Println(err)
	}
	return template.HTML(buf.String()), err
}

// 反向代理到服务器.
func reverseProxy(w http.ResponseWriter, r *http.Request) {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = *host
		req.Host = *host
		req.AddCookie(&http.Cookie{
			Name:  "sid",
			Value: token,
			Path:  "/",
		})
	}
	defer func() {
		if err := recover(); err != nil {
			b := make([]byte, 1024)
			n := runtime.Stack(b, true)

			fmt.Println(string(b[:n]))
		}
	}()
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func login() {
	var params = fmt.Sprintf(`{"account":"%s","password":"%s","code":"1234"}`, *username, *password)
	// 请求验证码

	resp, err := http.Get(fmt.Sprintf("http://%s/captcha", *host))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var captcha = ""
	ck := resp.Cookies()
	for _, v := range ck {
		if v.Name == "captcha" {
			captcha = v.Value
		}
	}

	if captcha == "" {
		panic("验证码获取失败")
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api/login", *host), strings.NewReader(params))
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "captcha",
		Value: captcha,
		Path:  "/",
	})

	fmt.Println(params, "验证码key", captcha)

	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	cks := resp.Cookies()
	for _, v := range cks {
		if v.Name == "sid" {
			token = v.Value
		}
	}
	if token == "" {
		panic("登录失败")
	}

	fmt.Println("登录成功：sid", token)
}

