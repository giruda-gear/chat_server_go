package repository

import (
	"chat_controller_go/config"
	"chat_controller_go/repository/kafka"
	"chat_controller_go/types/table"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config

	db *sql.DB

	Kafka *kafka.Kafka
}

const (
	room       = "chatting.room"
	chat       = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}
	var err error

	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	} else if r.Kafka, err = kafka.NewKafka(cfg); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (r *Repository) GetAvailableServers() ([]*table.ServerInfo, error) {
	qs := query([]string{"SELECT * FROM", serverInfo, "WHERE available = 1"})

	rows, err := r.db.Query(qs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []*table.ServerInfo

	for rows.Next() {
		d := new(table.ServerInfo)

		if err = rows.Scan(
			&d.IP,
			&d.Available,
		); err != nil {
			return nil, err
		}
		servers = append(servers, d)
	}

	if len(servers) == 0 {
		return []*table.ServerInfo{}, nil
	}
	return servers, rows.Err()
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
