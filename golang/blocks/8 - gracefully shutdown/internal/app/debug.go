package app

import (
	"github.com/dimfeld/httptreemux/v5"
	_ "github.com/dimfeld/httptreemux/v5"
	"github.com/fkaanoz/gracefully-shutdown/internal/controllers"
	"net/http"
	"net/http/pprof"
)

func DebugServer() *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	// standard library
	mux.Handle(http.MethodGet, "/debug/", pprof.Index)
	mux.Handle(http.MethodGet, "/debug/cmdline", pprof.Cmdline)
	mux.Handle(http.MethodGet, "/debug/symbol", pprof.Symbol)
	mux.Handle(http.MethodGet, "/debug/trace", pprof.Trace)

	// custom end points
	mux.Handle(http.MethodGet, "/debug/test", controllers.TestDebug)

	return mux
}
