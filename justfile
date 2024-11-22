run: 
  go run .

status:
  git status

add PATH:
  git add {{PATH}}

commit MESSAGE:
  git commit -m "{{ MESSAGE }}"

push BRANCH='master':
  git push origin {{ BRANCH }}
  git push github {{ BRANCH }}
