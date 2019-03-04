package main

import (
    "context"
    "flag"
    "log"
    "io/ioutil"

    "github.com/buty4649/Sibyl/nlc"
)

type Sibyl struct {
    config  Config
}

var sibyl *Sibyl

func main() {
    flag.Parse()
    yamlfile, err := ioutil.ReadFile(flag.Arg(0))
    if err != nil {
        log.Fatal(err)
    }

    config, err := LoadConfig(yamlfile)
    if err != nil {
        log.Fatal(err)
    }

    sibyl = &Sibyl{}
    sibyl.config = *config

    var opts nlc.Options
    opts.Hook = sibyl

    n, err := nlc.Init(opts)
    if err != nil {
        log.Fatal(err)
    }
    defer n.Close()

    ctx := context.Background()
    n.Start(ctx)
    <-ctx.Done()
}
