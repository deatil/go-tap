package main

import (
    "os"
    "fmt"
    "log"

    "github.com/urfave/cli/v2"

    "github.com/deatil/go-tap/ws"
    "github.com/deatil/go-tap/tcp"
    "github.com/deatil/go-tap/http"
)

const Version = "v0.0.2"

type IServer interface {
    Server()
}

// > go run main.go --type=tcp --addr=0.0.0.0:7755 --proxy=127.0.0.5:233
// > go run main.go --type=ws --addr=0.0.0.0:8082 --proxy=http://127.0.0.1:8002
// > go run main.go --type=http --addr=0.0.0.0:8083
// > ./proxy --type=tcp --addr=0.0.0.0:7755 --proxy=127.0.0.5:233
func main() {
    app := &cli.App{
        Name:      "proxy",
        Usage:     "proxy server",
        UsageText: "--type=tcp --proxy=127.0.0.1:9999,127.0.0.2:9999 [--addr=0.0.0.0:7755]",
        Version:   Version,
        EnableBashCompletion: true,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:    "conf",
                EnvVars: []string{""},
                Hidden:  true,
                Value:   "cfg.toml",
            },

            &cli.StringFlag{
                Name:    "type",
                Aliases: []string{"t"},
                Hidden:  false,
            },
            &cli.StringFlag{
                Name:   "addr",
                Hidden: false,
            },
            &cli.StringFlag{
                Name:   "proxy",
                Hidden: false,
            },
        },
        Action: func(ctx *cli.Context) error {
            runProxy(ctx)
            return nil
        },
        Commands: []*cli.Command{
            VersionCommand,
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Printf("ERROR: %s\n\n", err)
        os.Exit(1)
    }
}

// version
// > go run main.go version
var VersionCommand = &cli.Command{
    Name:      "version",
    Aliases:   []string{""},
    Usage:     "version info",
    UsageText: "version",
    Flags:  []cli.Flag{},
    Action: func(cctx *cli.Context) error {
        fmt.Printf("proxy version is %s \n", Version)

        return nil
    },
}

func runProxy(ctx *cli.Context) {
    addr := ctx.String("addr")
    proxy := ctx.String("proxy")
    if addr == "" {
        addr = "0.0.0.0:7755"
    }

    var s IServer
    var err error

    typ := ctx.String("type")
    switch typ {
        case "tcp":
            s, err = tcp.New(addr, proxy)
        case "ws":
            s, err = ws.New(addr, proxy)
        case "http":
            s, err = http.New(addr)
        default:
            log.Printf("%s unsupported", typ)
            return
    }

    if err != nil {
        log.Println(s)
        return
    }

    log.Printf("[%s] listen at %s", typ, addr)

    s.Server()
}
