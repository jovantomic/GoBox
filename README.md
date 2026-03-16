# GoBox

container runtime from scratch in Go. Why? It sounds fun and wanted to learn go

basically building what docker does under the hood — namespaces, cgroups, overlayfs, networking.

## status

just started. can spawn a child process, next up is actual isolation.

```bash
sudo ./gobox run /bin/sh
```

## todo

- [x] basic process execution
- [ ] namespace isolation (PID, UTS, MNT, NET)
- [ ] cgroups v2 (mem/cpu limits)
- [ ] overlayfs + alpine rootfs
- [ ] bridge networking
- [ ] cleanup

## dev log

[diary.md](diary.md)
