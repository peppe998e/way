# Way
HTTP router for GoLang

* Deliberately simple
* Extremely fast
* Route based on HTTP methods and path
* Path parameters via `Context` (e.g. `/music/:band/:song`)
* Trailing `/` matches path prefixes

## Install

There's no need to add a dependency to Way, just copy `way.go` and `way_test.go` into your project, or [drop](https://github.com/matryer/drop) them in:

```
drop github.com/peppe998e/way
```

If you prefer, it is go gettable:

```
go get github.com/peppe998e/way
```

## Usage

* Use `NewRouter` to make a new `Router`
* Call `Handle`, `ALL`, `GET`, `POST`... to add handlers
* Specify HTTP method and path pattern for each route
* Use `Param` function to get the path parameters from the context

```go
func main() {
	router := way.NewRouter()

	router.GET("/music/:song", handleReadSong)
	router.Handle(way.WAY_GET, "/music/:band/:song", handleReadSong)
	
	router.Handle(way.WAY_GET|way.WAY_POST, "/author/:name", handleUpdateSong)
	router.DELETE("/author/:name", handleDeleteSong)

	log.Fatalln(http.ListenAndServe(":8080", router))
}

func handleReadSong(w http.ResponseWriter, r *http.Request) {
	band := way.Param(r.Context(), "band")
	song := way.Param(r.Context(), "song")
	// use 'band' and 'song' parameters...
}
```

* Prefix matching

To match any path that has a specific prefix, use the `...` prefix indicator:

```go
func main() {
	router := way.NewRouter()

	router.GET("/images...", handleImages)
	log.Fatalln(http.ListenAndServe(":8080", router))
}
```

In the above example, the following paths will match:

* `/images`
* `/images/`
* `/images/one/two/three.jpg`

* Set `Router.NotFound` to handle 404 errors manually

```go
func main() {
	router := way.NewRouter()
	
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not the page you are looking for")
	})
	
	log.Fatalln(http.ListenAndServe(":8080", router))
}
```

## Why another HTTP router?

I know, I know. But no routers offer the simplicity of path parameters via Context, and HTTP method matching. Which covers 100% of my use cases so far.
