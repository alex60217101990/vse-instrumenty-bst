package fast_http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"testing"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"
	"github.com/alex60217101990/vse-instrumenty-bst/external/helpers"
	server "github.com/alex60217101990/vse-instrumenty-bst/external/http-server"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

var (
	testServer = NewFastHttpServer()
)

func init() {
	helpers.InitConfigs("")
	logger.InitLoggerSettings()
	configs.Conf.BST.UseCompression = false
	configs.Conf.BST.SnapshotPath = "../../../tmp/test-data.json"
	err := testServer.initFromFile()
	if err != nil {
		logger.DefaultLogger.Fatal(err)
	}
	// configs.Conf.IsDebug = true
}

func serveHttp(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return client.Do(req)
}

func serve(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return client.Do(req, res)
}

func TestInsert(t *testing.T) {
	bts, err := json.Marshal(&server.InsertRequest{
		Key:   3,
		Value: "djfhrkghrkg",
	})
	if err != nil {
		t.Error(err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://localhost:8079/v1/insert")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(bts)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err = serve(testServer.Insert, req, resp)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode() == fasthttp.StatusOK {
		t.Log("success query")
	}
}

func TestSearch(t *testing.T) {
	type tmpData struct {
		Data interface{} `json:"data"`
	}

	testServer.bst.Insert(5, 5)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://localhost:8079/v1/search")
	req.Header.SetMethod("GET")
	req.URI().QueryArgs().Set("key", "5")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := serve(testServer.Get, req, resp)
	if err != nil {
		t.Error(err)
		return
	}

	var tmp tmpData
	err = json.Unmarshal(resp.Body(), &tmp)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, tmp.Data, "5")
}

func TestDelete(t *testing.T) {
	testServer.bst.Insert(5, 5)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://localhost:8079/v1/delete")
	req.Header.SetMethod("DELETE")
	req.URI().QueryArgs().Set("key", "5")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := serve(testServer.Delete, req, resp)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode() == fasthttp.StatusOK {
		t.Log("success query")
	}
}

func TestLoad(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("snapshot", "test.json")
	if err != nil {
		t.Error(err)
		return
	}

	part.Write([]byte(`{
		"15": "1",
		"17": "10",
		"30": "2"
	}`))

	err = writer.Close()
	if err != nil {
		t.Error(err)
		return
	}

	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest("POST", "http://localhost:8079/v1/load", body)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = serveHttp(testServer.Load, req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Log("success test")
	}

	testServer.bst.String()
}

func TestDump(t *testing.T) {
	testServer.bst.Insert(2, 2)
	testServer.bst.Insert(7, 7)

	m, _ := json.MarshalIndent(testServer.bst, "", "\t")

	bts, err := json.Marshal(&server.DumpRequest{
		FilePath: "/tmp/test.json",
	})
	if err != nil {
		t.Error(err)
		return
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("http://localhost:8079/v1/dump")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBody(bts)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = serve(testServer.Dump, req, resp)
	if err != nil {
		t.Error(err)
	}

	if assert.NotEqual(t, string(m), string(resp.Body())) {
		t.Log("success test result")
	}
}
