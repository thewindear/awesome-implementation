package service

import (
    "bytes"
    "errors"
    "fmt"
    "log"
    "net"
    "strconv"
    "strings"
    "testing"
)

func TestConnectRedis(t *testing.T) {
    
    conn, err := net.Dial("tcp", "127.0.0.1:6379")
    if err != nil {
        t.Fatalf("connect error: %s", err)
    }
    defer conn.Close()
    
    err = auth(conn)
    if err != nil {
        t.Logf(err.Error())
    }
    
    err = send(conn, "set", "name", "good")
    log.Println(parseResponse(conn))
}

func auth(conn net.Conn) error {
    err := send(conn, "AUTH", "123456")
    if err != nil {
        return err
    }
    resp, err := parseResponse(conn)
    if err != nil {
        return err
    } else {
        if resp == "OK" {
            return nil
        }
        return fmt.Errorf("auth error %s", resp)
    }
}

func send(conn net.Conn, commands ...string) error {
    var write []byte
    var wBuf = bytes.NewBuffer(write)
    commandsLen := len(commands)
    wBuf.WriteString(fmt.Sprintf("*%d\r\n", commandsLen))
    
    for _, command := range commands {
        wBuf.WriteString(fmt.Sprintf("$%d\r\n", len(command)))
        wBuf.WriteString(fmt.Sprintf("%s\r\n", command))
    }
    _, err := conn.Write(wBuf.Bytes())
    if err != nil {
        return err
    }
    return nil
}

func parseResponse(conn net.Conn) (resp string, err error) {
    var read = make([]byte, 1024)
    rBuf := bytes.NewBuffer(read)
    var readLen int
    //解析响应
    readLen, err = conn.Read(read)
    prefix := rBuf.Next(1)
    
    var content = string(rBuf.Next(readLen - 1))
    
    strs := strings.Split(content, "\r\n")
    
    switch string(prefix[0]) {
    case "+": //直接读取单行
        resp = strings.Join(strs, "")
        break
    case "-":
        err = errors.New(strings.Join(strs, ""))
        break
    case ":":
        resp = strs[0]
        break
    case "$":
        //多行字符串先读取长度再读取数据
        strLen, _ := strconv.Atoi(strs[0])
        if strLen == 0 {
            resp = ""
        }
        resp = strings.Join(strs[1:len(strs)-1], "")
        break
    case "*":
        break
    default:
        err = errors.New("invalid data type")
    }
    return
}
