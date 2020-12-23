package container

import (
    "syscall"
    "os/exec"
    "os"
    "log"
    "path"
    "strings"
    "fmt"
)

var (
    RUNNING string             = "running"
    STOP string                = "stopped"
    EXIT string                = "exited"
    DefaultInfoLocation string = "/var/run/yyfdocker/%s/"
    ConfigName string          = "config.txt"
    ContainerLog string        = "container.log"
    RootPath string            = "/root"
    MntPath string             = "/root/mnt"
    WriteLayerPath string      = "/root/writeLayer"
)

type ContainerInfo struct {
    Pid         string
    Id          string
    Name        string
    Command     string
    CreatedTime string
    Status      string
    Volume      string
}


/** 
 * NewParentProcess: fork an isolated process which will execute the "yyfdocker init"
 * @para tty: if connect the stdin, stdout, stderr between my process and the os
 * @para volume: the path of host directory which will be mapped to /root/mnt
 */
func NewParentProcess(tty bool, volume string, containerName string, imageName string) (*exec.Cmd, *os.File) {
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
    if tty {                                 // -it
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
	} else {
        logDirPath := fmt.Sprintf(DefaultInfoLocation, containerName)
        if err = os.MkdirAll(logDirPath, 0622); err != nil {
            log.Panicf("[NewParentProcess] make all directories failed: %v", err)
            return nil, nil
        }
        logFilePath := path.Join(logDirPath, ContainerLog)
        logFile, err := os.Create(logFilePath)
        if err != nil {
            log.Panicf("[NewParentProcess] create log file error: %v", err)
            return nil, nil
        }
        cmd.Stdout = logFile
    }

    cmd.ExtraFiles = []*os.File{readPipe}    // extra file descriptor besides in, out and err

    NewWorkSpace(volume, imageName, containerName)
    cmd.Dir = path.Join(MntPath, containerName)

    return cmd, writePipe
}


/** Create an AUFS filesystem as a container root workspace */
func NewWorkSpace(volume string, imageName string, containerName string) {
    CreateReadOnlyLayer(imageName)
    CreateWriteOnlyLayer(containerName)
    CreateMountPoint(containerName, imageName)
    if volume != "" {
        volumePaths := strings.Split(volume, ":")
        if len(volumePaths) == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
            MountVolume(volumePaths, containerName)
        } else {
            log.Panicf("[NewWorkSpace] volume parameter is incorrect: %s", volume)
        }

    }
}


func CreateReadOnlyLayer(imageName string) {
    imagePath := path.Join(RootPath, imageName)
    imageTarPath := path.Join(RootPath, imageName + ".tar")

    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        err := os.MkdirAll(imagePath, 0777)
        if err != nil {
            log.Panicf("[CreateReadOnlyLayer] make directory failed: %v", err)
        }

        _, err = exec.Command("tar", "-xvf", imageTarPath, "-C", imagePath).CombinedOutput()
        if err != nil {
            log.Panicf("[CreateReadOnlyLayer] tar -xvf failed: %v", err)
        }
    }
}

func CreateWriteOnlyLayer(containerName string) {
    if err := os.MkdirAll(path.Join(WriteLayerPath, containerName), 0777); err != nil {
        log.Printf("[CreateWriteOnlyLayer] make directory failed: %v", err)
    }
}

func CreateMountPoint(containerName string, imageName string) {
    if err := os.MkdirAll(path.Join(MntPath, containerName), 0777); err != nil {
        log.Printf("[CreateMountPoint] make directory failed: %v", err)
    }

    // mount -t aufs -o dirs=/root/writeLayer:/root/busy/box none /root/mnt
    containerMntPath := path.Join(MntPath, containerName)
    imagePath := path.Join(RootPath, imageName)
    dirs := "dirs=" + path.Join(WriteLayerPath, containerName) + ":" + imagePath
    if  _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerMntPath).CombinedOutput(); err != nil {
        log.Panicf("[CreateMountPoint] mount cmd failed: %v", err)
    }
}

func MountVolume(volumePaths []string, containerName string) {
    hostPath := volumePaths[0]; containerPath := volumePaths[1]

    containerVolumePath := path.Join(path.Join(MntPath, containerName), containerPath)
    if err1 := os.Mkdir(hostPath, 0777);  err1 != nil {
        log.Printf("[MountVolume] make directory failed: %v", err1)
    }
    if err2 := os.Mkdir(containerVolumePath, 0777); err2 != nil {
        log.Printf("[MountVolume] make directory failed: %v", err2)
    }

    _, err := exec.Command("mount", "-t", "aufs", "-o", "dirs="+hostPath, "none", containerVolumePath).CombinedOutput()
    if err != nil {
        log.Panicf("[MountVolume] run mount command failed: %v", err)
    }
}

func DeleteWorkSpace(volume string, containerName string) {
    if volume != "" {
        volumePaths := strings.Split(volume, ":")
        if len(volumePaths) == 2 && volumePaths[0] != "" && volumePaths[1] != "" {
            DeleteMountPointWithVolume(volumePaths, containerName)
        } else {
            DeleteMountPoint(containerName)
        }
    } else {
        DeleteMountPoint(containerName)
    }
    DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) {
    if _, err := exec.Command("umount", "-l", path.Join(MntPath, containerName)).CombinedOutput(); err != nil {
        log.Panicf("[DeleteMountPoint] umount failed: %v", err)
    }
    if err := os.RemoveAll(path.Join(MntPath, containerName)); err != nil {
        log.Panicf("[DeleteMountPoint] remove files failed: %v", err)
    }
}

func DeleteMountPointWithVolume(volumePaths []string, containerName string) {
    containerMntPath := path.Join(MntPath, containerName)
    if _, err := exec.Command("umount", "-l", path.Join(containerMntPath, volumePaths[1])).CombinedOutput(); err != nil {
        log.Printf("[DeleteMountPointWithVolume] umount failed: %v", err)
    }
    DeleteMountPoint(containerName)
}

func DeleteWriteLayer(containerName string) {
    if err := os.RemoveAll(path.Join(WriteLayerPath, containerName)); err != nil {
        log.Panicf("[DeleteWriteLayer] remove files failed: %v", err)
    }
}


