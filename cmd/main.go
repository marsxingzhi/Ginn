package main

import (
	"fmt"
	"net/http"

	"marsxingzhi.github.com/ginn/ginn"
)

func main() {
	fmt.Println("hello world")

	// 使用http构造请求
	{
		// 注册handler function到DefaultServeMux
		// http.HandleFunc("/", defaultHandler)
		// http.HandleFunc("/hello", helloHandler)

		// err := http.ListenAndServe(":8888", nil)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	{
		// handleFunc是路由与handler的映射
		// http.HandleFunc("/", defaultHandler)
		// http.HandleFunc("/hello", helloHandler)

		// router := ginn.New()

		// err := http.ListenAndServe(":8001", router)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	{
		router := ginn.New()

		router.GET("/", defaultHandler)
		router.GET("/hello", helloHandler)

		router.Run(":8002")
	}

}

// %q：该值对应的单引号括起来的go语法字符字面值，必要时会采用安全的转义表示
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%v] = %v\n", k, v)
	}
}