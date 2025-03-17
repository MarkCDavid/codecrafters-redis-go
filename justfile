BUILD_PATH := "/tmp/codecrafters-build-redis-go"
GIT := if path_exists(home_directory() + "/.ssh/id_mcd_ed25519") == "true" { "GIT_SSH_COMMAND='ssh -i " + home_directory() + "/.ssh/id_mcd_ed25519 -o IdentitiesOnly=yes' git" } else { "git" }

build:
  go build -o {{ BUILD_PATH }} app/*.go

run:
  just build
  {{ BUILD_PATH }}

commit MESSAGE:
  {{ GIT }} add .
  {{ GIT }} commit --allow-empty -m "{{ MESSAGE }}"

push:
  {{ GIT }} push codecrafters master
  {{ GIT }} push github master

test:
  just commit "Testing"
  just push
