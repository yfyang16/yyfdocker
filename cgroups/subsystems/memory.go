package subsystems

type MemorySubSystem struct {}

func (sys *MemorySubSystem) Name() string {
    return "memory"
}

func (sys *MemorySubSystem) Set(cgroupPath string, cfg *ResourceConfig) error {
    return SubsysTemplateSet(cgroupPath, cfg, sys.Name(), "memory.limit_in_bytes")
}

func (sys *MemorySubSystem) Apply(cgroupPath string, pid int) error {
    return SubsysTemplateApply(cgroupPath, pid, sys.Name())
}

func (sys *MemorySubSystem) Remove(cgroupPath string) error {
    return SubsysTemplateRemove(cgroupPath, sys.Name())
}