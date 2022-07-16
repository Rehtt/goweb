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

func (g *gweb) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	match, handleFunc, grep := g.PathMatch(request.RequestURI, request.Method)
	if handleFunc == nil {
		http.Redirect(writer, request, "/404", 404)
		return
	}

	ctx := &Context{
		Request: request,
		Writer:  writer,
		param:   match,
		survive: true,
	}
	for grep != nil {
		for i := range grep.middlewares {
			ctx.runFunc(grep.middlewares[i])
		}
		grep = grep.parent
	}

	ctx.runFunc(handleFunc)
}

func New() *gweb {
	return new(gweb)
}
