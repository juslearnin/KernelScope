package tree

import (
	"fmt"
	"sort"
	"strings"

	"kernelscope/models"
)

func BuildProcessTree(processes []models.Process) map[int][]models.Process {
	tree := make(map[int][]models.Process)

	for _, process := range processes {
		tree[process.PPID] = append(tree[process.PPID], process)
	}

	for ppid := range tree {
		sort.Slice(tree[ppid], func(i, j int) bool {
			return tree[ppid][i].PID < tree[ppid][j].PID
		})
	}

	return tree
}

func PrintTree(pid int, childrenByPPID map[int][]models.Process, processes []models.Process, depth int) {
	process, found := findProcessByPID(processes, pid)
	if !found {
		return
	}

	indent := strings.Repeat("  ", depth)

	displayName := process.Name
	if process.Cmdline != "" {
		displayName = process.Cmdline
	}

	fmt.Printf("%s└── %s [pid=%d, ram=%dKB, threads=%d, cpu=%.2f]\n",
		indent,
		displayName,
		process.PID,
		process.MemoryKB,
		process.Threads,
		process.CPUPercent,
	)

	children := childrenByPPID[pid]

	for _, child := range children {
		PrintTree(child.PID, childrenByPPID, processes, depth+1)
	}
}

func findProcessByPID(processes []models.Process, pid int) (models.Process, bool) {
	for _, process := range processes {
		if process.PID == pid {
			return process, true
		}
	}

	return models.Process{}, false
} 