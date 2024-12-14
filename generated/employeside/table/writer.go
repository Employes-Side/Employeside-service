//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/mysql"
)

var Writer = newWriterTable("employeside", "writer", "")

type writerTable struct {
	mysql.Table

	// Columns
	ID         mysql.ColumnString
	FirstName  mysql.ColumnString
	LastName   mysql.ColumnString
	Email      mysql.ColumnString
	UserName   mysql.ColumnString
	Password   mysql.ColumnString
	IsVerified mysql.ColumnBool
	IsActive   mysql.ColumnBool
	CreatedAt  mysql.ColumnTimestamp
	UpdatedAt  mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type WriterTable struct {
	writerTable

	NEW writerTable
}

// AS creates new WriterTable with assigned alias
func (a WriterTable) AS(alias string) *WriterTable {
	return newWriterTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WriterTable with assigned schema name
func (a WriterTable) FromSchema(schemaName string) *WriterTable {
	return newWriterTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WriterTable with assigned table prefix
func (a WriterTable) WithPrefix(prefix string) *WriterTable {
	return newWriterTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WriterTable with assigned table suffix
func (a WriterTable) WithSuffix(suffix string) *WriterTable {
	return newWriterTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWriterTable(schemaName, tableName, alias string) *WriterTable {
	return &WriterTable{
		writerTable: newWriterTableImpl(schemaName, tableName, alias),
		NEW:         newWriterTableImpl("", "new", ""),
	}
}

func newWriterTableImpl(schemaName, tableName, alias string) writerTable {
	var (
		IDColumn         = mysql.StringColumn("id")
		FirstNameColumn  = mysql.StringColumn("first_name")
		LastNameColumn   = mysql.StringColumn("last_name")
		EmailColumn      = mysql.StringColumn("email")
		UserNameColumn   = mysql.StringColumn("user_name")
		PasswordColumn   = mysql.StringColumn("password")
		IsVerifiedColumn = mysql.BoolColumn("is_Verified")
		IsActiveColumn   = mysql.BoolColumn("is_Active")
		CreatedAtColumn  = mysql.TimestampColumn("created_at")
		UpdatedAtColumn  = mysql.TimestampColumn("updated_at")
		allColumns       = mysql.ColumnList{IDColumn, FirstNameColumn, LastNameColumn, EmailColumn, UserNameColumn, PasswordColumn, IsVerifiedColumn, IsActiveColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns   = mysql.ColumnList{FirstNameColumn, LastNameColumn, EmailColumn, UserNameColumn, PasswordColumn, IsVerifiedColumn, IsActiveColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return writerTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		FirstName:  FirstNameColumn,
		LastName:   LastNameColumn,
		Email:      EmailColumn,
		UserName:   UserNameColumn,
		Password:   PasswordColumn,
		IsVerified: IsVerifiedColumn,
		IsActive:   IsActiveColumn,
		CreatedAt:  CreatedAtColumn,
		UpdatedAt:  UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
