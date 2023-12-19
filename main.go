package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func main() {
	conn, err := connect()
	if err != nil {
		panic((err))
	}
	ctx := context.Background()

	err = BatchInsert(conn, ctx)
	if err != nil {
		return
	}

	//rows, err := conn.Query(ctx, "SELECT name,toString(uuid) as uuid_str FROM system.tables LIMIT 500")
	rows, err := conn.Query(ctx, "SELECT Col3,Col4,Col5 FROM example LIMIT 500")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(rows.Columns()))
	for rows.Next() {
		var (
			name, uuid string
		)

		if err := rows.Scan(
			&name,
			&uuid,
		); err != nil {
			log.Fatal(err)
		}
		log.Printf("name: %s, uuid: %s",
			name, uuid)
	}

}

func connect() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"localhost:9000"},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: "default",
				Password: "default",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func BatchInsert(conn driver.Conn, ctx context.Context) error {
	conn.Exec(ctx, "DROP TABLE IF EXISTS example")
	err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS example (
			  Col1 UInt8
			, Col2 String
			, Col3 FixedString(3)
			, Col4 UUID
			, Col5 Map(String, UInt8)
			, Col6 Array(String)
			, Col7 Tuple(String, UInt8, Array(Map(String, String)))
			, Col8 DateTime
		) Engine = Memory
	`)
	if err != nil {
		return err
	}

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO example")
	if err != nil {
		return err
	}
	for i := 0; i < 1000; i++ {
		err := batch.Append(
			uint8(42),
			"ClickHouse",
			"Inc",
			uuid.New(),
			map[string]uint8{"key": 1},             // Map(String, UInt8)
			[]string{"Q", "W", "E", "R", "T", "Y"}, // Array(String)
			[]any{ // Tuple(String, UInt8, Array(Map(String, String)))
				"String Value", uint8(5), []map[string]string{
					{"key": "value"},
					{"key": "value"},
					{"key": "value"},
				},
			},
			time.Now(),
		)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}
