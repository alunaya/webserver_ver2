package cws

import (
	"fmt"
	"net/http"
)

type Server struct {
	port uint
	middlewarePreRegistry []http.HandlerFunc
	middlewarePostRegistry []http.HandlerFunc
	Router Router
}

func NewServer(port uint) *Server{
	s := Server{
		port: port,
		middlewarePreRegistry: make([]http.HandlerFunc, 0),
		middlewarePostRegistry: make([]http.HandlerFunc, 0),
	}
	return &s
}

func (s *Server) Run(){
	wrappedMux := http.NewServeMux()
	trueMux := http.NewServeMux()
	var preMiddlewareChain http.Handler
	var postMiddlewareChain http.Handler

	trueMux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		w.Write([]byte(r.Method))
	})

	if len(s.middlewarePreRegistry)>0 {
		for i := len(s.middlewarePreRegistry) - 1; i >= 0; i-- {
			if i == len(s.middlewarePreRegistry)-1 {
				preMiddlewareChain = s.middlewarePreRegistry[i]
				continue
			}

			preMiddlewareChain = preMiddlewareBuilder(preMiddlewareChain, s.middlewarePreRegistry[i])
		}
	}

	if len(s.middlewarePostRegistry)>0 {
		for i,v := range s.middlewarePostRegistry {
			if i == 0 {
				postMiddlewareChain = v
				continue
			}

			postMiddlewareChain = postMiddlewareBuilder(postMiddlewareChain, s.middlewarePostRegistry[i])
		}
	}

	rootHandlerFunc := func (w http.ResponseWriter, r *http.Request) {
		if preMiddlewareChain != nil {
			preMiddlewareChain.ServeHTTP(w, r)
		}

		trueMux.ServeHTTP(w,r)

		if postMiddlewareChain != nil {
			postMiddlewareChain.ServeHTTP(w, r)
		}
	}

	wrappedMux.HandleFunc("/", rootHandlerFunc)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", s.port),
		Handler: wrappedMux,
	}

	fmt.Printf("Listening on port %v\n", s.port)
	server.ListenAndServe()
}

func preMiddlewareBuilder (next http.Handler, middleware http.HandlerFunc) http.Handler {
	fn := func (w http.ResponseWriter, r *http.Request){
		middleware(w,r)
		next.ServeHTTP(w,r)
	}
	return http.HandlerFunc(fn)
}

func postMiddlewareBuilder (previous http.Handler, middleware http.HandlerFunc) http.Handler {
	fn := func (w http.ResponseWriter, r *http.Request){
		previous.ServeHTTP(w,r)
		middleware(w,r)
	}
	return http.HandlerFunc(fn)
}

func (s *Server) RegisterPreMiddleware (middleware http.HandlerFunc){
	s.middlewarePreRegistry = append (s.middlewarePreRegistry, middleware)
}

func (s *Server) RegisterPostMiddleware (middleware http.HandlerFunc){
	s.middlewarePostRegistry = append(s.middlewarePostRegistry, middleware)
}