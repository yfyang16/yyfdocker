package subsystems

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

/** Find the mount point (prefix of an absolute path) of a subsystem:
 ** e.g. memory, cpushare, cpuset
 */
func FindCgroupMountPoint(subsys string) string {
	log.Printf("** FindCgroupMountPoint START **\n")
    defer log.Printf("** FindCgroupMountPoint END **\n")

	/** The mount information of the current process can be found in
	 ** "/proc/self/mountinfo" file. The content of this file is like:
	 ** ...
	 ** 30 27 0:24 I /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:l3 - cgroup cgroup rw,memory
	 ** 31 27 0:25 I /sys/fs/cgroup/freezer rw, nosuid, nodev, noexec, relatime shared:l4 - cgroup cgroup rw,freezer
	 ** 32 27 0:26 I /sys/fs/cgroup/hugetlb rw,nosuid,nodev,noexec,relatime shared:lS - cgroup cgroup rw,hugetlb
	 ** ...
	 */
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Println("[FindCgroupMountPoint] Cannot open /proc/self/mountinfo!")
		return ""
	}
	defer f.Close()

	// Iterate every line of this file and look for the corresponding subsystem's cgroup path.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()                      // line text
		fieldList := strings.Split(txt, " ")       // split by space: [30, 27, 0:24, I, /sys/fs/cgroup/memroy, ...]
		lastOption := fieldList[len(fieldList)-1]  // "rw,memory"
		for _, opt := range strings.Split(lastOption, ",") {
			if opt == subsys {
				return fieldList[4]
			}                                      // return "memory"
		}
	}

	err = scanner.Err()
	if err != nil {
		log.Printf("[FindCgroupMountPoint] Scanner Error: %s!\n", err.Error())
		return ""
	}

	return ""                                      // do not find the corresponding subsystem
}

/** Get the complete cgroup path */
func GetCgroupPath(subsys string, cgroupPath string, autoCreate bool) (string, error) {
    log.Printf("** GetCgroupPath START **\n")
    defer log.Printf("** GetCgroupPath END **\n")

	cgroupRoot := FindCgroupMountPoint(subsys)
	completeCgroupPath := path.Join(cgroupRoot, cgroupPath)
	_, err := os.Stat(completeCgroupPath)

    log.Printf("[GetCgroupPath] autoCreate: %v; error_bool: %v", autoCreate, os.IsNotExist(err))

	if err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			err := os.Mkdir(completeCgroupPath, 0755)
			if err != nil {
				log.Panicf("[GetCgroupPath] Cannot create path: %v\n", err)
				return "", fmt.Errorf("[GetCgroupPath] Cannot create path: %v", err)
			}
		}
		return completeCgroupPath, nil
	} else {
		log.Panicf("[GetCgroupPath] Cgroup path error: %v\n", err)
		return "", fmt.Errorf("[GetCgroupPath] Cgroup path error: %v", err)
	}
}
