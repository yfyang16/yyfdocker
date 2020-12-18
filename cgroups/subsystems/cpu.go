package subsystems

type CpuSubSystem struct {}

func (sys *CpuSubSystem) Name() string {
    return "cpu"
}

func (sys *CpuSubSystem) Set(cgroupPath string, cfg *ResourceConfig) error {
    return SubsysTemplateSet(cgroupPath, cfg, sys.Name(), "cpu.shares")
}

func (sys *CpuSubSystem) Apply(cgroupPath string, pid int) error {
    return SubsysTemplateApply(cgroupPath, pid, sys.Name())
}

func (sys *CpuSubSystem) Remove(cgroupPath string) error {
    return SubsysTemplateRemove(cgroupPath, sys.Name())
}