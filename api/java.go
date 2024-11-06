package handler

import (
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "strings"
    "time"
)

// JavaStatusResponse 定义 Java 版服务器状态的结构体
type JavaStatusResponse struct {
    Version     string `json:"version"`
    Protocol    int    `json:"protocol"`
    Players     int    `json:"players"`
    MaxPlayers  int    `json:"max_players"`
    Description string `json:"description"`
    Favicon     string `json:"favicon,omitempty"`
}

// queryJavaServer 查询 Java 版 Minecraft 服务器
func queryJavaServer(host string, port int) (JavaStatusResponse, error) {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second*5)
    if err != nil {
        return JavaStatusResponse{}, err
    }
    defer conn.Close()

    // 构建协议包
    handshakePacket := []byte{0x00, 0x04, 0x09, byte(len(host))}
    handshakePacket = append(handshakePacket, []byte(host)...)
    handshakePacket = append(handshakePacket, byte(port>>8), byte(port), 0x01)

    conn.Write(handshakePacket)

    // 请求状态
    conn.Write([]byte{0x00})

    // 读取数据长度
    var length int32
    binary.Read(conn, binary.BigEndian, &length)

    // 读取数据包
    data := make([]byte, length)
    _, err = io.ReadFull(conn, data)
    if err != nil {
        return JavaStatusResponse{}, err
    }

    // 解码 JSON 数据
    var result JavaStatusResponse
    err = json.Unmarshal(data[3:], &result)
    return result, err
}

// JavaServerStatusHandler 是 Java 版请求的处理器
func JavaServerStatusHandler(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        http.Error(w, "Invalid path", http.StatusBadRequest)
        return
    }

    hostPort := pathParts[2]
    host, port, err := parseHostPort(hostPort, 25565)
    if err != nil {
        http.Error(w, "Invalid host or port", http.StatusBadRequest)
        return
    }

    result, err := queryJavaServer(host, port)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 输出 JSON 响应
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
