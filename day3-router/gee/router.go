package gee

import (
	"log"
	"net/http"
)

//type router struct {
//	handlers map[string]HandlerFunc
//}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
