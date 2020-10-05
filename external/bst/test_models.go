package bst

type testDataItem struct {
	fileName    string
	useCompress bool
}

type baseTest struct {
	name string
	err  string
}

type testLoadDump struct {
	name     string
	testData *testDataItem
	body     func(td *testDataItem) error
	err      string
}

type testUnmarshal struct {
	name     string
	testData []byte
	body     func() error
	err      string
}

type testDeleteKey struct {
	*baseTest
	testKey int
}
