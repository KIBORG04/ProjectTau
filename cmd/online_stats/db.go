package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

var dialerCounter atomic.Int64

func FetchNewRoundsFromDB(cfg *OnlineStatsConfig, lastID int) ([]Round, error) {
	dbCfg := &cfg.RemoteDB

	// 1. Establish SSH Client
	sshConfig := &ssh.ClientConfig{
		User: dbCfg.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(dbCfg.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	sshAddr := fmt.Sprintf("%s:%d", dbCfg.SSHHost, dbCfg.SSHPort)
	log.Printf("Connecting to SSH %s...\n", sshAddr)
	sshClient, err := ssh.Dial("tcp", sshAddr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh: %w", err)
	}
	defer sshClient.Close()

	// 2. Register custom MySQL dialer to route through SSH (unique name to avoid panic on re-registration)
	networkName := fmt.Sprintf("mysql+ssh+%s+%d", sshAddr, dialerCounter.Add(1))
	mysql.RegisterDialContext(networkName, func(ctx context.Context, addr string) (net.Conn, error) {
		return sshClient.Dial("tcp", addr)
	})
	defer mysql.DeregisterDialContext(networkName)

	// 3. Connect to MySQL via the custom dialer
	mysqlAddr := fmt.Sprintf("%s:%d", dbCfg.Host, dbCfg.Port)
	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true",
		dbCfg.User, dbCfg.Password, networkName, mysqlAddr, dbCfg.DBName)

	log.Printf("Connecting to MySQL via SSH tunnel...\n")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	// 4. Execute query
	query := `
SELECT
    id,
    server_port,
    start_datetime,
    end_datetime,
    players
FROM (
    SELECT
        r.id,
        r.server_port,
        r.start_datetime,
        r.end_datetime,
        f.var_value AS players,
        ROW_NUMBER() OVER (PARTITION BY r.id ORDER BY f.id DESC) AS rn
    FROM erro_round r
    JOIN erro_feedback f ON r.id = f.round_id
    WHERE f.var_name = 'round_end_clients'
      AND r.start_datetime IS NOT NULL
      AND r.end_datetime IS NOT NULL
      AND f.var_value > 0
      AND r.id > ?
) t
WHERE rn = 1
ORDER BY id ASC;`

	log.Printf("Fetching rounds > %d...\n", lastID)
	
	// Timeout for query execution
	queryCtx, queryCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer queryCancel()

	rows, err := db.QueryContext(queryCtx, query, lastID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var newRounds []Round
	for rows.Next() {
		var r Round
		var playersStr string
		var start, end sql.NullTime

		if err := rows.Scan(&r.ID, &r.ServerPort, &start, &end, &playersStr); err != nil {
			log.Printf("Warning: failed to scan row: %v\n", err)
			continue
		}

		if !start.Valid || !end.Valid {
			// Ignore null dates
			continue
		}
		
		r.StartDatetime = start.Time
		r.EndDatetime = end.Time

		fmt.Sscanf(playersStr, "%d", &r.Players)
		
		if r.Players > 0 {
			newRounds = append(newRounds, r)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Printf("Fetched %d new rounds from DB.\n", len(newRounds))
	return newRounds, nil
}
