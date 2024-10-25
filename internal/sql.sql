SELECT ccr.srcIp,
       ccr.id,
       pr.param1, -- 正则表达式 是parsing中的Regex值
       pr.ignoreReg, -- 是否去除头信息
        pr.feature, --
        pr.check, -- 规则的启用或者停用
        pr.igEscape,-- 是否去除转义
        pr.parseType, -- 解析方式 parserFile=Format 0 正则，2 JSON
        pr.logSlice, -- 去除首尾字节
        li.indexName -- topic == parserFile=name
FROM collect_conf AS ccf
         RIGHT JOIN collect_conf_relation AS ccr ON ccf.srcIp = ccr.srcIp
         LEFT JOIN parsing_rule AS pr ON ccr.id = pr.rid  LEFT JOIN log_index AS li ON pr.id = li.id
WHERE ccr.agentPath != '';




[parser]
 Name: xxx
 Format: xxx (regex / json)
 Regex: xxx (只有Format为regex时才有，)
[parser]
 Name: xxx
 Format: xxx
 Regex: xxx
[parser]
 Name: xxx
 Format: xxx
 Regex: xxx





===========================




[INPUT]
    name (tail) // fixed
    Path (xxx) agentPath
    tag (indexName)


[FILTER]
    Name (parser) // fixed
    Match (indexName)
    Key_Name (log)
    Parser (xxx) indexName

[OUTPUT]
    name (kafka)
    match (indexName)
    Brokers (目标IP+端口号)
    Topics (data_indexName)

clickhouse-client -u default --password Safe.app