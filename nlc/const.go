package nlc

const (
    NETLINK_CONNECTOR = 11

    CN_IDX_PROC = 0x1
    CN_VAL_PROC = 0x1

    PROC_CN_MCAST_LISTEN = 1

    NLMSG_NOOP    = 0x1
    NLMSG_ERROR   = 0x2
    NLMSG_DONE    = 0x3
    NLMSG_OVERRUN = 0x4

    PROC_EVENT_NONE     = 0x00000000
    PROC_EVENT_FORK     = 0x00000001
    PROC_EVENT_EXEC     = 0x00000002
    PROC_EVENT_UID      = 0x00000004
    PROC_EVENT_GID      = 0x00000040
    PROC_EVENT_SID      = 0x00000080
    PROC_EVENT_PTRACE   = 0x00000100
    PROC_EVENT_COMM     = 0x00000200
    PROC_EVENT_COREDUMP = 0x40000000
    PROC_EVENT_EXIT     = 0x80000000
)
