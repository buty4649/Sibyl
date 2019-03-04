package nlc

import (
    "bytes"
    "context"
    "os"
    "unsafe"
    "log"

    "encoding/binary"
    "golang.org/x/sys/unix"
)

func Init(opt Options) (*NetlinkConnector, error) {
    fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, NETLINK_CONNECTOR)
    if err != nil {
        return nil, err
    }
    nlc := NetlinkConnector{fd:fd, option:opt}

    if err = nlc.regist(); err != nil {
        nlc.Close()
        return nil, err
    }

    return &nlc, nil
}

func (nlc *NetlinkConnector) regist() error {
    sockaddr := unix.SockaddrNetlink{
        Family: unix.AF_NETLINK,
        Groups: CN_IDX_PROC,
        Pid:    uint32(os.Getpid()),
    }

    if err := unix.Bind(nlc.fd, &sockaddr); err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    if err := buildPacket(buf); err != nil {
        return err
    }

    if err := unix.Sendto(nlc.fd, buf.Bytes(), 0, &sockaddr); err != nil {
        return err
    }

    return nil
}

func buildPacket(buf *bytes.Buffer) (error) {
    nlcmh := NetlinkConnectorMessageHeader{}
    nlmh := NetlinkMessageHeader{}
    var data uint32 = PROC_CN_MCAST_LISTEN

    nlmh.Length = uint32(unsafe.Sizeof(nlmh) + unsafe.Sizeof(nlcmh) + unsafe.Sizeof(data))
    nlmh.Type = NLMSG_DONE
    nlmh.Pid  = uint32(os.Getpid())
    if err := binary.Write(buf, binary.LittleEndian, &nlmh); err != nil {
        return err
    }

    nlcmh.Index = CN_IDX_PROC
    nlcmh.Value = CN_VAL_PROC
    nlcmh.Message.Length = uint16(unsafe.Sizeof(data))
    if err := binary.Write(buf, binary.LittleEndian, &nlcmh); err != nil {
        return err
    }

    if err := binary.Write(buf, binary.LittleEndian, &data); err != nil {
        return err
    }

    return nil
}

func (nlc *NetlinkConnector) Close() {
    unix.Close(nlc.fd)
}

func (nlc *NetlinkConnector) Start(ctxParent context.Context) {
    ctx, cancel := context.WithCancel(ctxParent)

    chanSelect := make(chan bool)

    go func() {
        defer cancel()
        defer close(chanSelect)

        for {
            var readfds unix.FdSet

            FD_ZERO(&readfds)
            FD_SET(nlc.fd, &readfds)
            _, err := unix.Select(nlc.fd+1, &readfds, nil, nil, nil)
            if err != nil {
                log.Fatal(err)
                break
            }

            if FD_ISSET(nlc.fd, &readfds) {
                chanSelect <- true
            }
        }
    }()

    go func() {
        for {
            select {
            case <-ctx.Done():
                break

            case <-chanSelect:
                reader, err := nlc.recv()
                if err != nil {
                    log.Fatal(err)
                    ctx.Err()
                    break
                }
                if reader == nil {
                    continue
                }

                for {
                    ph, next := readPacket(reader)

                    switch ph.What {
                    case PROC_EVENT_FORK:
                        var ev ForkEvent
                        err = binary.Read(reader, binary.LittleEndian, &ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }

                        err = nlc.option.Hook.Fork(ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }
                    case PROC_EVENT_EXEC:
                        var ev ExecEvent
                        err = binary.Read(reader, binary.LittleEndian, &ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }

                        err = nlc.option.Hook.Exec(ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }
                    case PROC_EVENT_COREDUMP:
                        var ev CoredumpEvent
                        err = binary.Read(reader, binary.LittleEndian, &ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }

                        err = nlc.option.Hook.Coredump(ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }
                    case PROC_EVENT_EXIT:
                        var ev ExitEvent
                        err = binary.Read(reader, binary.LittleEndian, &ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }

                        err = nlc.option.Hook.Exit(ev)
                        if err != nil {
                            log.Fatal(err)
                            continue
                        }
                    }

                    if next == false {
                        break
                    }
                }
            }
        }
    }()
}

func (nlc *NetlinkConnector) recv() (*bytes.Reader, error) {
    buf := make([]byte, 1024)

    n, sa, err := unix.Recvfrom(nlc.fd, buf, 0)
    if err != nil || n < 1 {
        return nil, err
    }

    var sanl = sa.(*unix.SockaddrNetlink)
    if sanl.Groups != CN_IDX_PROC || sanl.Pid != 0 {
        return nil, nil
    }

    return bytes.NewReader(buf), nil
}

func readPacket(r *bytes.Reader) (ph *NetlinkConnectorMessagePacketHeader, next bool) {
    next = false

    for {
        var nlmh NetlinkMessageHeader
        if binary.Read(r, binary.LittleEndian, &nlmh) != nil{
            // EOF
            break
        }
        if nlmh.Type == NLMSG_ERROR || nlmh.Type == NLMSG_OVERRUN {
            break
        }
        if nlmh.Type == NLMSG_NOOP {
            // Drop unnecessary data
            readPaddingData(r, nlmh.Length - uint32(unsafe.Sizeof(nlmh)))
            continue
        }

        // Drop unnecessary data
        readPaddingData(r, uint32(unsafe.Sizeof(NetlinkConnectorMessageHeader{})))

        ph = &NetlinkConnectorMessagePacketHeader{}
        if binary.Read(r, binary.LittleEndian, ph) != nil {
            // EOF
            break
        }

        if nlmh.Type == NLMSG_DONE {
            break
        }

        next = true
        break
    }
    return ph, next
}

func readPaddingData(r *bytes.Reader,n uint32) {
    buf := make([]byte, n)
    r.Read(buf)
}

func FD_SET(i int,p *unix.FdSet) {
    p.Bits[i/64] |= 1 << uint(i) % 64
}

func FD_ISSET(i int,p *unix.FdSet) bool {
	return (p.Bits[i/64] & (1 << uint(i) % 64)) != 0
}

func FD_ZERO(p *unix.FdSet) {
    for i := range p.Bits {
        p.Bits[i] = 0
    }
}
