package way

import (
	"context"
	"net/http"
	"strings"
)

// wayContextKey is the context key type for storing
// parameters in context.Context.
type wayContextKey string

// Router routes HTTP requests.
type Router interface {
	Handle(method, pattern string, handler http.Handler)
	HandleFunc(method, pattern string, fn http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type router1 struct {
	routes []*route
	// NotFound is the http.Handler to call when no routes
	// match. By default uses http.NotFoundHandler().
	NotFound http.Handler
}

// NewRouter makes a new Router.
func NewRouter() Router {
	return &router1{
		NotFound: http.NotFoundHandler(),
	}
}

func (r *router1) pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

// Handle adds a handler with the specified method and pattern.
// Method can be any HTTP method string or "*" to match all methods.
// Pattern can contain path segments such as: /item/:id which is
// accessible via context.Value("id").
// If pattern ends with trailing /, it acts as a prefix.
func (r *router1) Handle(method, pattern string, handler http.Handler) {
	route := &route{
		method:  strings.ToLower(method),
		segs:    r.pathSegments(pattern),
		handler: handler,
		prefix:  strings.HasSuffix(pattern, "/"),
	}
	r.routes = append(r.routes, route)
}

// HandleFunc is the http.HandlerFunc alternative to http.Handle.
func (r *router1) HandleFunc(method, pattern string, fn http.HandlerFunc) {
	r.Handle(method, pattern, fn)
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (r *router1) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := strings.ToLower(req.Method)
	segs := r.pathSegments(req.URL.Path)
	for _, route := range r.routes {
		if route.method != method && route.method != "*" {
			continue
		}
		if ctx, ok := route.match(req.Context(), r, segs); ok {
			route.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.NotFound.ServeHTTP(w, req)
}

// Param gets the path parameter from the specified Context.
// Returns an empty string if the parameter was not found.
func Param(ctx context.Context, param string) string {
	v := ctx.Value(wayContextKey(param))
	if v == nil {
		return ""
	}
	vStr, ok := v.(string)
	if !ok {
		return ""
	}
	return vStr
}

type route struct {
	method  string
	segs    []string
	handler http.Handler
	prefix  bool
}

func (r *route) match(ctx context.Context, router *router1, segs []string) (context.Context, bool) {
	if len(segs) > len(r.segs) && !r.prefix {
		return nil, false
	}
	for i, seg := range r.segs {
		if i > len(segs)-1 {
			return nil, false
		}
		isParam := false
		if strings.HasPrefix(seg, ":") {
			isParam = true
			seg = strings.TrimPrefix(seg, ":")
		}
		if !isParam { // verbatim check
			if seg != segs[i] {
				return nil, false
			}
		}
		if isParam {
			ctx = context.WithValue(ctx, wayContextKey(seg), segs[i])
		}
	}
	return ctx, true
}
