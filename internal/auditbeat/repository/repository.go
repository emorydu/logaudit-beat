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
	InsertOrUpdateMonitor(ctx context.Context, ip string, cpuUse, memUse float64, status, updated int) error
	QueryMonitorInfo(ctx context.Context, ip string) (int, error)
	Update(context.Context, string) error
}

type repository struct {
	db driver.Conn
}

func NewRepository(orm driver.Conn) Repository {
	return &repository{
		db: orm,
	}
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

func (r *repository) InsertOrUpdateMonitor(ctx context.Context, ip string, cpuUse, memUse float64, status, _ int) error {
	q := "SELECT ip FROM monitor WHERE ip = ?"
	err := r.db.QueryRow(ctx, q, ip).Scan(&ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			q = "INSERT INTO monitor (ip) VALUES (?)"
			_ = r.db.Exec(ctx, q, ip)
		}
	} else {

		q = "ALTER TABLE monitor UPDATE cpuUse = ?, memUse = ?, status = ? WHERE ip = ?"
		err = r.db.Exec(ctx, q, cpuUse, memUse, status, ip)

	}
	return nil
}

func (r *repository) FetchConfInfo(ctx context.Context, ip string) ([]model.ConfigInfo, error) {
	items := make([]model.ConfigInfo, 0)
	rows, err := r.db.Query(ctx, `
SELECT ccr.srcIp,
       ccr.agentPath,
       pr.mutiParse,
       pr.param1 AS regexValue,
       pr.check,
       pr.parseType,
       pr.indexName
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
		err := rows.Scan(&item.IP, &item.AgentPath, &item.MultiParse, &item.RegexParamValue, &item.Check, &item.ParseType, &item.IndexName)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
