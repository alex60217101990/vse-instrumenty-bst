package bst

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/alex60217101990/vse-instrumenty-bst/external/helpers"
	"github.com/alex60217101990/vse-instrumenty-bst/external/logger"

	"github.com/stretchr/testify/assert"
)

var bst BinarySearchTree

func init() {
	helpers.InitConfigs("")
	logger.InitLoggerSettings()

	// configs.Conf.IsDebug = true
}

func fillTree(bst *BinarySearchTree) {
	bst.Insert(8, "8")
	bst.Insert(4, "4")
	bst.Insert(10, "10")
	bst.Insert(2, "2")
	bst.Insert(6, "6")
	bst.Insert(1, "1")
	bst.Insert(3, "3")
	bst.Insert(5, "5")
	bst.Insert(7, 15)
	bst.Insert(8, 28)
}

func TestMarshal(t *testing.T) {
	defer logger.CloseLoggers()

	fillTree(&bst)

	bts, err := json.Marshal(&bst)

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(bts))
}

func TestUnmarshal(t *testing.T) {
	defer logger.CloseLoggers()

	unmarshalTest := func() error {
		var bts BinarySearchTree

		return json.Unmarshal([]byte(`[
		"foo",
		100
	  ]`), &bts)
	}

	var tests = []*testUnmarshal{
		&testUnmarshal{
			name: "success",
			testData: []byte(`{
		"1": "foo",
		"2": 100
	  }`),
			body: unmarshalTest,
			err:  "",
		},
		&testUnmarshal{
			name: "fail",
			testData: []byte(`[
			"foo",
			100
		  ]`),
			body: unmarshalTest,
			err:  "cannot unmarshal array into Go value of type map[int]interface {}",
		},
	}

	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			err := st.body()
			if err != nil {
				if !assert.Contains(t, err.Error(), st.err) {
					t.Error(err)
					return
				}
			}

			t.Logf("test: %s was success\n", st.name)
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	defer logger.CloseLoggers()

	loadTest := func(td *testDataItem) (err error) {
		var (
			bst BinarySearchTree
			f   *os.File
		)

		f, err = os.OpenFile(td.fileName, os.O_RDWR, 644)
		defer f.Close()
		if err != nil {
			return err
		}

		if td.useCompress {
			err = bst.Load(f, true)
		} else {
			err = bst.Load(f)
		}

		bst.String()

		return err
	}

	var tests = []*testLoadDump{
		&testLoadDump{
			name: "success load with compress",
			testData: &testDataItem{
				fmt.Sprintf("..%s..%stmp%stest-dump-data.zst", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator)),
				true,
			},
			body: loadTest,
			err:  "",
		},
		&testLoadDump{
			name: "success load without compress",
			testData: &testDataItem{
				fmt.Sprintf("..%s..%stmp%stest-data.json", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator)),
				false,
			},
			body: loadTest,
			err:  "",
		},
	}

	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			err := st.body(st.testData)
			if err != nil {
				if !assert.Contains(t, err.Error(), st.err) {
					t.Error(err)
					return
				}
			}

			t.Logf("test: %s was success\n", st.name)
		})
	}
}

func TestDumpToFile(t *testing.T) {
	defer logger.CloseLoggers()

	dumpTest := func(td *testDataItem) (err error) {
		var (
			bst BinarySearchTree
			f   *os.File
		)

		fillTree(&bst)

		f, err = os.OpenFile(td.fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		defer f.Close()
		if err != nil {
			return err
		}

		if td.useCompress {
			return bst.Dump(f, true)
		}

		return bst.Dump(f)
	}

	var tests = []*testLoadDump{
		&testLoadDump{
			name: "success dump with compress",
			testData: &testDataItem{
				fmt.Sprintf("..%s..%stmp%stest-data.zst", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator)),
				true,
			},
			body: dumpTest,
			err:  "",
		},
		&testLoadDump{
			name: "success dump without compress",
			testData: &testDataItem{
				fmt.Sprintf("..%s..%stmp%stest-data.json", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator)),
				false,
			},
			body: dumpTest,
			err:  "",
		},
	}

	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			err := st.body(st.testData)
			if err != nil {
				if !assert.Contains(t, err.Error(), st.err) {
					t.Error(err)
					return
				}
			}

			t.Logf("test: %s was success\n", st.name)
		})
	}
}

func TestSearch(t *testing.T) {
	fillTree(&bst)

	var tests = []*testDeleteKey{
		&testDeleteKey{
			&baseTest{
				name: "success",
			}, 2,
		},
		&testDeleteKey{
			&baseTest{
				name: "failed",
			}, 15,
		},
	}

	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			_, ok := bst.Search(st.testKey)
			switch st.name {
			case "failed":
				assert.Equal(t, false, ok)
			case "success":
				assert.Equal(t, true, ok)
			}

			t.Logf("test: %s was success\n", st.name)
		})
	}
}

func TestDeleteValue(t *testing.T) {
	defer logger.CloseLoggers()

	fillTree(&bst)

	var tests = []*testDeleteKey{
		&testDeleteKey{
			&baseTest{
				name: "success",
				err:  "",
			}, 2,
		},
		&testDeleteKey{
			&baseTest{
				name: "failed",
				err:  DelNotExistKey.Error(),
			}, 15,
		},
	}

	for _, st := range tests {
		t.Run(st.name, func(t *testing.T) {
			err := bst.Delete(st.testKey)
			if err != nil {
				if !assert.Equal(t, err.Error(), st.err) {
					t.Error(err)
					return
				}
			}

			t.Logf("test: %s was success\n", st.name)
		})
	}
}
