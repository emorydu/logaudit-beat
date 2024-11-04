# DbAudit-Beat

1. internal/auditbeat (server)
2. internal/beatcli (linux)
3. internal/winbeatcli (windows)
4. internal/common (share pkg)



/var/gbk.log:::/var/gbk_utf8.log 1
客户端转编码定时任务只做一件事情 （存在数据的时候进行转编码并且更新该条数据的行号）