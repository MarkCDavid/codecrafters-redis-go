BUILD_PATH := "/tmp/codecrafters-build-redis-go"

build:
  go build -o {{ BUILD_PATH }} app/*.go

run:
  just build
  {{ BUILD_PATH }}



commit MESSAGE:
  git add .
  git commit --allow-empty -m "{{ MESSAGE }}"

push:
  git push codecrafters master
  git push github master

test:
  just commit "Testing"
  just push
