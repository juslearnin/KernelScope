package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"kernelscope/models"
)

func readTCPConnections() map[string]models.NetworkConnection {
	file, err := os.Open("/proc/net/tcp")
	
	if err != nil {
		return map[string]models.NetworkConnection{}
	}
	defer file.Close()

	connections := make(map[string]models.NetworkConnection)

	scanner := bufio.NewScanner(file)

	// Skip header
	if scanner.Scan() {
		// header ignored
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if len(fields) < 10 {
			continue
		}

		localAddress, localPort := parseAddress(fields[1])
		remoteAddress, remotePort := parseAddress(fields[2])
		state := parseTCPState(fields[3])
		inode := fields[9]

		connections[inode] = models.NetworkConnection{
			Inode:         inode,
			LocalAddress:  localAddress,
			LocalPort:     localPort,
			RemoteAddress: remoteAddress,
			RemotePort:    remotePort,
			State:         state,
		}
	}

	return connections
}

func parseAddress(value string) (string, int) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return "0.0.0.0", 0
	}

	ipHex := parts[0]
	portHex := parts[1]

	port64, _ := strconv.ParseInt(portHex, 16, 32)
	port := int(port64)

	if len(ipHex) != 8 {
		return "0.0.0.0", port
	}

	// IPv4 is little-endian in /proc/net/tcp
	b1 := ipHex[6:8]
	b2 := ipHex[4:6]
	b3 := ipHex[2:4]
	b4 := ipHex[0:2]

	n1, _ := strconv.ParseInt(b1, 16, 32)
	n2, _ := strconv.ParseInt(b2, 16, 32)
	n3, _ := strconv.ParseInt(b3, 16, 32)
	n4, _ := strconv.ParseInt(b4, 16, 32)

	ip := strconv.Itoa(int(n1)) + "." +
		strconv.Itoa(int(n2)) + "." +
		strconv.Itoa(int(n3)) + "." +
		strconv.Itoa(int(n4))

	return ip, port
}

func parseTCPState(code string) string {
	states := map[string]string{
		"01": "ESTABLISHED",
		"02": "SYN_SENT",
		"03": "SYN_RECV",
		"04": "FIN_WAIT1",
		"05": "FIN_WAIT2",
		"06": "TIME_WAIT",
		"07": "CLOSE",
		"08": "CLOSE_WAIT",
		"09": "LAST_ACK",
		"0A": "LISTEN",
		"0B": "CLOSING",
	}

	if state, ok := states[code]; ok {
		return state
	}

	return "UNKNOWN"
}

func matchConnections(
	files []models.FileDescriptor,
	tcpConnections map[string]models.NetworkConnection,
) []models.NetworkConnection {
	var connections []models.NetworkConnection

	for _, file := range files {
		if file.Type != "socket" {
			continue
		}

		inode := extractSocketInode(file.Target)
		if inode == "" {
			continue
		}

		connection, found := tcpConnections[inode]
		if !found {
			continue
		}

		connections = append(connections, connection)
	}

	return connections
}