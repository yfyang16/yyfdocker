package subsystems

type CpusetSubSystem struct {}

func (sys *CpusetSubSystem) Name() string {
    return "cpuset"
}

func (sys *CpusetSubSystem) Set(cgroupPath string, cfg *ResourceConfig) error {
    return SubsysTemplateSet(cgroupPath, cfg, sys.Name(), "cpuset.cpus")
}

func (sys *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
    return SubsysTemplateApply(cgroupPath, pid, sys.Name())
}

func (sys *CpusetSubSystem) Remove(cgroupPath string) error {
    return SubsysTemplateRemove(cgroupPath, sys.Name())
}