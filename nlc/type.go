package nlc

type HookFunc interface {
    Fork(ForkEvent)         error
    Exec(ExecEvent)         error
    Coredump(CoredumpEvent) error
    Exit(ExitEvent)         error
}

type Options struct {
    Hook    HookFunc
}

type NetlinkConnector struct {
    fd     int
    option Options
    ErrCh  chan error
}

type NetlinkMessageHeader struct {
    Length      uint32
    Type        uint16
    Flags       uint16
    Sequence    uint32
    Pid         uint32
}

type NetlinkConnectorMessageHeader struct {
    Index   uint32
    Value   uint32
    Message struct {
        Sequence    uint32
        Ack         uint32
        Length      uint16
        Flags       uint16
    }
}

type NetlinkConnectorMessagePacketHeader struct {
    What      uint32
    Cpu       uint32
    Timestamp uint64
}

type ForkEvent struct {
    ParentPID  uint32
    ParentTGID uint32
    ChildPID   uint32
    ChildTGID  uint32
}

type ExecEvent struct {
    ProcessPID  uint32
    ProcessTGID uint32
}

type CoredumpEvent struct {
    ProcessPID  uint32
    ProcessTGID uint32
    ParentPID   uint32
    ParentTGID  uint32
}

type ExitEvent struct {
    ProcessPID  uint32
    ProcessTGID uint32
    ExitCode    uint32
    ParentPID   uint32
    ParentTGID  uint32
}
