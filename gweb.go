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
	noRouter HandlerFunc
}

func (g *gweb) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Request: request,
		Writer:  writer,
		survive: true,
	}
	match, handleFunc, grep := g.PathMatch(request.RequestURI, request.Method)
	if handleFunc == nil {
		g.handler404(ctx)
		return
	}

	ctx.param = match

	for grep != nil {
		for i := range grep.middlewares {
			ctx.runFunc(grep.middlewares[i])
		}
		grep = grep.parent
	}

	ctx.runFunc(handleFunc)
}

func (g *gweb) NoRoute(handlerFunc HandlerFunc) {
	g.noRouter = handlerFunc
}
func (g *gweb) handler404(ctx *Context) {
	ctx.Writer.WriteHeader(http.StatusNotFound)
	if g.noRouter != nil {
		g.noRouter(ctx)
	} else {
		http.NotFound(ctx.Writer, ctx.Request)
	}
}

func New() *gweb {
	return new(gweb)
}
