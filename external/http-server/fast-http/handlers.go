package fast_http

import (
	"fmt"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	server "github.com/alex60217101990/vse-instrumenty-bst/external/http-server"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"

	"github.com/valyala/fasthttp"
)

func (s *FastHttpServer) Delete(ctx *fasthttp.RequestCtx) {
	key, err := parseKeyQuery(ctx)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "delete")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.bst.Delete(key)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "delete")
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
}

func (s *FastHttpServer) Get(ctx *fasthttp.RequestCtx) {
	key, err := parseKeyQuery(ctx)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "get")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	value, ok := s.bst.Search(key)
	if !ok {
		logger.AppLogger.Errorw(server.ValueNotFound.Error(), "fast-http-server", "get")
		errorPrint(ctx, server.ValueNotFound, fasthttp.StatusInternalServerError)
		return
	}

	messagePrint(ctx, value, fasthttp.StatusOK)
}

func (s *FastHttpServer) Insert(ctx *fasthttp.RequestCtx) {
	request, err := parseInsertRequest(ctx)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "insert")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	err = s.bst.Insert(request.Key, request.Value)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "insert")
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
}

func (s *FastHttpServer) Load(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("snapshot")
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "load")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	f, err := header.Open()
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "load")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}
	defer func() {
		f.Close()
		ctx.Request.RemoveMultipartFormFiles()
	}()

	if configs.Conf.BST.UseCompression {
		err = s.bst.Load(f, true)
	} else {
		err = s.bst.Load(f)
	}

	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "load")
		errorPrint(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
}

func (s *FastHttpServer) Dump(ctx *fasthttp.RequestCtx) {
	dp, err := parseDumpRequest(ctx)
	if err != nil {
		logger.AppLogger.Errorw(err.Error(), "fast-http-server", "dump")
		errorPrint(ctx, err, fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("application/octet-stream")
	ctx.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, dp.FilePath))
	// w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))

	rwc := &RequestWriteCloser{ctx}
	if configs.Conf.BST.UseCompression {
		err = s.bst.Dump(rwc, true)
	} else {
		err = s.bst.Dump(rwc)
	}
}
