# Running

For running linters stages and application on containers, run: 
```sh
$ make -f ./deploy/main.Makefile deploy
```

For running coverage util command, run:
```sh
$ make -f ./deploy/main.Makefile coverage
```

  - Type some Markdown on the left
  - See HTML in the right
  - Magic

# Testing

 - Dump endpoint:
```sh
curl -X POST   http://localhost:8077/v1/dump   -H 'cache-control: no-cache'   -H 'content-type: application/json' -d '{
"file_path": "~/tmp.zts"
}' > ~/tmp.zts
```


Unlike the test task, the tree is a structure: `key -> int, val -> interface{}`.

Loadind of init data you can do from snapshot file or from network endpoint `load`.

# P.S.
I don't have time to describe all the endpoints, but everything is clear by the code and especially by the tests in the file `handlers_test.go`.