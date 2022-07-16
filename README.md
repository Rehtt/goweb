一个简单的web api框架

支持中间件，路由编写更友好

```go
package main

import (
	"fmt"
	"gweb"
	"net/http"
)

func main() {
	web := gweb.New()
	web.Middleware(func(ctx *gweb.Context) {
		fmt.Println("中间件")
	})

	web.Any("/123/#asd/234", func(ctx *gweb.Context) {
		fmt.Println(ctx.GetParam("asd"), "获取动态路由参数")
	})
	api := web.Grep("/api")
	api.GET("/test", func(ctx *gweb.Context) {
		fmt.Println("test")
	})

	http.ListenAndServe(":9090", web)
}

```
