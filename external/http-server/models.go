package http_server

type InsertRequest struct {
	Key   int         `json:"key"`
	Value interface{} `json:"value"`
}

type DumpRequest struct {
	FilePath string `json:"file_path"`
}
