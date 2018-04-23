# tsingakbar's fork of Leanote

## Dev steps

1. set up golang dev env(like GOROOT GOPATH), mongodb and revel like the offical guide described
2. `HTTP_PROXY=http://if.u.r.behind.http.proxy:8080 go get -v github.com/tsingakbar/leanote/app`
3. `revel run github.com/tsingakbar/leanote` to run in dev mode
4. `revel package github.com/tsingakbar/leanote prod` to release a package
5. upload the package to your server and `run.sh`
