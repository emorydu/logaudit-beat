// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/emorydu/dbaudit/internal/auditbeat/model"
)

type Repository interface {
	FetchConfInfo(context.Context, string) ([]model.ConfigInfo, error)
	InsertOrUpdateMonitor(ctx context.Context, ip string, cpuUse, memUse float64, status int, timestamp int64) error
	QueryMonitorInfo(ctx context.Context, ip string) (int, error)
	Update(context.Context, string) error
	QueryMonitorTimestamp(ctx context.Context) (map[string]int64, error)
	UpdateStatus(ctx context.Context, ip string, status int) error
	QueryCollectConfig(ctx context.Context, ip string) (model.CollectInfo, error)
}

type repository struct {
	db driver.Conn
}

func NewRepository(orm driver.Conn) Repository {
	return &repository{
		db: orm,
	}
}

func (r *repository) QueryCollectConfig(ctx context.Context, ip string) (model.CollectInfo, error) {
	// TODO
	//q := "SELECT mapStatus, mapIp, kafkaPort FROM collect_conf WHERE ip = ?"
	//var collectConf model.CollectInfo
	//err := r.db.QueryRow(ctx, q, ip).Scan(&collectConf.MappingStatus, &collectConf.MappingIP, &collectConf.KafkaPort)
	//if err != nil {
	//	return model.CollectInfo{}, err
	//}
	//if collectConf.MappingStatus != 1 {
	//	q1 := `SELECT param_value FROM param_config WHERE param_id = 2001`
	//	q2 := `SELECT * FROM net_config`
	//}

	return model.CollectInfo{}, nil
}

func (r *repository) UpdateStatus(ctx context.Context, ip string, status int) error {
	q := "ALTER TABLE monitor UPDATE status = 2, cpuUse = 0, memUse = 0 WHERE ip = ?"
	return r.db.Exec(ctx, q, ip)
}

func (r *repository) QueryMonitorTimestamp(ctx context.Context) (map[string]int64, error) {
	q := "SELECT ip, timestamp FROM monitor WHERE operator != 2;"

	data := make(map[string]int64)
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			ip        string
			timestamp int64
		)
		err = rows.Scan(&ip, &timestamp)
		if err != nil {
			return nil, err
		}

		data[ip] = timestamp
	}

	return data, nil
}
func (r *repository) Update(ctx context.Context, ip string) error {
	q := "ALTER TABLE monitor UPDATE operator = 0 WHERE ip = ?"
	return r.db.Exec(ctx, q, ip)
}

func (r *repository) QueryMonitorInfo(ctx context.Context, ip string) (int, error) {
	q := `
SELECT operator FROM monitor WHERE ip = ?;
`
	var operator uint8
	err := r.db.QueryRow(ctx, q, ip).Scan(&operator)
	if err != nil {
		return -1, err
	}

	return int(operator), nil
}

func (r *repository) InsertOrUpdateMonitor(ctx context.Context, ip string, cpuUse, memUse float64, status int, timestamp int64) error {
	q := "SELECT ip FROM monitor WHERE ip = ?"
	err := r.db.QueryRow(ctx, q, ip).Scan(&ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			q = "INSERT INTO monitor (ip) VALUES (?)"
			return r.db.Exec(ctx, q, ip)
		}
	} else {

		q = "ALTER TABLE monitor UPDATE cpuUse = ?, memUse = ?, status = ?, timestamp = ? WHERE ip = ?"
		return r.db.Exec(ctx, q, cpuUse, memUse, status, timestamp, ip)

	}
	return nil
}

func (r *repository) FetchConfInfo(ctx context.Context, ip string) ([]model.ConfigInfo, error) {
	items := make([]model.ConfigInfo, 0)
	rows, err := r.db.Query(ctx, `
SELECT ccr.srcIp,
       ccf.mapStatus,
       ccf.mapIp,
       ccf.kafkaPort,
       ccr.agentPath,
       pr.mutiParse,
       pr.param1 AS regexValue,
       pr.check,
       pr.parseType,
       pr.indexName,
       pr.secondary,
       pr.secondaryState,
       pr.parseType2,
       pr.param1_2,
       pr.rid,
       ccr.encoding
FROM collect_conf AS ccf
         RIGHT JOIN collect_conf_relation AS ccr ON ccf.srcIp = ccr.srcIp
         LEFT JOIN parsing_rule AS pr ON ccr.rid = pr.rid 
WHERE ccr.agentPath != '' AND ccr.srcIp = ?;
`, ip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.ConfigInfo{}
		err := rows.Scan(&item.IP, &item.MappingStatus, &item.MappingIP, &item.KafkaPort, &item.AgentPath, &item.MultiParse, &item.RegexParamValue, &item.Check, &item.ParseType, &item.IndexName, &item.Secondary, &item.SecondaryState, &item.SecondaryParsingType, &item.SecondaryRegexValue, &item.RID, &item.Encoding)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
