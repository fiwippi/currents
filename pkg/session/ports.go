package session

import (
	"go.bug.st/serial/enumerator"
)

func GetAvailablePorts() ([]string, error) {
	portData, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}

	ports := make([]string, 0)
	for _, port := range portData {
		if port.IsUSB {
			ports = append(ports, port.Name)
		}
	}
	return ports, nil
}
