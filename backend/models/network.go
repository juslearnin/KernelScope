package models

type NetworkConnection struct {
	Inode         string `json:"inode"`

	LocalAddress  string `json:"localAddress"`
	LocalPort     int    `json:"localPort"`

	RemoteAddress string `json:"remoteAddress"`
	RemotePort    int    `json:"remotePort"`

	State         string `json:"state"`
}