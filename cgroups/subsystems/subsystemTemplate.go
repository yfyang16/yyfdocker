package subsystems

import (
    "io/ioutil"
    "path"
    "os"
    "fmt"
    "strconv"
    "log"
)


func SubsysTemplateSet(cgroupPath string, cfg *ResourceConfig, subsysName string, fileName string) error {
    log.Printf("** %sSubsystem:Set START **\n", subsysName)
    defer log.Printf("** %sSubsystem:Set END **\n", subsysName)

    valueLimited := ""
    switch subsysName {
    case "memory":
        valueLimited = cfg.MemoryLimit
    case "cpuset":
        valueLimited = cfg.CpuSet
    case "cpu":
        valueLimited = cfg.CpuShare
    }

    subsysCgroupPath, err := GetCgroupPath(subsysName, cgroupPath, true)
    if err == nil {
        if valueLimited != "" {
            limitConfigFile := path.Join(subsysCgroupPath, fileName)
            err := ioutil.WriteFile(limitConfigFile, []byte(valueLimited), 0644)
            if err != nil {
                log.Panicf("[%sSubsystem:Set] Cannot write into file: %v", subsysName, err)
                return fmt.Errorf("[%sSubsystem:Set] Cannot write into file: %v", subsysName, err)
            } else {
                return nil
            }
        }
    } else {
        log.Panicf("[%sSubsystem:Set] GetCgroupPath throws error: %v", subsysName, err)
    }

    return err
}

func SubsysTemplateApply(cgroupPath string, pid int, subsysName string) error {
    log.Printf("** %sSubsystem:Apply START **\n", subsysName)
    defer log.Printf("** %sSubsystem:Apply END **\n", subsysName)

    subsysCgroupPath, err := GetCgroupPath(subsysName, cgroupPath, false)
    if err == nil {
        taskProcessFile := path.Join(subsysCgroupPath, "tasks")
        err := ioutil.WriteFile(taskProcessFile, []byte(strconv.Itoa(pid)), 0644)
        if err != nil {
            log.Panicf("[%sSubsystem:Apply] Write file error: %v", subsysName, err)
            return fmt.Errorf("[%sSubsystem:Apply] Write file error: %v", subsysName, err)
        } else {
            return nil
        }
    } else {
        log.Panicf("[%sSubsystem:Apply] GetCgroupPath throws error: %v", subsysName, err)
        return fmt.Errorf("[%SSubsystem:Apply] GetCgroupPath throws error: %v", subsysName, err)
    }
}

func SubsysTemplateRemove(cgroupPath string, subsysName string) error {
    log.Printf("** %sSubsystem:Remove START **\n", subsysName)
    defer log.Printf("** %sSubsystem:Remove END **\n")

    subsysCgroupPath, err := GetCgroupPath(subsysName, cgroupPath, false)
    if err == nil {
        return os.RemoveAll(subsysCgroupPath)
    } else {
        log.Panicf("[%sSubsystem:Remove] GetCgroupPath throws error: %v", subsysName, err)
        return fmt.Errorf("[%sSubsystem:Remove] GetCgroupPath throws error: %v", subsysName, err)
    }

}