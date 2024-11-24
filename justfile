run *ARGS: 
  go build -o /tmp/codecrafters-build-redis-go app/*.go && /tmp/codecrafters-build-redis-go {{ ARGS }}

clean:
  rm -f /tmp/codecrafters-build-redis-go

status:
  git status

test:
  just add .
  just commit "Testing commit"
  just push

add PATH:
  git add {{PATH}}

commit MESSAGE:
  git commit --allow-empty -m "{{ MESSAGE }}"

push BRANCH='master':
  git push origin {{ BRANCH }}
  git push github {{ BRANCH }}
