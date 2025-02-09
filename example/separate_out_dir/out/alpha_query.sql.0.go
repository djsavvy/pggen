// Code generated by pggen. DO NOT EDIT.

package out

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
	AlphaNested(ctx context.Context) (string, error)
	// AlphaNestedBatch enqueues a AlphaNested query into batch to be executed
	// later by the batch.
	AlphaNestedBatch(batch genericBatch)
	// AlphaNestedScan scans the result of an executed AlphaNestedBatch query.
	AlphaNestedScan(results pgx.BatchResults) (string, error)

	AlphaCompositeArray(ctx context.Context) ([]Alpha, error)
	// AlphaCompositeArrayBatch enqueues a AlphaCompositeArray query into batch to be executed
	// later by the batch.
	AlphaCompositeArrayBatch(batch genericBatch)
	// AlphaCompositeArrayScan scans the result of an executed AlphaCompositeArrayBatch query.
	AlphaCompositeArrayScan(results pgx.BatchResults) ([]Alpha, error)

	Alpha(ctx context.Context) (string, error)
	// AlphaBatch enqueues a Alpha query into batch to be executed
	// later by the batch.
	AlphaBatch(batch genericBatch)
	// AlphaScan scans the result of an executed AlphaBatch query.
	AlphaScan(results pgx.BatchResults) (string, error)

	Bravo(ctx context.Context) (string, error)
	// BravoBatch enqueues a Bravo query into batch to be executed
	// later by the batch.
	BravoBatch(batch genericBatch)
	// BravoScan scans the result of an executed BravoBatch query.
	BravoScan(results pgx.BatchResults) (string, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
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

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...interface{})
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, alphaNestedSQL, alphaNestedSQL); err != nil {
		return fmt.Errorf("prepare query 'AlphaNested': %w", err)
	}
	if _, err := p.Prepare(ctx, alphaCompositeArraySQL, alphaCompositeArraySQL); err != nil {
		return fmt.Errorf("prepare query 'AlphaCompositeArray': %w", err)
	}
	if _, err := p.Prepare(ctx, alphaSQL, alphaSQL); err != nil {
		return fmt.Errorf("prepare query 'Alpha': %w", err)
	}
	if _, err := p.Prepare(ctx, bravoSQL, bravoSQL); err != nil {
		return fmt.Errorf("prepare query 'Bravo': %w", err)
	}
	return nil
}

// Alpha represents the Postgres composite type "alpha".
type Alpha struct {
	Key *string `json:"key"`
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

type compositeField struct {
	name       string                 // name of the field
	typeName   string                 // Postgres type name
	defaultVal pgtype.ValueTranscoder // default value to use
}

func (tr *typeResolver) newCompositeValue(name string, fields ...compositeField) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	fs := make([]pgtype.CompositeTypeField, len(fields))
	vals := make([]pgtype.ValueTranscoder, len(fields))
	isBinaryOk := true
	for i, field := range fields {
		oid, val, ok := tr.findValue(field.typeName)
		if !ok {
			oid = unknownOID
			val = field.defaultVal
		}
		isBinaryOk = isBinaryOk && oid != unknownOID
		fs[i] = pgtype.CompositeTypeField{Name: field.name, OID: oid}
		vals[i] = val
	}
	// Okay to ignore error because it's only thrown when the number of field
	// names does not equal the number of ValueTranscoders.
	typ, _ := pgtype.NewCompositeTypeValues(name, fs, vals)
	if !isBinaryOk {
		return textPreferrer{typ, name}
	}
	return typ
}

func (tr *typeResolver) newArrayValue(name, elemName string, defaultVal func() pgtype.ValueTranscoder) pgtype.ValueTranscoder {
	if _, val, ok := tr.findValue(name); ok {
		return val
	}
	elemOID, elemVal, ok := tr.findValue(elemName)
	elemValFunc := func() pgtype.ValueTranscoder {
		return pgtype.NewValue(elemVal).(pgtype.ValueTranscoder)
	}
	if !ok {
		elemOID = unknownOID
		elemValFunc = defaultVal
	}
	typ := pgtype.NewArrayType(name, elemOID, elemValFunc)
	if elemOID == unknownOID {
		return textPreferrer{typ, name}
	}
	return typ
}

// newAlpha creates a new pgtype.ValueTranscoder for the Postgres
// composite type 'alpha'.
func (tr *typeResolver) newAlpha() pgtype.ValueTranscoder {
	return tr.newCompositeValue(
		"alpha",
		compositeField{"key", "text", &pgtype.Text{}},
	)
}

// newAlphaArray creates a new pgtype.ValueTranscoder for the Postgres
// '_alpha' array type.
func (tr *typeResolver) newAlphaArray() pgtype.ValueTranscoder {
	return tr.newArrayValue("_alpha", "alpha", tr.newAlpha)
}

const alphaNestedSQL = `SELECT 'alpha_nested' as output;`

// AlphaNested implements Querier.AlphaNested.
func (q *DBQuerier) AlphaNested(ctx context.Context) (string, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "AlphaNested")
	row := q.conn.QueryRow(ctx, alphaNestedSQL)
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query AlphaNested: %w", err)
	}
	return item, nil
}

// AlphaNestedBatch implements Querier.AlphaNestedBatch.
func (q *DBQuerier) AlphaNestedBatch(batch genericBatch) {
	batch.Queue(alphaNestedSQL)
}

// AlphaNestedScan implements Querier.AlphaNestedScan.
func (q *DBQuerier) AlphaNestedScan(results pgx.BatchResults) (string, error) {
	row := results.QueryRow()
	var item string
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan AlphaNestedBatch row: %w", err)
	}
	return item, nil
}

const alphaCompositeArraySQL = `SELECT ARRAY[ROW('key')]::alpha[];`

// AlphaCompositeArray implements Querier.AlphaCompositeArray.
func (q *DBQuerier) AlphaCompositeArray(ctx context.Context) ([]Alpha, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "AlphaCompositeArray")
	row := q.conn.QueryRow(ctx, alphaCompositeArraySQL)
	item := []Alpha{}
	arrayArray := q.types.newAlphaArray()
	if err := row.Scan(arrayArray); err != nil {
		return item, fmt.Errorf("query AlphaCompositeArray: %w", err)
	}
	if err := arrayArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign AlphaCompositeArray row: %w", err)
	}
	return item, nil
}

// AlphaCompositeArrayBatch implements Querier.AlphaCompositeArrayBatch.
func (q *DBQuerier) AlphaCompositeArrayBatch(batch genericBatch) {
	batch.Queue(alphaCompositeArraySQL)
}

// AlphaCompositeArrayScan implements Querier.AlphaCompositeArrayScan.
func (q *DBQuerier) AlphaCompositeArrayScan(results pgx.BatchResults) ([]Alpha, error) {
	row := results.QueryRow()
	item := []Alpha{}
	arrayArray := q.types.newAlphaArray()
	if err := row.Scan(arrayArray); err != nil {
		return item, fmt.Errorf("scan AlphaCompositeArrayBatch row: %w", err)
	}
	if err := arrayArray.AssignTo(&item); err != nil {
		return item, fmt.Errorf("assign AlphaCompositeArray row: %w", err)
	}
	return item, nil
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
