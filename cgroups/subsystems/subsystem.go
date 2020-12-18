package subsystems

type ResourceConfig struct {
	MemoryLimit string                          // memory limit
	CpuShare    string                          // weight of cpu time share
	CpuSet      string                          // core number
}

type Subsystem interface {
	Name() string                               // subsystem name
	Set(path string, res *ResourceConfig) error // set the resource configure of a cgroup
	Apply(path string, pid int) error           // add a process to this cgroup
	Remove(path string) error                   // remove a cgroup
}

/** SubsystemsIns is a list of several subsystems.*/
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
