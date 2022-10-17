// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/testifysec/archivist/ent/attestationcollection"
	"github.com/testifysec/archivist/ent/dsse"
	"github.com/testifysec/archivist/ent/predicate"
	"github.com/testifysec/archivist/ent/statement"
	"github.com/testifysec/archivist/ent/subject"
)

// StatementQuery is the builder for querying Statement entities.
type StatementQuery struct {
	config
	limit                      *int
	offset                     *int
	unique                     *bool
	order                      []OrderFunc
	fields                     []string
	predicates                 []predicate.Statement
	withSubjects               *SubjectQuery
	withAttestationCollections *AttestationCollectionQuery
	withDsse                   *DsseQuery
	modifiers                  []func(*sql.Selector)
	loadTotal                  []func(context.Context, []*Statement) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the StatementQuery builder.
func (sq *StatementQuery) Where(ps ...predicate.Statement) *StatementQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit adds a limit step to the query.
func (sq *StatementQuery) Limit(limit int) *StatementQuery {
	sq.limit = &limit
	return sq
}

// Offset adds an offset step to the query.
func (sq *StatementQuery) Offset(offset int) *StatementQuery {
	sq.offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *StatementQuery) Unique(unique bool) *StatementQuery {
	sq.unique = &unique
	return sq
}

// Order adds an order step to the query.
func (sq *StatementQuery) Order(o ...OrderFunc) *StatementQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// QuerySubjects chains the current query on the "subjects" edge.
func (sq *StatementQuery) QuerySubjects() *SubjectQuery {
	query := &SubjectQuery{config: sq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(statement.Table, statement.FieldID, selector),
			sqlgraph.To(subject.Table, subject.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, statement.SubjectsTable, statement.SubjectsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryAttestationCollections chains the current query on the "attestation_collections" edge.
func (sq *StatementQuery) QueryAttestationCollections() *AttestationCollectionQuery {
	query := &AttestationCollectionQuery{config: sq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(statement.Table, statement.FieldID, selector),
			sqlgraph.To(attestationcollection.Table, attestationcollection.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, statement.AttestationCollectionsTable, statement.AttestationCollectionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryDsse chains the current query on the "dsse" edge.
func (sq *StatementQuery) QueryDsse() *DsseQuery {
	query := &DsseQuery{config: sq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := sq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(statement.Table, statement.FieldID, selector),
			sqlgraph.To(dsse.Table, dsse.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, statement.DsseTable, statement.DsseColumn),
		)
		fromU = sqlgraph.SetNeighbors(sq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Statement entity from the query.
// Returns a *NotFoundError when no Statement was found.
func (sq *StatementQuery) First(ctx context.Context) (*Statement, error) {
	nodes, err := sq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{statement.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *StatementQuery) FirstX(ctx context.Context) *Statement {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Statement ID from the query.
// Returns a *NotFoundError when no Statement ID was found.
func (sq *StatementQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{statement.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *StatementQuery) FirstIDX(ctx context.Context) int {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Statement entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Statement entity is found.
// Returns a *NotFoundError when no Statement entities are found.
func (sq *StatementQuery) Only(ctx context.Context) (*Statement, error) {
	nodes, err := sq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{statement.Label}
	default:
		return nil, &NotSingularError{statement.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *StatementQuery) OnlyX(ctx context.Context) *Statement {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Statement ID in the query.
// Returns a *NotSingularError when more than one Statement ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *StatementQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = sq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{statement.Label}
	default:
		err = &NotSingularError{statement.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *StatementQuery) OnlyIDX(ctx context.Context) int {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Statements.
func (sq *StatementQuery) All(ctx context.Context) ([]*Statement, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return sq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (sq *StatementQuery) AllX(ctx context.Context) []*Statement {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Statement IDs.
func (sq *StatementQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := sq.Select(statement.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *StatementQuery) IDsX(ctx context.Context) []int {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *StatementQuery) Count(ctx context.Context) (int, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return sq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (sq *StatementQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *StatementQuery) Exist(ctx context.Context) (bool, error) {
	if err := sq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return sq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *StatementQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the StatementQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *StatementQuery) Clone() *StatementQuery {
	if sq == nil {
		return nil
	}
	return &StatementQuery{
		config:                     sq.config,
		limit:                      sq.limit,
		offset:                     sq.offset,
		order:                      append([]OrderFunc{}, sq.order...),
		predicates:                 append([]predicate.Statement{}, sq.predicates...),
		withSubjects:               sq.withSubjects.Clone(),
		withAttestationCollections: sq.withAttestationCollections.Clone(),
		withDsse:                   sq.withDsse.Clone(),
		// clone intermediate query.
		sql:    sq.sql.Clone(),
		path:   sq.path,
		unique: sq.unique,
	}
}

// WithSubjects tells the query-builder to eager-load the nodes that are connected to
// the "subjects" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *StatementQuery) WithSubjects(opts ...func(*SubjectQuery)) *StatementQuery {
	query := &SubjectQuery{config: sq.config}
	for _, opt := range opts {
		opt(query)
	}
	sq.withSubjects = query
	return sq
}

// WithAttestationCollections tells the query-builder to eager-load the nodes that are connected to
// the "attestation_collections" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *StatementQuery) WithAttestationCollections(opts ...func(*AttestationCollectionQuery)) *StatementQuery {
	query := &AttestationCollectionQuery{config: sq.config}
	for _, opt := range opts {
		opt(query)
	}
	sq.withAttestationCollections = query
	return sq
}

// WithDsse tells the query-builder to eager-load the nodes that are connected to
// the "dsse" edge. The optional arguments are used to configure the query builder of the edge.
func (sq *StatementQuery) WithDsse(opts ...func(*DsseQuery)) *StatementQuery {
	query := &DsseQuery{config: sq.config}
	for _, opt := range opts {
		opt(query)
	}
	sq.withDsse = query
	return sq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Predicate string `json:"predicate,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Statement.Query().
//		GroupBy(statement.FieldPredicate).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (sq *StatementQuery) GroupBy(field string, fields ...string) *StatementGroupBy {
	grbuild := &StatementGroupBy{config: sq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := sq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return sq.sqlQuery(ctx), nil
	}
	grbuild.label = statement.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Predicate string `json:"predicate,omitempty"`
//	}
//
//	client.Statement.Query().
//		Select(statement.FieldPredicate).
//		Scan(ctx, &v)
func (sq *StatementQuery) Select(fields ...string) *StatementSelect {
	sq.fields = append(sq.fields, fields...)
	selbuild := &StatementSelect{StatementQuery: sq}
	selbuild.label = statement.Label
	selbuild.flds, selbuild.scan = &sq.fields, selbuild.Scan
	return selbuild
}

func (sq *StatementQuery) prepareQuery(ctx context.Context) error {
	for _, f := range sq.fields {
		if !statement.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *StatementQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Statement, error) {
	var (
		nodes       = []*Statement{}
		_spec       = sq.querySpec()
		loadedTypes = [3]bool{
			sq.withSubjects != nil,
			sq.withAttestationCollections != nil,
			sq.withDsse != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Statement).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Statement{config: sq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(sq.modifiers) > 0 {
		_spec.Modifiers = sq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := sq.withSubjects; query != nil {
		if err := sq.loadSubjects(ctx, query, nodes,
			func(n *Statement) { n.Edges.Subjects = []*Subject{} },
			func(n *Statement, e *Subject) { n.Edges.Subjects = append(n.Edges.Subjects, e) }); err != nil {
			return nil, err
		}
	}
	if query := sq.withAttestationCollections; query != nil {
		if err := sq.loadAttestationCollections(ctx, query, nodes, nil,
			func(n *Statement, e *AttestationCollection) { n.Edges.AttestationCollections = e }); err != nil {
			return nil, err
		}
	}
	if query := sq.withDsse; query != nil {
		if err := sq.loadDsse(ctx, query, nodes,
			func(n *Statement) { n.Edges.Dsse = []*Dsse{} },
			func(n *Statement, e *Dsse) { n.Edges.Dsse = append(n.Edges.Dsse, e) }); err != nil {
			return nil, err
		}
	}
	for i := range sq.loadTotal {
		if err := sq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (sq *StatementQuery) loadSubjects(ctx context.Context, query *SubjectQuery, nodes []*Statement, init func(*Statement), assign func(*Statement, *Subject)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Statement)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Subject(func(s *sql.Selector) {
		s.Where(sql.InValues(statement.SubjectsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.statement_subjects
		if fk == nil {
			return fmt.Errorf(`foreign-key "statement_subjects" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "statement_subjects" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (sq *StatementQuery) loadAttestationCollections(ctx context.Context, query *AttestationCollectionQuery, nodes []*Statement, init func(*Statement), assign func(*Statement, *AttestationCollection)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Statement)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.AttestationCollection(func(s *sql.Selector) {
		s.Where(sql.InValues(statement.AttestationCollectionsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.statement_attestation_collections
		if fk == nil {
			return fmt.Errorf(`foreign-key "statement_attestation_collections" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "statement_attestation_collections" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (sq *StatementQuery) loadDsse(ctx context.Context, query *DsseQuery, nodes []*Statement, init func(*Statement), assign func(*Statement, *Dsse)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*Statement)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.InValues(statement.DsseColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.dsse_statement
		if fk == nil {
			return fmt.Errorf(`foreign-key "dsse_statement" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "dsse_statement" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (sq *StatementQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	if len(sq.modifiers) > 0 {
		_spec.Modifiers = sq.modifiers
	}
	_spec.Node.Columns = sq.fields
	if len(sq.fields) > 0 {
		_spec.Unique = sq.unique != nil && *sq.unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *StatementQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (sq *StatementQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   statement.Table,
			Columns: statement.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: statement.FieldID,
			},
		},
		From:   sq.sql,
		Unique: true,
	}
	if unique := sq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := sq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, statement.FieldID)
		for i := range fields {
			if fields[i] != statement.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *StatementQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(statement.Table)
	columns := sq.fields
	if len(columns) == 0 {
		columns = statement.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.unique != nil && *sq.unique {
		selector.Distinct()
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// StatementGroupBy is the group-by builder for Statement entities.
type StatementGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *StatementGroupBy) Aggregate(fns ...AggregateFunc) *StatementGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the group-by query and scans the result into the given value.
func (sgb *StatementGroupBy) Scan(ctx context.Context, v any) error {
	query, err := sgb.path(ctx)
	if err != nil {
		return err
	}
	sgb.sql = query
	return sgb.sqlScan(ctx, v)
}

func (sgb *StatementGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range sgb.fields {
		if !statement.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := sgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (sgb *StatementGroupBy) sqlQuery() *sql.Selector {
	selector := sgb.sql.Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(sgb.fields)+len(sgb.fns))
		for _, f := range sgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(sgb.fields...)...)
}

// StatementSelect is the builder for selecting fields of Statement entities.
type StatementSelect struct {
	*StatementQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ss *StatementSelect) Scan(ctx context.Context, v any) error {
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	ss.sql = ss.StatementQuery.sqlQuery(ctx)
	return ss.sqlScan(ctx, v)
}

func (ss *StatementSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ss.sql.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
