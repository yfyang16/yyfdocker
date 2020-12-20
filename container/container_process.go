package container

import (
    "syscall"
    "os/exec"
    "os"
    "log"
    "path"
    "strings"
)

/** 
 * NewParentProcess: fork an isolated process which will execute the "yyfdocker init"
 * @para tty: if connect the stdin, stdout, stderr between my process and the os
 * @para volume: the path of host directory which will be mapped to /root/mnt
 */
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
    log.Printf("** NewParentProcess START **\n")
    defer log.Printf("** NewParentProcess END **\n")

    readPipe, writePipe, err := os.Pipe()
    if err != nil {
        log.Fatalf("[NewParentProcess] create anonymous pipe failure!")
        return nil, nil
    }

    cmd := exec.Command("/proc/self/exe", "init")

    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
        syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
    }
    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
	}

    cmd.ExtraFiles = []*os.File{readPipe}    // extra file descriptor besides in, out and err

    rootPath := "/root"; mntPath := "/root/mnt"
    NewWorkSpace(rootPath, mntPath, volume)
    cmd.Dir = mntPath

    return cmd, writePipe
}


/** Create an AUFS filesystem as a container root workspace */
func NewWorkSpace(rootPath string, mntPath string, volume string) {
    CreateReadOnlyLayer(rootPath, "busybox")
    CreateWriteOnlyLayer(rootPath)
    CreateMountPoint(rootPath, mntPath)
    if volume != "" {
        volumePaths := strings.Split(volume, ":")
        if len(volumePaths) == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
            MountVolume(rootPath, mntPath, volumePaths)
        } else {
            log.Panicf("[NewWorkSpace] volume parameter is incorrect: %s", volume)
        }

    }
}


func CreateReadOnlyLayer(rootPath string, imageName string) {
    imagePath := path.Join(rootPath, imageName)
    imageTarPath := path.Join(rootPath, imageName + ".tar")

    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        err := os.Mkdir(imagePath, 0777)
        if err != nil {
            log.Panicf("[CreateReadOnlyLayer] make directory failed: %v", err)
        }

        _, err = exec.Command("tar", "-xvf", imageTarPath, "-C", imagePath).CombinedOutput()
        if err != nil {
            log.Panicf("[CreateReadOnlyLayer] tar -xvf failed: %v", err)
        }
    }
}

func CreateWriteOnlyLayer(rootPath string) {
    if err := os.Mkdir(path.Join(rootPath, "writeLayer"), 0777); err != nil {
        log.Panicf("[CreateWriteOnlyLayer] make directory failed: %v", err)
    }
}

func CreateMountPoint(rootPath string, mntPath string) {
    if err := os.Mkdir(mntPath, 0777); err != nil {
        log.Panicf("[CreateMountPoint] make directory failed: %v", err)
    }

    // mount -t aufs -o dirs=/root/writeLayer:/root/busy/box none /root/mnt
    dirs := "dirs="+path.Join(rootPath, "writeLayer")+":"+path.Join(rootPath, "busybox")
    if  _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntPath).CombinedOutput(); err != nil {
        log.Panicf("[CreateMountPoint] mount cmd failed: %v", err)
    }
}

func MountVolume(rootPath string, mntPath string, volumePaths []string) {
    hostPath := volumePaths[0]; containerPath := volumePaths[1]

    containerVolumePath := path.Join(mntPath, containerPath)
    if err1 := os.Mkdir(hostPath, 0777);  err1 != nil {
        log.Panicf("[MountVolume] make directory failed: %v", err1)
    }
    if err2 := os.Mkdir(containerVolumePath, 0777); err2 != nil {
        log.Panicf("[MountVolume] make directory failed: %v", err2)
    }

    _, err := exec.Command("mount", "-t", "aufs", "-o", "dirs="+hostPath, "none", containerVolumePath).CombinedOutput()
    if err != nil {
        log.Panicf("[MountVolume] run mount command failed: %v", err)
    }
}

func DeleteWorkSpace(rootPath string, mntPath string, volume string) {
    if volume != "" {
        volumePaths := strings.Split(volume, ":")
        if len(volumePaths) == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
            DeleteMountPointWithVolume(rootPath, mntPath, volumePaths)
        } else {
            DeleteMountPoint(rootPath, mntPath)
        }
    } else {
        DeleteMountPoint(rootPath, mntPath)
    }
    DeleteWriteLayer(rootPath)
}

func DeleteMountPoint(rootPath string, mntPath string) {
    if _, err := exec.Command("umount", mntPath).CombinedOutput(); err != nil {
        log.Panicf("[DeleteMountPoint] umount failed: %v", err)
    }
    if err := os.RemoveAll(mntPath); err != nil {
        log.Panicf("[DeleteMountPoint] remove files failed: %v", err)
    }
}

func DeleteMountPointWithVolume(rootPath string, mntPath string, volumePaths []string) {
    if _, err := exec.Command("umount", path.Join(mntPath, volumePaths[1])).CombinedOutput(); err != nil {
        log.Panicf("[DeleteMountPointWithVolume] umount failed: %v", err)
    }
    DeleteMountPoint(rootPath, mntPath)
}

func DeleteWriteLayer(rootPath string) {
    if err := os.RemoveAll(path.Join(rootPath, "writeLayer")); err != nil {
        log.Panicf("[DeleteWriteLayer] remove files failed: %v", err)
    }
}


