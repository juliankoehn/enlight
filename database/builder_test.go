package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplode(t *testing.T) {
	columns := []string{"FOO", "BAR"}
	table := "table"
	prefix := "test_"
	typ := "index"

	imploded := strings.ToLower(prefix + table + "_" + strings.Join(columns, "_") + "_" + typ)
	t.Log(imploded)
}

func TestBuilder(t *testing.T) {
	manager := New()

	sql := ConnectionConfig{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Username: "testUser",
		Password: "yfLpFsBG2uMRhMaG",
		Database: "test",
		Prefix:   "test_",
	}

	manager.AddConnection(&sql, "mysql")
	conn, err := manager.GetConnection("mysql")
	if err != nil {
		t.Error(err)
	}

	builder := NewBuilder(*conn)
	builder.HasTable("test")

	// valid table
	exists, err := builder.HasColumn("test", "test")
	if err != nil {
		t.Error(err)
	}
	if !exists {
		t.Errorf("hasColumn should return true but got [%v]", exists)
	}

	exists, err = builder.HasColumn("test", "notExistent")
	if err != nil {
		t.Error(err)
	}
	if exists {
		t.Errorf("hasColumn should return false but got [%v]", exists)
	}

	result, err := builder.HasColumns("test", []string{"test", "age", "nope"})
	if err != nil {
		t.Error(err)
	}
	if result {
		t.Errorf("hasColumns should return false: 'nope' does not exists, got [true]")
	}

	resultStringSlice, err := builder.GetColumnListing("test")
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, resultStringSlice)

	// lets test the table creation process

	stmt := builder.Create("awesome_table", func(table *Blueprint) {
		table.UnsignedBigInteger("id", true)
		table.String("label", 255)
		table.Text("long_description")
		table.MediumText("short_description")
		table.TinyInteger("is_active", false, false)
	})

	t.Log(stmt)
}
