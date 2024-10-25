package model

type ConfigInfo struct {
	IP              string
	AgentPath       string
	MultiParse      int8   // 0 stop 1 startup
	RegexParamValue string // regexValue (parser_file)
	Check           int8   // rule open / stop
	ParseType       int8   // format (0 regex, 1 json)
	IndexName       string // topic = parse_file (name)
}

/*
===============================

collect_conf 采集配置表
srcIp String // 采集IP
sysType Int8	// 系统类型
encoding Int8	// 编码
agentPort Int32	// Agent通信端口
kafkaPort Int32	// Kafka通信端口
mapIp String	// ？？？
create_time DataTime

===============================

collect_conf_relation 采集配置关联关系表
srcIp String
id Int32
kafkaType Int8
agentPath String

===============================

log_index 日志索引表
id Int32
indexName String
fieldJson String
pjJson String
createTime DateTime
indexNameCN String
gid Int32

===============================

parsing_rule 规则解析表
rid Int32
id Int32
rname String
rdesc String
logSample String
parseType Int8
param1 String
param2 String
ignoreReg String
feature String
check Int8
secondary String
igEscape Int8
logSlice Int8
secondaryState Int8
parseType2 Int8
param1_2 String
param2_2 String
mapJson String
cType Int8
mutiParse Int8
*/
