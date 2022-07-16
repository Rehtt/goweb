/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/7/16 下午 04:30
 */

package gweb

import (
	"net/http"
)

type gweb struct {
	RouterGroup
}

// todo 细化，修bug
func (g *gweb) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match, handleFunc, grep := g.PathMatch(request.RequestURI, request.Method)
	if handleFunc == nil {
		http.Redirect(writer, request, "/404", 404)
		return
	}
	for _, m := range grep.middlewares {
		m.ServeHTTP(writer, request)
	}
	handleFunc.ServeHTTP(writer, request)
}

func New() *gweb {
	return new(gweb)
}
