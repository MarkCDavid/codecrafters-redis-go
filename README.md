[![progress-banner](https://backend.codecrafters.io/progress/redis/d6beea00-82d3-4411-aca5-2eae9afe7c54)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

In this challenge, you'll build a toy Redis clone that's capable of handling
basic commands like `PING`, `SET` and `GET`. Along the way we'll learn about
event loops, the Redis protocol and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Instructions

The repository contains a justfile to assist with performing the challenge.

Requirements:
1. `go`
1. `git`
1. `just`

`just add`, `just commit` and `just push` provides assistance with `git` commands, and pushing to two different remotes at the same time.

`just run` builds and runs the application. `just cleans` removes the built executable from `/tmp` folder.

# Reference

[Redis serialization protocol (RESP) specification](https://redis.io/docs/latest/develop/reference/protocol-spec)
