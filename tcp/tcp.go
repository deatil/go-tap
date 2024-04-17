package tcp

import (
    "io"
    "log"
    "net"
    "fmt"
    "sync"
    "errors"
    "strings"
)

type tcp struct {
    lock   sync.Mutex
    addr   string
    proxys []string
}

func New(addr, proxy string) (*tcp, error) {
    if addr == "" {
        return nil, errors.New("need addr flag")
    }
    if proxy == "" {
        return nil, errors.New("need proxy flag")
    }

    tp := &tcp{
        addr:   addr,
        proxys: strings.Split(proxy, ","),
    }

    return tp, nil
}

func (p *tcp) Server() {
    listen, err := net.Listen("tcp", p.addr)
    if err != nil {
        fmt.Println(err)
        return
    }

    defer listen.Close()

    for {
        conn, err := listen.Accept()
        if err != nil {
            log.Printf("client: %s, local: %s, listen accept error: %v\n", conn.RemoteAddr(), conn.LocalAddr(), err)
            continue
        }

        go p.handle(conn)
    }
}

func (p *tcp) handle(sconn net.Conn) {
    defer sconn.Close()
    proxy, ok := p.selectProxy()
    if !ok {
        return
    }

    dconn, err := net.Dial("tcp", proxy)
    if err != nil {
        log.Printf("dial %v fail: %v\n", proxy, err)
        return
    }
    defer dconn.Close()

    ExitChan := make(chan bool, 1)

    // 转发到目标服务器
    go func(sconn net.Conn, dconn net.Conn, exit chan bool) {
        _, err := io.Copy(dconn, sconn)
        if err != nil {
            log.Printf("give message to %v fail: %v\n", proxy, err)
            ExitChan <- true
        }
    }(sconn, dconn, ExitChan)

    // 从目标服务器返回数据到客户端
    go func(sconn net.Conn, dconn net.Conn, exit chan bool) {
        _, err := io.Copy(sconn, dconn)
        if err != nil {
            log.Printf("get message fail from %v: %v\n", proxy, err)
            ExitChan <- true
        }
    }(sconn, dconn, ExitChan)

    <-ExitChan
}

// ip 轮询
func (p *tcp) selectProxy() (string, bool) {
    p.lock.Lock()
    defer p.lock.Unlock()

    if len(p.proxys) < 1 {
        return "", false
    }

    proxy := p.proxys[0]
    p.proxys = append(p.proxys[1:], proxy)
    return proxy, true
}
