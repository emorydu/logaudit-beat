package main

import (
	"fmt"
	"testing"
)

func TestAppendContent(t *testing.T) {
	content := AppendContent(`
[INPUT]
	Name tail
	Path /helloword, /test000
	Tag agent_json
	Read_From_Head true
	(insert) agent_json

[FILTER]
	Name parser
	Match agent_json
	Key_Name log
	Parser agent_json

[OUTPUT]
	Name kafka
	Match agent_json
	Brokers kafka:9200
	Topics data_agent_json

[INPUT]
	Name tail
	Path /helloword, /test000
	Tag agent_json
	Read_From_Head true
	(insert) agent_json

[FILTER]
	Name parser
	Match agent_json
	Key_Name log
	Parser agent_json

[OUTPUT]
	Name kafka
	Match agent_json
	Brokers kafka:9200
	Topics data_agent_json

`, "192.168.1.32")
	fmt.Println(string(content))
}
