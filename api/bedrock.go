package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// BedrockStatusResponse 定义基岩版服务器状态的结构体
type BedrockStatusResponse struct {
	ServerID    string `json:"server_id"`
	MOTD        string `json:"motd"`
	Gamemode    string `json:"gamemode"`
	PlayerCount int    `json:"player_count"`
	MaxPlayers  int    `json:"max_players"`
	Version     string `json:"version"`
}

// queryBedrockServer 查询基岩版 Minecraft 服务器
func queryBedrockServer(host string, port int) (BedrockStatusResponse, error) {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", host, port), time.Second*5)
	if err != nil {
		return BedrockStatusResponse{}, err
	}
	defer conn.Close()

	// 发送 ping 请求
	packet := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	conn.Write(packet)

	// 接收响应
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		return BedrockStatusResponse{}, err
	}

	// 解析数据
	data := buffer[35:n]
	fields := string(data)
	var response BedrockStatusResponse
	fmt.Sscanf(fields, "%s;%s;%s;%d;%d;%s", &response.ServerID, &response.MOTD, &response.Gamemode, &response.PlayerCount, &response.MaxPlayers, &response.Version)

	return response, nil
}

// BedrockServerStatusHandler 是基岩版请求的处理器
func BedrockServerStatusHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	hostPort := pathParts[2]
	host, port, err := parseHostPort(hostPort, 19132)
	if err != nil {
		http.Error(w, "Invalid host or port", http.StatusBadRequest)
		return
	}

	result, err := queryBedrockServer(host, port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 输出 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

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
