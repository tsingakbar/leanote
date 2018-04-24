# tsingakbar's fork of Leanote

## Golang server dev steps

1. set up golang dev env(like GOROOT GOPATH), mongodb and revel like the offical guide described
2. `HTTP_PROXY=http://if.u.r.behind.http.proxy:8080 go get -v github.com/tsingakbar/leanote/app`
3. `revel run github.com/tsingakbar/leanote` to run in dev mode
4. `revel package github.com/tsingakbar/leanote prod` to release a package
5. upload the package to your server and `run.sh`


## Operation and maintenance

* `conf` and `files` folder under `revel.BasePath`(`$GOPATH/src/github.com/tsingakbar/leanote` whose GOPATH is pacakge's root folder) need to be configured or regularly backuped.
* For server having extreamly small memory, the packaging step should disable cgo to forbidden spawning thread by `pthread_create`, which is `CGO_ENABLED=0 revel package github.com/tsingakbar/leanote prod`. But this way usually requires packaging user clone an writable go setup such as `$HOME/goroot` to execute `GOROOT=$HOME/goroot CGO_ENABLED=0 revel package github.com/tsingakbar/leanote prod` because of rebuilding/installing `$GOROOT/pkg/linux_amd64/os/user.a` during packaing process.

> why avoid using pthread_create on openVZ: 
`pthread_create`, `fork` are all implement by syscall `clone()` on linux, which coresponse to kernel LWP. `pthread_create` by default use several MiB as thread stack(`man pthread_create`), and if cgo disabled, golang will use its own thread implementation which calls `clone()` directly with a small thread stack(but of course this thread will not able to run c code anymore). Anyway, several MiB is only VM usage, it won't cause a problem even on machines with little RAM, but openVZ fucked up memory accounting for hosted VPS, which is our concerns.
