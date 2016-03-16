#!/bin/bash

APP_NAME=find-yourself
TMP_DIR=$PWD/tmp
IMAGE_NAME=local/$APP_NAME

build_on_docker() {
  docker build --tag $IMAGE_NAME:built .
}

build_on_local() {
  env GOOS=linux go build -a -tags netgo -installsuffix cgo -ldflags '-w' -o "${TMP_DIR}/vulcand" .
}

copy() {
  cp $TMP_DIR/$APP_NAME .
  cp $TMP_DIR/$APP_NAME entrypoint/
}

init() {
  rm -rf $TMP_DIR/ \
   && mkdir -p $TMP_DIR/
}

package() {
  docker build --tag $IMAGE_NAME:latest entrypoint
}

run() {
  docker run --rm \
    --volume $TMP_DIR:/export/ \
    $IMAGE_NAME:built \
      cp $APP_NAME /export
}

panic() {
  local message=$1
  echo $message
  exit 1
}

docker_build() {
  init    || panic "init failed"
  build_on_docker || panic "build_on_docker failed"
  run     || panic "run failed"
  copy    || panic "copy failed"
  package || panic "package failed"
}

local_build() {
  init    || panic "init failed"
  build_on_local || panic "build_on_local failed"
  copy    || panic "copy failed"
  package || panic "package failed"
}

main() {
  local mode="$1"
  if [ "$mode" == "local" ]; then
    echo "Local Build"
    local_build
  else
    echo "Docker Build"
    docker_build
  fi
  exit $?
}
main
