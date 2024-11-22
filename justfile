run: 
  go build -o /tmp/codecrafters-build-redis-go app/*.go && /tmp/codecrafters-build-redis-go "$@"

clean:
  rm -f /tmp/codecrafters-build-redis-go

status:
  git status

add PATH:
  git add {{PATH}}

commit MESSAGE:
  git commit -m "{{ MESSAGE }}"

push BRANCH='master':
  git push origin {{ BRANCH }}
  git push github {{ BRANCH }}
