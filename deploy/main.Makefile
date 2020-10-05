APP_NAME=vse-instrumenty-bst
GIT_HASH ?= $(shell git rev-parse --short HEAD)
PKG_NAME=github.com/alex60217101990/${APP_NAME}
VERSION ?= 0.1
BUILD_OS ?= linux
BUILD_PATCH ?= develop
BUILD_VERSION ?= ${VERSION}.${BUILD_PATCH}

TMP := ''
deploy_dir := ''
tmp_procedure:
ifeq ($(notdir $(patsubst %/,%,$(CURDIR))),vse-instrumenty-bst)
	$(eval CUR_DIR=$(shell pwd | rev | cut -d'/' -f1- | rev))
	$(eval TMP=$(shell find ${CUR_DIR} -name "tmp" -print))
	$(eval deploy_dir=$(shell find ${CUR_DIR} -name "deploy" -print))
else
	$(eval CUR_DIR=$(shell pwd | rev | cut -d'/' -f2- | rev))
	$(eval TMP=$(shell find ${CUR_DIR} -name "tmp" -print))
	$(eval deploy_dir=$(shell find ${CUR_DIR} -name "deploy" -print))
endif

test: tmp_procedure
	@echo "Running unit tests with coverage..."
ifeq ($(notdir $(patsubst %/,%,$(CURDIR))),vse-instrumenty-bst)
	go test -v -cover -coverprofile=${TMP}/${APP_NAME}.coverprofile -coverpkg="${PKG_NAME}/..." ./...
else
	cd .. && $(MAKE) -f ./deploy/main.Makefile test
endif

cover_test: test
	go tool cover -html=${TMP}/${APP_NAME}.coverprofile -o ${TMP}/${APP_NAME}.coverprofile.html
	go tool cover -func=${TMP}/${APP_NAME}.coverprofile

coverage: test
	go tool cover -html=${TMP}/${APP_NAME}.coverprofile

rm_cover_profiles: tmp_procedure
	find ${TMP} -type f -name '*.coverprofile' -delete

print-%: tmp_procedure
	@echo '$*=$($*)'

go_lint: tmp_procedure
	@printf "Validate golang code...\n"
	docker build -t alex6021710/go-linter:v0.0.1 -f ${deploy_dir}/dockerfiles/Dockerfile.golint ${CUR_DIR} && \
	docker run --rm -it alex6021710/go-linter:v0.0.1 golangci-lint run /app/cmd 


docker_lint: tmp_procedure
	@printf "Validate dockerfiles...\n"
	for file in $(shell ls ${deploy_dir}/dockerfiles/Dockerfile.*); do \
		docker run --rm -i hadolint/hadolint < $${file} ; \
	done

lint: docker_lint go_lint

go_push: tmp_procedure
	@printf "Build app image and push them...\n"
	docker build -t alex6021710/vse-instrumenty-bst:v0.0.1 -f ${deploy_dir}/dockerfiles/Dockerfile.golang ${CUR_DIR} && \
	docker login && \
	docker push alex6021710/vse-instrumenty-bst:v0.0.1 

run: 
	docker run --rm -p 127.0.0.1:8077:8077 -i alex6021710/vse-instrumenty-bst:v0.0.1

deploy: lint run