/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/7/16 下午 05:02
 */

package gweb

import (
	"net/http"
	"strings"
)

type RouterGroup struct {
	index       int
	middlewares []http.Handler
	path        string
	method      map[string]http.HandlerFunc
	child       map[string]*RouterGroup
	parent      *RouterGroup
}

func (g *RouterGroup) Grep(path string) *RouterGroup {
	return g.position(path)
}
func (g *RouterGroup) Middleware(handlers ...http.Handler) {
	if len(g.middlewares) == 0 {
		g.middlewares = make([]http.Handler, 0, len(handlers)+5)
	}
	g.middlewares = append(g.middlewares, handlers...)
}

func (g *RouterGroup) position(path string) *RouterGroup {
	for _, p := range strings.Split(path, "/") {
		if p == "" {
			continue
		}

		if _, ok := g.child[p]; ok {
			g = g.child[p]
			continue
		}
		if g.child == nil {
			g.child = make(map[string]*RouterGroup, 10)
		}
		if p[0] == '#' {
			for child := range g.child {
				if child[0] == '#' {
					panic(path + " 地址泛匹配重复")
				}
			}
		}

		g.child[p] = &RouterGroup{
			path:   p,
			parent: g,
			index:  g.index + 1,
		}
		g = g.child[p]
	}
	return g
}

func (g *RouterGroup) completePath() string {
	completePath := make([]string, g.index)
	for g.index != 0 {
		completePath[g.index-1] = g.path
		g = g.parent
	}
	return "/" + strings.Join(completePath, "/")
}
func (g *RouterGroup) handle(method string, path string, handlerFunc http.HandlerFunc) {
	g = g.position(path)
	if method == ANY && len(g.method) > 1 {
		if _, ok := g.method[ANY]; ok {
			panic(g.completePath() + "该路由any方法冲突")
		}
	}
	if _, ok := g.method[method]; ok {
		panic(g.completePath() + "该路由method重复")
	}
	if g.method == nil {
		g.method = make(map[string]http.HandlerFunc, 5)
	}
	g.method[method] = handlerFunc
}
func (g *RouterGroup) GET(path string, handlerFunc http.HandlerFunc) {
	g.handle(GET, path, handlerFunc)
}
func (g *RouterGroup) POST(path string, handlerFunc http.HandlerFunc) {
	g.handle(POST, path, handlerFunc)
}
func (g *RouterGroup) PUT(path string, handlerFunc http.HandlerFunc) {
	g.handle(PUT, path, handlerFunc)
}
func (g *RouterGroup) DELETE(path string, handlerFunc http.HandlerFunc) {
	g.handle(DELETE, path, handlerFunc)
}
func (g *RouterGroup) Any(path string, handlerFunc http.HandlerFunc) {
	g.handle(ANY, path, handlerFunc)
}

func (g *RouterGroup) PathMatch(path, method string) (match map[string]string, handle http.HandlerFunc, group *RouterGroup) {
	match = make(map[string]string)
	paths := strings.Split(path, "/")
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, ok := g.child[p]; ok {
			g = g.child[p]
			continue
		}
		var find bool
		for child := range g.child {
			if child[0] == '#' {
				match[child[1:]] = p
				g = g.child[child]
				find = true
				break
			}
		}
		if !find {
			return nil, nil, nil
		}
	}
	h, ok := g.method[method]
	if !ok {
		return nil, nil, nil
	}
	return match, h, g
}
