// Code generated by pggen. DO NOT EDIT.

package enums

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error)
	// FindAllDevicesBatch enqueues a FindAllDevices query into batch to be executed
	// later by the batch.
	FindAllDevicesBatch(ctx context.Context, batch *pgx.Batch)
	// FindAllDevicesScan scans the result of an executed FindAllDevicesBatch query.
	FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error)

	InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error)
	// InsertDeviceBatch enqueues a InsertDevice query into batch to be executed
	// later by the batch.
	InsertDeviceBatch(ctx context.Context, batch *pgx.Batch, mac pgtype.Macaddr, typePg DeviceType)
	// InsertDeviceScan scans the result of an executed InsertDeviceBatch query.
	InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error)
}

type DBQuerier struct {
	conn genericConn
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// DeviceType represents the Postgres enum device_type.
type DeviceType string

const (
	Undefined DeviceType = "undefined"
	Phone     DeviceType = "phone"
	Laptop    DeviceType = "laptop"
	Ipad      DeviceType = "ipad"
	Desktop   DeviceType = "desktop"
	Iot       DeviceType = "iot"
)

func (d DeviceType) String() string { return string(d) }

const findAllDevicesSQL = `SELECT mac, type from device;`

type FindAllDevicesRow struct {
	Mac  pgtype.Macaddr `json:"mac"`
	Type DeviceType     `json:"type"`
}

// FindAllDevices implements Querier.FindAllDevices.
func (q *DBQuerier) FindAllDevices(ctx context.Context) ([]FindAllDevicesRow, error) {
	rows, err := q.conn.Query(ctx, findAllDevicesSQL)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindAllDevices: %w", err)
	}
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevices row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindAllDevicesBatch implements Querier.FindAllDevicesBatch.
func (q *DBQuerier) FindAllDevicesBatch(ctx context.Context, batch *pgx.Batch) {
	batch.Queue(findAllDevicesSQL)
}

// FindAllDevicesScan implements Querier.FindAllDevicesScan.
func (q *DBQuerier) FindAllDevicesScan(results pgx.BatchResults) ([]FindAllDevicesRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	items := []FindAllDevicesRow{}
	for rows.Next() {
		var item FindAllDevicesRow
		if err := rows.Scan(&item.Mac, &item.Type); err != nil {
			return nil, fmt.Errorf("scan FindAllDevicesBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const insertDeviceSQL = `INSERT INTO device (mac, type) VALUES ($1, $2);`

// InsertDevice implements Querier.InsertDevice.
func (q *DBQuerier) InsertDevice(ctx context.Context, mac pgtype.Macaddr, typePg DeviceType) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, insertDeviceSQL, mac, typePg)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertDevice: %w", err)
	}
	return cmdTag, err
}

// InsertDeviceBatch implements Querier.InsertDeviceBatch.
func (q *DBQuerier) InsertDeviceBatch(ctx context.Context, batch *pgx.Batch, mac pgtype.Macaddr, typePg DeviceType) {
	batch.Queue(insertDeviceSQL, mac, typePg)
}

// InsertDeviceScan implements Querier.InsertDeviceScan.
func (q *DBQuerier) InsertDeviceScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertDeviceBatch: %w", err)
	}
	return cmdTag, err
}
