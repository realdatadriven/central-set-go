package flight

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"google.golang.org/grpc"

	"github.com/hugr-lab/airport-go"
	"github.com/hugr-lab/airport-go/catalog"
)

// ─────────────────────────────────────────────────────────────────────────────
// DynamicTable — the only table type used everywhere
// ─────────────────────────────────────────────────────────────────────────────

type DynamicTable struct {
	mu       sync.RWMutex
	name     string
	schema   *arrow.Schema
	comment  string
	scanFunc catalog.ScanFunc // nil → empty table (DDL-created)
}

func NewDynamicTable(name string, schema *arrow.Schema, comment string, scanFunc catalog.ScanFunc) *DynamicTable {
	return &DynamicTable{
		name:     name,
		schema:   schema,
		comment:  comment,
		scanFunc: scanFunc,
	}
}

func (t *DynamicTable) Name() string    { return t.name }
func (t *DynamicTable) Comment() string { return t.comment }

func (t *DynamicTable) ArrowSchema(columns []string) *arrow.Schema {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return catalog.ProjectSchema(t.schema, columns) // assuming helper exists; otherwise implement projection
}

func (t *DynamicTable) Scan(ctx context.Context, opts *catalog.ScanOptions) (array.RecordReader, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.scanFunc != nil {
		return t.scanFunc(ctx, opts)
	}

	// DDL-created table → returns empty
	return array.NewRecordReader(t.schema, nil)
}

// DML stubs
func (t *DynamicTable) Insert(_ context.Context, _ array.RecordReader, _ *catalog.DMLOptions) (*catalog.DMLResult, error) {
	return nil, fmt.Errorf("INSERT not supported yet")
}

func (t *DynamicTable) Update(_ context.Context, _ []int64, _ array.RecordReader, _ *catalog.DMLOptions) (*catalog.DMLResult, error) {
	return nil, fmt.Errorf("UPDATE not supported yet")
}

func (t *DynamicTable) Delete(_ context.Context, _ []int64, _ *catalog.DMLOptions) (*catalog.DMLResult, error) {
	return nil, fmt.Errorf("DELETE not supported yet")
}

// Optional: column DDL stubs
func (t *DynamicTable) AddColumn(_ context.Context, _ *arrow.Field, _ catalog.AddColumnOptions) error {
	return fmt.Errorf("ADD COLUMN not supported yet")
}

// ─────────────────────────────────────────────────────────────────────────────
// Dynamic Catalog & Schema
// ─────────────────────────────────────────────────────────────────────────────

type DynamicCatalog struct {
	mu      sync.RWMutex
	schemas map[string]*DynamicSchema
}

func NewDynamicCatalog() *DynamicCatalog {
	return &DynamicCatalog{
		schemas: make(map[string]*DynamicSchema),
	}
}

func (c *DynamicCatalog) Schemas(ctx context.Context) ([]catalog.Schema, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]catalog.Schema, 0, len(c.schemas))
	for _, s := range c.schemas {
		res = append(res, s)
	}
	return res, nil
}

func (c *DynamicCatalog) Schema(ctx context.Context, name string) (catalog.Schema, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.schemas[name]
	if !ok {
		return nil, nil // airport-go often returns nil instead of error for "not found"
	}
	return s, nil
}

func (c *DynamicCatalog) CreateSchema(_ context.Context, name string, opts catalog.CreateSchemaOptions) (catalog.Schema, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.schemas[name]; exists {
		return nil, catalog.ErrAlreadyExists
	}
	s := NewDynamicSchema(name, opts.Comment)
	c.schemas[name] = s
	return s, nil
}

func (c *DynamicCatalog) DropSchema(_ context.Context, name string, opts catalog.DropSchemaOptions) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.schemas[name]; !ok {
		if opts.IgnoreNotFound {
			return nil
		}
		return catalog.ErrNotFound
	}
	delete(c.schemas, name)
	return nil
}

type DynamicSchema struct {
	mu      sync.RWMutex
	name    string
	comment string
	tables  map[string]*DynamicTable
}

func NewDynamicSchema(name, comment string) *DynamicSchema {
	return &DynamicSchema{
		name:    name,
		comment: comment,
		tables:  make(map[string]*DynamicTable),
	}
}

func (s *DynamicSchema) Name() string    { return s.name }
func (s *DynamicSchema) Comment() string { return s.comment }

func (s *DynamicSchema) Tables(ctx context.Context) ([]catalog.Table, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]catalog.Table, 0, len(s.tables))
	for _, t := range s.tables {
		res = append(res, t)
	}
	return res, nil
}

