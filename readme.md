## yyfdocker

A simple implementation of docker

### Run

```
cd yyfdocker
go build .
./yyfdocker run -it /bin/bash
```

### Log file
In YYFdocker.log file

### Problem
If we don't add the code "syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")"
in `init.go`, we cannot see /proc directory after we exit out namespace. (https://github.com/xianlubird/mydocker/issues/41)

### Idea
"yyfdocker run -it /bin/bash" will first call `Run(true, "/bin/bash")` in main.go.
Inside `Run`, it will call `container.NewParentProcess(true, "/bin/bash")`.
`container.NewParentProcess(true, "/bin/bash")` will help create a isolated namespace (a new process in an isolated namespace) which is so-called container. The container's first process is
an `init` process (NewParentProcess -> call yyfdocker init /bin/bash). After that, in `init.go`, 
`RunContainerInitProcess` will mount "/proc" in the namespace of container and replace `init` process with `/bin/bash` process using `syscall.Exec`.