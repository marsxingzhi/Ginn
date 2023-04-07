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
		// router := ginn.New()

		// router.GET("/", defaultHandler)
		// router.GET("/hello", helloHandler)

		// router.Run(":8002")
	}

	{
		// router := ginn.New()

		// router.GET("/ping", func(ctx *ginn.Context) {
		// 	ctx.JSON(http.StatusOK, ginn.H{
		// 		"status":  0,
		// 		"message": "ok",
		// 	})
		// })

		// router.Run(":8003")
	}

	{
		// router := ginn.New()
		// router.GET("/ping", func(ctx *ginn.Context) {
		// 	ctx.JSON(http.StatusOK, ginn.H{
		// 		"status":  0,
		// 		"message": "OK",
		// 	})
		// })

		// v1 := router.Group("/v1")

		// {
		// 	v1.GET("/user/register", func(ctx *ginn.Context) {
		// 		ctx.JSON(http.StatusOK, ginn.H{
		// 			"status":  0,
		// 			"message": "this is response of /v1/user/register",
		// 		})
		// 	})
		// }

		// router.Run(":8004")
	}

	{
		router := ginn.New()
		router.Use(mockGlobalLoggerMiddleware)

		router.GET("/ping", func(ctx *ginn.Context) {
			ctx.JSON(http.StatusOK, ginn.H{
				"status":  0,
				"message": "OK",
			})
		})

		v1 := router.Group("/v1")
		v1.Use(mockLoggerMiddlewareV1)

		{
			v1.GET("/user/register", func(ctx *ginn.Context) {
				ctx.JSON(http.StatusOK, ginn.H{
					"status":  0,
					"message": "this is response of /v1/user/register",
				})
			})
		}

		v2 := router.Group("/v2")
		v2.Use(mockLoggerMiddlewareV2)
		{
			v2.GET("/user/register", func(ctx *ginn.Context) {
				ctx.JSON(http.StatusOK, ginn.H{
					"status":  0,
					"message": "this is response of /v2/user/register",
				})
			})
		}

		router.Run(":8004")
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

func mockGlobalLoggerMiddleware(ctx *ginn.Context) {
	fmt.Println("this is mock logger middleware")
}

func mockLoggerMiddlewareV1(ctx *ginn.Context) {
	fmt.Println("this is mock logger middleware v1")
}

func mockLoggerMiddlewareV2(ctx *ginn.Context) {
	fmt.Println("this is mock logger middleware v2")
}
