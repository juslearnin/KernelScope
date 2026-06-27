package collector

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"fmt"

	"kernelscope/models"
)

func CollectProcesses() ([]models.Process, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var processes []models.Process

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		process, err := readProcessStatus(pid)
		if err != nil {
			continue
		}

		processes = append(processes, process)
	}

	return processes, nil
}

func readProcessStatus(pid int) (models.Process, error) {
    statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")

    file, err := os.Open(statusPath)
    if err != nil {
        return models.Process{}, err
    }
    defer file.Close()

    process := models.Process{PID: pid}

    // 🚨 MOVE THESE THREE LINES UP HERE (Outside and before the text scanner loop)
    process.Cmdline = readCmdline(pid)
    process.CPUTime = readCPUTime(pid)
    process.OpenFiles = readFileDescriptors(pid)
    fmt.Printf("=== PID %d -> %d open files tracked ===\n", pid, len(process.OpenFiles))

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        switch {
        case strings.HasPrefix(line, "Name:"):
            process.Name = cleanValue(line, "Name:")
        case strings.HasPrefix(line, "State:"):
            process.State = cleanValue(line, "State:")
        case strings.HasPrefix(line, "PPid:"):
            process.PPID = parseIntValue(line, "PPid:")
        case strings.HasPrefix(line, "Threads:"):
            process.Threads = parseIntValue(line, "Threads:")
        case strings.HasPrefix(line, "VmRSS:"):
            process.MemoryKB = parseMemoryKB(line)
        }
    }

    return process, scanner.Err()
}

func readCmdline(pid int) string {
	cmdlinePath := filepath.Join("/proc", strconv.Itoa(pid), "cmdline")

	data, err := os.ReadFile(cmdlinePath)
	if err != nil || len(data) == 0 {
		return ""
	}

	cmdline := strings.ReplaceAll(string(data), "\x00", " ")
	return strings.TrimSpace(cmdline)
}

func readCPUTime(pid int) uint64 {
	statPath := filepath.Join("/proc", strconv.Itoa(pid), "stat")

	data, err := os.ReadFile(statPath)
	if err != nil {
		return 0
	}

	content := string(data)

	closingParen := strings.LastIndex(content, ")")
	if closingParen == -1 {
		return 0
	}

	afterName := strings.TrimSpace(content[closingParen+1:])
	fields := strings.Fields(afterName)

	if len(fields) < 13 {
		return 0
	}

	utime, _ := strconv.ParseUint(fields[11], 10, 64)
	stime, _ := strconv.ParseUint(fields[12], 10, 64)

	return utime + stime
}

func cleanValue(line string, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func parseIntValue(line string, prefix string) int {
	value := cleanValue(line, prefix)
	num, _ := strconv.Atoi(value)
	return num
}

func parseMemoryKB(line string) int {
	value := cleanValue(line, "VmRSS:")
	parts := strings.Fields(value)

	if len(parts) == 0 {
		return 0
	}

	num, _ := strconv.Atoi(parts[0])
	return num
}
func readFileDescriptors(pid int) []models.FileDescriptor {
    fdPath := filepath.Join("/proc", strconv.Itoa(pid), "fd")
    
    // 1. CHANGE THIS LINE: Always initialize to a clean, empty slice
    descriptors := make([]models.FileDescriptor, 0)

    entries, err := os.ReadDir(fdPath)
    if err != nil {
        // 2. CHANGE THIS LINE: Return the empty slice instead of nil
        return descriptors 
    }

    for _, entry := range entries {
        fd, err := strconv.Atoi(entry.Name())
        if err != nil {
            continue
        }
        
        linkPath := filepath.Join("/proc", strconv.Itoa(pid), "fd", strconv.Itoa(fd))
        target, err := os.Readlink(linkPath)
        if err != nil {
            continue
        }
        
        descriptors = append(descriptors, models.FileDescriptor{
            FD:     fd,
            Target: target,
            Type:   detectFDType(target),
        })
    }

    return descriptors
}

func detectFDType(target string) string {
	switch {
	case strings.HasPrefix(target, "socket:"):
		return "socket"
	case strings.HasPrefix(target, "pipe:"):
		return "pipe"
	case strings.HasPrefix(target, "anon_inode:"):
		return "kernel"
	case strings.HasPrefix(target, "/dev/"):
		return "device"
	case strings.HasPrefix(target, "/proc/"):
		return "proc"
	case strings.HasPrefix(target, "/"):
		return "file"
	default:
		return "unknown"
	}
}