func (s *DynamicSchema) Table(ctx context.Context, name string) (catalog.Table, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tables[name]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (s *DynamicSchema) CreateTable(_ context.Context, name string, schema *arrow.Schema, opts catalog.CreateTableOptions) (catalog.Table, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, exists := s.tables[name]; exists {
		switch opts.OnConflict {
		case catalog.OnConflictIgnore:
			return existing, nil
		case catalog.OnConflictReplace:
			// fall through to replace
		default:
			return nil, catalog.ErrAlreadyExists
		}
	}

	t := NewDynamicTable(name, schema, opts.Comment, nil) // nil scanFunc → empty
	s.tables[name] = t
	return t, nil
}

func (s *DynamicSchema) DropTable(_ context.Context, name string, opts catalog.DropTableOptions) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tables[name]; !ok {
		if opts.IgnoreNotFound {
			return nil
		}
		return catalog.ErrNotFound
	}
	delete(s.tables, name)
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Builder (only uses DynamicTable)
// ─────────────────────────────────────────────────────────────────────────────

type CatalogBuilder struct {
	schemas []*schemaBuilder
	dynamic bool
	built   bool
}

func NewCatalogBuilder() *CatalogBuilder {
	return &CatalogBuilder{
		schemas: make([]*schemaBuilder, 0),
		dynamic: false,
		built:   false,
	}
}

func (cb *CatalogBuilder) Dynamic() *CatalogBuilder {
	cb.dynamic = true
	return cb
}

func (cb *CatalogBuilder) Schema(name string) *SchemaBuilder {
	sb := &schemaBuilder{
		name:           name,
		comment:        "",
		tables:         make([]*DynamicTable, 0),
		catalogBuilder: cb,
	}
	cb.schemas = append(cb.schemas, sb)
	return &SchemaBuilder{builder: sb}
}

func (cb *CatalogBuilder) Build() (catalog.Catalog, error) {
	if cb.built {
		return nil, fmt.Errorf("catalog already built")
	}
	cb.built = true

	// Basic validation
	seenSchemas := make(map[string]bool)
	for _, sb := range cb.schemas {
		if sb.name == "" || seenSchemas[sb.name] {
			return nil, fmt.Errorf("invalid/duplicate schema name: %q", sb.name)
		}
		seenSchemas[sb.name] = true

		seenTables := make(map[string]bool)
		for _, t := range sb.tables {
			if t.Name() == "" || seenTables[t.Name()] {
				return nil, fmt.Errorf("invalid/duplicate table %q in schema %q", t.Name(), sb.name)
			}
			seenTables[t.Name()] = true
			if t.schema == nil {
				return nil, fmt.Errorf("table %s.%s has nil schema", sb.name, t.Name())
			}
		}
	}

	if !cb.dynamic {
		// You could still support static path if desired — omitted here for simplicity
		return nil, fmt.Errorf("only dynamic mode is implemented in this example")
	}

	// Dynamic catalog with pre-loaded tables
	dcat := NewDynamicCatalog()

	for _, sb := range cb.schemas {
		s := &DynamicSchema{
			name:    sb.name,
			comment: sb.comment,
			tables:  make(map[string]*DynamicTable),
		}

		for _, t := range sb.tables {
			s.tables[t.Name()] = t
		}

		dcat.schemas[sb.name] = s
	}

	return dcat, nil
}

type schemaBuilder struct {
	name           string
	comment        string
	tables         []*DynamicTable
	catalogBuilder *CatalogBuilder
}

type SchemaBuilder struct {
	builder *schemaBuilder
}

func (sb *SchemaBuilder) Comment(c string) *SchemaBuilder {
	sb.builder.comment = c
	return sb
}

func (sb *SchemaBuilder) Table(t *DynamicTable) *SchemaBuilder {
	sb.builder.tables = append(sb.builder.tables, t)
	return sb
}

func (sb *SchemaBuilder) Build() (catalog.Catalog, error) {
	return sb.builder.catalogBuilder.Build()
}

/*
// ─────────────────────────────────────────────────────────────────────────────
// Example static data function
// ─────────────────────────────────────────────────────────────────────────────

var pool = memory.NewGoAllocator()

func exampleUsersScan(_ context.Context, _ *catalog.ScanOptions) (array.RecordReader, error) {
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
		nil,
	)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	b.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	b.Field(1).(*array.StringBuilder).AppendValues([]string{"Alice", "Bob", "Charlie"}, nil)

	rec := b.NewRecord()
	defer rec.Release()

	return array.NewSingleRecordReader(rec), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Main — how to use it
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	usersTable := NewDynamicTable(
		"users",
		arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: false},
				{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
			},
			nil,
		),
		"Pre-loaded example users",
		exampleUsersScan,
	)

	cat, err := NewCatalogBuilder().
		Dynamic().
		Schema("public").
		Comment("Main application schema").
		Table(usersTable).
		Schema("logs").
		Comment("Application logs").
		Build()

	if err != nil {
		log.Fatalf("Failed to build catalog: %v", err)
	}

	// Start Airport server
	debug := slog.LevelDebug
	config := airport.ServerConfig{
		Catalog:  cat,
		LogLevel: &debug,
	}

	srv := grpc.NewServer(airport.ServerOptions(config)...)
	if err := airport.NewServer(srv, config); err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	log.Println("Airport server listening on :50051")
	log.Println("Pre-loaded: public.users (has data)")
	log.Println("You can now:")
	log.Println("  CREATE SCHEMA demo;")
	log.Println("  CREATE TABLE demo.events (ts TIMESTAMP, msg VARCHAR);")
	log.Println("  SELECT * FROM public.users;          → returns 3 rows")
	log.Println("  SELECT * FROM demo.events;           → empty")
	log.Println("  INSERT INTO demo.events ...          → not supported yet")

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Serve failed: %v", err)
	}
}
*/