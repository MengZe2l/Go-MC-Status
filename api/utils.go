// api/utils.go
package handler

import (
    "fmt"
    "strings"
)

// parseHostPort 解析主机和端口
func parseHostPort(hostPort string, defaultPort int) (string, int, error) {
    parts := strings.Split(hostPort, ":")
    host := parts[0]
    port := defaultPort
    if len(parts) > 1 {
        _, err := fmt.Sscanf(parts[1], "%d", &port)
        if err != nil {
            return "", 0, err
        }
    }
    return host, port, nil
}
