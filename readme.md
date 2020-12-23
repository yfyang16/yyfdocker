## yyfdocker

A very simple implementation of docker

### Usage

```
cd yyfdocker
go build .

./yyfdocker run [-it/-d] [-m {v}] [-cpushare {v}] [-cpuset {v}] [-v {host path}:{container path}] [--name {containerName}] {imageName} {commands}
./yyfdocker commit {containerName} {imageName}
./yyfdocker ps
./yyfdocker logs {containerName}
./yyfdocker stop {containerName}
./yyfdocker rm {containerName}
./yyfdcoker exec {containerName} {commands}
```

### Log file
In YYFdocker.log file

### Problem
If we don't add the code "syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")"
in `init.go`, we cannot see /proc directory after we exit out namespace. (https://github.com/xianlubird/mydocker/issues/41)

When writing pid value into "/sys/fs/cgroup/cpuset/yyfdocker-cgroup/tasks", it would return no space left error. In order to solve
this problem, we should write "0" in "/sys/fs/cgroup/cpuset/yyfdocker-cgroup/cpuset.mems" and make sure "/sys/fs/cgroup/cpuset/yyfdocker-cgroup/cpuset.cpus" is not empty.

### Idea
"yyfdocker run -it /bin/bash" will first call `Run(true, "/bin/bash")` in main.go.
Inside `Run`, it will call `container.NewParentProcess(true, "/bin/bash")`.
`container.NewParentProcess(true, "/bin/bash")` will help create a isolated namespace (a new process in an isolated namespace) which is so-called container. The container's first process is
an `init` process (NewParentProcess -> call yyfdocker init /bin/bash). After that, in `init.go`, 
`RunContainerInitProcess` will mount "/proc" in the namespace of container and replace `init` process with `/bin/bash` process using `syscall.Exec`.