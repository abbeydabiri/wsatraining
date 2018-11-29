package api

import (
	"compress/gzip"
	"context"
	"io"
	"path"
	"strings"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"

	"wsatraining/config"
)

//Message ...
type Message struct {
	Code    int
	Message string
	Body    interface{}
}

//StartRouter ...
func StartRouter() {

	// total := 0
	// if config.Get().Sqlite3.Get(&total, "select count(id) from logs"); total == 0 {
	// 	utils.SaveFileToPath("adminurl", "config", []byte(config.Get().Adminurl))
	// }

	middlewares := alice.New()
	router := NewRouter()

	router.NotFound = middlewares.ThenFunc(
		func(httpRes http.ResponseWriter, httpReq *http.Request) {
			switch {
			default:
				fileServe(httpRes, httpReq)
			case strings.HasPrefix(httpReq.URL.Path, "/wp/"):
				apiProxy(httpRes, httpReq)
			}
		})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
	}).Handler(router)

	println("Now Serving")
	sMessage := "serving @ " + config.Get().Address
	println(sMessage)
	log.Println(sMessage)
	log.Fatal(http.ListenAndServe(config.Get().Address, handler))
}

func cut(name string) string {
	name = strings.TrimSuffix(name, "/")
	dir, _ := path.Split(name)
	return dir
}

func fileServe(httpRes http.ResponseWriter, httpReq *http.Request) {
	urlPath := strings.Replace(httpReq.URL.Path, "//", "/", -1)
	if strings.HasSuffix(urlPath, "/") {
		urlPath = path.Join(urlPath, "index.html")
	}

	if strings.HasPrefix(urlPath, "/") && len(urlPath) > 1 {
		urlPath = urlPath[1:]
	}

	var err error
	var dataBytes []byte

	if dataBytes, err = config.Asset(urlPath); err != nil {
		for urlPath != "/" {
			log.Printf("urlPath - %s", urlPath)
			urlPath = cut(urlPath)
			newPath := path.Join(urlPath, "index.html")
			if dataBytes, err = config.Asset(newPath); err == nil {
				break
			} else {
				log.Printf("err - %s", err.Error())
			}
		}
	}

	httpRes.Header().Set("Cache-Control", "max-age=0, must-revalidate")
	httpRes.Header().Set("Pragma", "no-cache")
	httpRes.Header().Set("Expires", "0")

	httpRes.Header().Add("Content-Type", config.ContentType(urlPath))
	if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
		httpRes.Write(dataBytes)
		return
	}
	gzipWrite(dataBytes, httpRes)
}

//Router ...
type Router struct { // Router struct would carry the httprouter instance,
	*httprouter.Router //so its methods could be verwritten and replaced with methds with wraphandler
}

//Get ...
func (router *Router) Get(path string, handler http.Handler) {
	router.GET(path, wrapHandler(handler)) // Get is an endpoint to only accept requests of method GET
}

//Post is an endpoint to only accept requests of method POST
func (router *Router) Post(path string, handler http.Handler) {
	router.POST(path, wrapHandler(handler))
}

//Put is an endpoint to only accept requests of method PUT
func (router *Router) Put(path string, handler http.Handler) {
	router.PUT(path, wrapHandler(handler))
}

//Delete is an endpoint to only accept requests of method DELETE
func (router *Router) Delete(path string, handler http.Handler) {
	router.DELETE(path, wrapHandler(handler))
}

//NewRouter is a wrapper that makes the httprouter struct a child of the router struct
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipWrite(dataBytes []byte, httpRes http.ResponseWriter) {
	httpRes.Header().Set("Content-Encoding", "gzip")
	gzipHandler := gzip.NewWriter(httpRes)
	defer gzipHandler.Close()
	httpResGzip := gzipResponseWriter{Writer: gzipHandler, ResponseWriter: httpRes}
	httpResGzip.Write(dataBytes)
}

func wrapHandler(httpHandler http.Handler) httprouter.Handle {
	return func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		ctx := context.WithValue(httpReq.Context(), "params", httpParams)
		httpReq = httpReq.WithContext(ctx)

		if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
			httpHandler.ServeHTTP(httpRes, httpReq)
			return
		}

		httpRes.Header().Set("Content-Encoding", "gzip")
		gzipHandler := gzip.NewWriter(httpRes)
		defer gzipHandler.Close()
		httpResGzip := gzipResponseWriter{Writer: gzipHandler, ResponseWriter: httpRes}
		httpHandler.ServeHTTP(httpResGzip, httpReq)
	}
}
