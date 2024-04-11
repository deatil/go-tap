package main

import (
    "os"
    "fmt"
    "log"

    "github.com/urfave/cli/v2"

    "github.com/deatil/go-tap/ws"
    "github.com/deatil/go-tap/tcp"
)

const Version = "v0.0.1"

// > go run main.go --type=tcp --src=0.0.0.0:7755 --dst=127.0.0.5:233
// > go run main.go --type=ws --src=0.0.0.0:8082 --dst=http://127.0.0.1:8002
// > ./proxy --type=tcp --src=0.0.0.0:7755 --dst=127.0.0.5:233
func main() {
    app := &cli.App{
        Name:      "proxy",
        Usage:     "proxy server",
        UsageText: "--type=tcp --dst=127.0.0.1:9999,127.0.0.2:9999 [--src=0.0.0.0:7755]",
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
                Name:   "src",
                Hidden: false,
            },
            &cli.StringFlag{
                Name:   "dst",
                Hidden: false,
            },
        },
        Action: func(ctx *cli.Context) error {
            src := ctx.String("src")
            dst := ctx.String("dst")
            if dst == "" {
                fmt.Println("need dst flag")
                return nil
            }

            if src == "" {
                dst = "0.0.0.0:7755"
            }

            log.Println("listen at:", src)

            typ := ctx.String("type")
            switch typ {
                case "tcp":
                    tcp.New(src, dst).Server()
                case "ws":
                    ws.New(src, dst).Server()
                default:
                    log.Println("type unsupported")
            }

            return nil
        },
        Commands: []*cli.Command{
            VersionCommand,
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Println("ERROR: %s\n\n", err)
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
