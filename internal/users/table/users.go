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

var Users = newUsersTable("users", "users", "")

type usersTable struct {
	mysql.Table

	// Columns
	ID        mysql.ColumnString
	UserName  mysql.ColumnString
	Email     mysql.ColumnString
	Password  mysql.ColumnString
	CreatedAt mysql.ColumnTimestamp
	UpdatedAt mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type UsersTable struct {
	usersTable

	NEW usersTable
}

// AS creates new UsersTable with assigned alias
func (a UsersTable) AS(alias string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UsersTable with assigned schema name
func (a UsersTable) FromSchema(schemaName string) *UsersTable {
	return newUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UsersTable with assigned table prefix
func (a UsersTable) WithPrefix(prefix string) *UsersTable {
	return newUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UsersTable with assigned table suffix
func (a UsersTable) WithSuffix(suffix string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUsersTable(schemaName, tableName, alias string) *UsersTable {
	return &UsersTable{
		usersTable: newUsersTableImpl(schemaName, tableName, alias),
		NEW:        newUsersTableImpl("", "new", ""),
	}
}

func newUsersTableImpl(schemaName, tableName, alias string) usersTable {
	var (
		IDColumn        = mysql.StringColumn("id")
		UserNameColumn  = mysql.StringColumn("user_name")
		EmailColumn     = mysql.StringColumn("email")
		PasswordColumn  = mysql.StringColumn("password")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		UpdatedAtColumn = mysql.TimestampColumn("updated_at")
		allColumns      = mysql.ColumnList{IDColumn, UserNameColumn, EmailColumn, PasswordColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = mysql.ColumnList{UserNameColumn, EmailColumn, PasswordColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return usersTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		UserName:  UserNameColumn,
		Email:     EmailColumn,
		Password:  PasswordColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
