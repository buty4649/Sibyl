# Sibyl
Sibyl is useful process event hook command.

## Usage

```
$ cat sibyl.yaml.sample
coredump:
  - commandline: .*
    exec: echo [$SIBYL_PID] $SIBYL_CMDLINE
exec:
  - commandline: .*
    exec: echo [$SIBYL_PID] $SIBYL_CMDLINE
exit:
  - commandline: .*
    exec: echo [$SIBYL_PID] $SIBYL_CMDLINE
fork:
  - commandline: .*
    exec: echo [$SIBYL_PID] $SIBYL_CMDLINE
$ sudo sibyl sibyl.yaml
```

## Support Hook

* fork
* exec
* exit
* coredump
