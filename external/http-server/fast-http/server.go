package fast_http

import (
	"fmt"
	"os"
	"time"

	"github.com/alex60217101990/vse-instrumenty-bst/external/bst"
	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type FastHttpServer struct {
	server *fasthttp.Server
	bst    bst.Tree
}

func NewFastHttpServer() *FastHttpServer {
	return &FastHttpServer{
		bst: &bst.BinarySearchTree{},
	}
}

func (s *FastHttpServer) routing() *router.Router {
	r := router.New()

	gr := r.Group("/v1")
	gr.DELETE("/delete", s.Delete)
	gr.GET("/search", s.Get)
	gr.POST("/insert", s.Insert)
	gr.POST("/load", s.Load)
	gr.POST("/dump", s.Dump)

	return r
}

func (s *FastHttpServer) initFromFile() (err error) {
	if len(configs.Conf.BST.SnapshotPath) > 0 {
		var f *os.File
		f, err = os.OpenFile(configs.Conf.BST.SnapshotPath, os.O_RDWR, 644)
		defer f.Close()
		if err != nil {
			logger.AppLogger.Fatalw(err.Error(), "fast-http-server", "init")
		}

		if configs.Conf.BST.UseCompression {
			err = s.bst.Load(f, true)
		} else {
			err = s.bst.Load(f)
		}
	}
	s.bst.String()
	return err
}

func (s *FastHttpServer) Init() {
	err := s.initFromFile()
	if err != nil {
		logger.AppLogger.Fatalw(err.Error(), "fast-http-server", "Init")
	}

	s.server = &fasthttp.Server{
		Name:               configs.Conf.ServiceName,
		Concurrency:        100000,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 3 * time.Second,
		ReadBufferSize:     1 << 10,
		WriteBufferSize:    1 << 10,
		ReadTimeout:        7 * time.Second,
		WriteTimeout:       7 * time.Second,
		IdleTimeout:        15 * time.Second,
		MaxRequestBodySize: 4 * 1024 * 1024 * 1024,
		//Logger:             logger.AppLogger,
		Handler: fasthttp.CompressHandler(s.PanicMiddleware(s.CorsMiddleware(s.routing().Handler))),
	}

	if configs.Conf.IsDebug {
		s.server.LogAllErrors = true
	}
}

func (s *FastHttpServer) Run() {
	logger.CmdInfo.Println("🔥 FastHttp server started.")
	logger.CmdError.Println(s.server.ListenAndServe(
		fmt.Sprintf("%s:%d", configs.Conf.Server.Host, configs.Conf.Server.Port)))
}

func (s *FastHttpServer) Close() error {
	defer logger.CmdInfo.Println("🔥 FastHttp server stoped.")
	return s.server.Shutdown()
}
