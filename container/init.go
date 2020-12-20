package container

import (
    "log"
    "os"
    "os/exec"
    "syscall"
    "io/ioutil"
    "fmt"
    "strings"
    "path"
)

/*
 * RunContainerInitProcess: This function will be executed in a yyfdocker inside.
 * Mount "/proc" to yydocker /proc
 */
func RunContainerInitProcess() error {
    log.Printf("** RunContainerInitProcess START **\n")

    // wait and read user cmds from the read pipe
    cmdArray := ReadUserCmd() 
    if cmdArray == nil {
        log.Panicf("[RunContainerInitProcess] Could not read pipe!")
        return fmt.Errorf("[RunContainerInitProcess] Could not read pipe!")
    }

    err := SetUpMount()
    if err != nil {
        log.Panicf("[RunContainerInitProcess] SetUpMount failed: %v", err)
        return fmt.Errorf("[RunContainerInitProcess] SetUpMount failed: %v", err)
    }

    // find the absolute path of a cmd in $PATH of OS
    // e.g. /bin/bash is the absolute path of bash cmd
    absolutePath, err := exec.LookPath(cmdArray[0])
    if err != nil {
        log.Panicf("[RunContainerInitProcess] Could not find the corresponding path of %v!", cmdArray[0])
        return err
    }

    log.Printf("** RunContainerInitProcess END **\n")

    // syscall.Exec will execute "command" and replace the init process with
    // "command" process. (So the first process (pid == 1) will be "command")
    err = syscall.Exec(absolutePath, cmdArray[0:], os.Environ())
    if err != nil {
        log.Printf(err.Error())
    }

    return err
}


/** Init process receive an extra file which is a read pipe recording the user cmds */
func ReadUserCmd() []string {
    readPipe := os.NewFile(uintptr(3), "readPipe")
    rawCommand, err := ioutil.ReadAll(readPipe)
    if err != nil {
        log.Panicf("[ReadUserCmd] Could not read pipe!")
        return nil
    }
    return strings.Split(string(rawCommand), " ")

}

/** Pivot root ,and mount proc and tmpfs*/
func SetUpMount() error {
    pwd, err := os.Getwd()
    if err != nil {
        log.Panicf("[SetUpMount] Get workin directory failed: %v", err)
        return fmt.Errorf("[SetUpMount] Get workin directory failed: %v", err)
    }

    log.Printf("[SetUpMount] working directory is: %s", pwd)

    // func Mount(source string, target string, fstype string, flags uintptr, data string)
    // We should explicitly state that this mount namespace is independent.
    err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
    if err != nil {
        log.Printf(err.Error())
        return fmt.Errorf(err.Error())
    }

    err = PivotRoot(pwd)
    if err != nil {
        log.Panicf("[SetUpMount] PivotRoot throws failure: %v", err)
        return fmt.Errorf("[SetUpMount] PivotRoot throws failure: %v", err)
    }

    // mount proc device on "/proc" directory
    // flag meaning:
    //   - MS_NOEXEC: no other program can be executed in this file system
    //   - MS_NOSUID: "set-user-ID" and "set-group-ID" are not allowed
    //   - MS_NODEV : default flag
    syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV, "")
    syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
    return nil
}

func PivotRoot(root string) error {
    log.Printf("In Pivotroot!")
    if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
        log.Panicf("[PivotRoot] mount -bind failed: %v", err)
        return fmt.Errorf("[PivotRoot] mount -bind failed: %v", err)
    }

    pivotDir := path.Join(root, ".pivot_root")
    // 创建 rootfs/.pivot_root 存储 old_root
    if err := os.Mkdir(pivotDir, 0777); err != nil {
        log.Panicf("[PivotRoot] make directory failed: %v", err)
        return fmt.Errorf("[PivotRoot] make directory failed: %v", err)
    }
    // pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
    // 挂载点现在依然可以在mount命令中看到
    if err := syscall.PivotRoot(root, pivotDir); err != nil {
        log.Panicf("[PivotRoot] PivotRoot failed: %v", err)
        return fmt.Errorf("[PivotRoot] PivotRoot failed: %v", err)
    }
    // 修改当前的工作目录到根目录
    if err := syscall.Chdir("/"); err != nil {
        log.Panicf("[PivotRoot] chdir failed: %v", err)
        return fmt.Errorf("[PivotRoot] chdir failed: %v", err)
    }

    pivotDir = path.Join("/", ".pivot_root")
    // umount rootfs/.pivot_root
    if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
        log.Panicf("[PivotRoot] umount failed: %v", err)
        return fmt.Errorf("[PivotRoot] umount failed: %v", err)
    }
    // 删除临时文件夹
    return os.Remove(pivotDir)
}


/** pivot to "root" filesystem*/
// func PivotRoot1(root string) error {
//     oldRoot, err1 := syscall.Open("/", syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
//     newRoot, err2 := syscall.Open(root, syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
//     if err1 != nil || err2 != nil {
//         log.Panicf("[PivotRoot] open directory failed. %v && %v", err1, err2)
//     }
//     defer syscall.Close(oldRoot)
//     defer syscall.Close(newRoot)

//     err := syscall.Fchdir(newRoot)
//     if err != nil {
//         log.Panicf("[PivotRoot] change directory failed: %v", err)
//         return fmt.Errorf("[PivotRoot] change directory failed: %v", err)
//     }

//     // pivot root to newRoot, and put oldRoot to /proc/self/cwd.
//     // Concretely see github opencontainers/run master #1125
//     err = syscall.PivotRoot(".", ".")
//     if err != nil {
//         log.Panicf("[PivotRoot] pivot root failed: %v", err)
//         return fmt.Errorf("[PivotRoot] pivot root failed: %v", err)
//     }

//     err = syscall.Fchdir(oldRoot)
//     if err != nil {
//         log.Panicf("[PivotRoot] change directory failed: %v", err)
//         return fmt.Errorf("[PivotRoot] change directory failed: %v", err)
//     }

//     // Make oldroot private to make sure our unmounts don't propogate to the host
//     if err = syscall.Mount("", ".", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
//         log.Panicf("[PivotRoot] mount failed: %v", err)
//         return fmt.Errorf("[PivotRoot] mount failed: %v", err)
//     }

//     // Must be unmount . before chdir to /
//     err = syscall.Unmount(".", syscall.MNT_DETACH)
//     if err != nil {
//         log.Panicf("[PivotRoot] unmount failed: %v", err)
//         return fmt.Errorf("[PivotRoot] unmount failed: %v", err)
//     }

//     err = syscall.Chdir("/")
//     if err != nil {
//         log.Panicf("[PivotRoot] change directory failed: %v", err)
//         return fmt.Errorf("[PivotRoot] change directory failed: %v", err)
//     }
//     return nil
// }


