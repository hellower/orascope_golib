package pgx

import (
	"context"
	"github.com/hellower/orascope_golib/orascopeLogger"
)

func (c *Conn) QueryExDebug(ctx context.Context, sql string, options *QueryExOptions, args ...interface{}) (rows *Rows, err error) {
	{
		orascopeLogger.Glog().Trace(sql)
	}
	return c.QueryEx(ctx, sql, options, args...)
}

func (c *Conn) QueryDebug(sql string, args ...interface{}) (*Rows, error) {
	return c.QueryExDebug(context.Background(), sql, nil, args...)
}

func (c *Conn) QueryRowDebug(sql string, args ...interface{}) *Row {
	rows, _ := c.QueryDebug(sql, args...)
	return (*Row)(rows)
}

func (c *Conn) QueryRowExDebug(ctx context.Context, sql string, options *QueryExOptions, args ...interface{}) *Row {
	rows, _ := c.QueryExDebug(ctx, sql, options, args...)
	return (*Row)(rows)
}

func (c *Conn) ExecDebug(sql string, arguments ...interface{}) (commandTag CommandTag, err error) {
	{
		orascopeLogger.Glog().Trace(sql)
	}
	return c.ExecEx(context.Background(), sql, nil, arguments...)
}

func (p *ConnPool) QueryRowByRole(a_role string, sql string, args ...interface{}) *Row {
	c, err := p.Acquire()
	if err != nil {
		// Because checking for errors can be deferred to the *Rows, build one with the error
		return (*Row)(&Rows{closed: true, err: err})
	}

	_, err = c.Exec("SET ROLE '" + a_role + "'")
	if err != nil {
		p.Release(c)
		return (*Row)(&Rows{closed: true, err: err})
	}

	rows, err := c.Query(sql, args...)
	if err != nil {
		p.Release(c)
		return (*Row)(rows)
	}

	rows.connPool = p
	return (*Row)(rows)
}

/*
func (p *ConnPool) QueryRow(sql string, args ...interface{}) *Row {
	rows, _ := p.Query(sql, args...)
	return (*Row)(rows)
}

func (p *ConnPool) Query(sql string, args ...interface{}) (*Rows, error) {
	c, err := p.Acquire()
	if err != nil {
		// Because checking for errors can be deferred to the *Rows, build one with the error
		return &Rows{closed: true, err: err}, err
	}

	rows, err := c.Query(sql, args...)
	if err != nil {
		p.Release(c)
		return rows, err
	}

	rows.connPool = p

	return rows, nil
}

*/
