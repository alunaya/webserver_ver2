package cws

import (
	"fmt"
	"html/template"
	"net/http"
)

type Router struct {
	routeMap map[string]map[string] http.HandlerFunc
	templates *template.Template
}
type Controller func(c *Context)
type routeMap map[string]map[string] http.HandlerFunc

func NewRouter() *Router{
	return &Router{
		routeMap: make(routeMap),
	}
}

func (r *Router)LoadHtmlGlob(glob string){
	r.templates = template.Must(template.ParseGlob(glob))
}

func registerRoute( routePath string, method string, handler http.HandlerFunc, routeMap routeMap){
	route := routeMap[routePath]

	if route == nil {
		route = make(map[string]http.HandlerFunc)
	}

	if method == MethodAny {
		for k, _ := range route {
			if route[k] != nil{
				panic(fmt.Sprintf("function for method \"%v\" of route \"%v\" are already registered", MethodAny, routePath))
			}
		}


	}

	if route[method] != nil {
		panic(fmt.Sprintf("function for method \"%v\" of route \"%v\" are already registered", method, routePath))
	}

	route[method] = handler
}

func(router *Router) Get(routePath string, controller Controller){
	registerRoute( routePath, http.MethodGet, controllerBuilder(router.templates, controller), router.routeMap)
}

func(router *Router) Post(routePath string, controller Controller){
	registerRoute( routePath, http.MethodPost, controllerBuilder(router.templates, controller), router.routeMap)
}

func(router *Router) Patch(routePath string, controller Controller){
	registerRoute( routePath, http.MethodPatch, controllerBuilder(router.templates, controller), router.routeMap)
}

func(router *Router) Put(routePath string, controller Controller){
	registerRoute( routePath, http.MethodPut, controllerBuilder(router.templates, controller), router.routeMap)
}

func(router *Router) Delete(routePath string, controller Controller){
	registerRoute( routePath, http.MethodDelete, controllerBuilder(router.templates, controller), router.routeMap)
}

func(router *Router) Any(routePath string, controller Controller){
	registerRoute( routePath, MethodAny, controllerBuilder(router.templates, controller), router.routeMap)
}

func controllerBuilder(templates *template.Template, controller Controller) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		var context = NewContext(w, r, templates)
		controller(context)
	}
}