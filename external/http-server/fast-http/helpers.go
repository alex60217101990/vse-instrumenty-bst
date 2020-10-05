package fast_http

import (
	"encoding/json"
	"strconv"

	server "github.com/alex60217101990/vse-instrumenty-bst/external/http-server"

	"github.com/valyala/fasthttp"
)

func messagePrint(ctx *fasthttp.RequestCtx, data interface{}, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	encoder := json.NewEncoder(ctx)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func errorPrint(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes([]byte("application/json"))
	if err1 := json.NewEncoder(ctx).Encode(map[string]string{
		"error": err.Error(),
	}); err1 != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func parsePostBody(ctx *fasthttp.RequestCtx, data interface{}) error {
	// Get the JSON body and decode into credentials
	return json.Unmarshal(ctx.PostBody(), data)
}

func parseKeyQuery(ctx *fasthttp.RequestCtx) (k int, err error) {
	key := string(ctx.QueryArgs().Peek("key"))
	if len(string(key)) > 0 {
		var i int64
		i, err = strconv.ParseInt(string(key), 10, 16)
		if err != nil {
			return k, err
		}
		return int(i), err
	}
	return k, server.EmptyKey
}

func parseInsertRequest(ctx *fasthttp.RequestCtx) (prod *server.InsertRequest, err error) {
	prod = &server.InsertRequest{}
	err = parsePostBody(ctx, prod)
	return prod, err
}

func parseDumpRequest(ctx *fasthttp.RequestCtx) (d *server.DumpRequest, err error) {
	d = &server.DumpRequest{}
	err = parsePostBody(ctx, d)
	return d, err
}
