package fast_http

import (
	"fmt"

	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func (s *FastHttpServer) CorsMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func (s *FastHttpServer) PanicMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					logger.AppLogger.Errorw(fmt.Sprintf("panic middleware detect, %v", err), "fast-http-server", "PanicMiddleware")
					errorPrint(ctx, err, fasthttp.StatusInternalServerError)
				}
			}
			ctx.Done()
		}()
		next(ctx)
	}
}
