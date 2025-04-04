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

var Modules = newModulesTable("employeside", "modules", "")

type modulesTable struct {
	mysql.Table

	// Columns
	ID              mysql.ColumnString
	UserID          mysql.ColumnString
	ModuleName      mysql.ColumnString
	ModuleType      mysql.ColumnString
	ModuleDesc      mysql.ColumnString
	ModuleShortName mysql.ColumnString
	ModulePrice     mysql.ColumnString
	Purchased       mysql.ColumnBool
	CreatedAt       mysql.ColumnTimestamp
	UpdatedAt       mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type ModulesTable struct {
	modulesTable

	NEW modulesTable
}

// AS creates new ModulesTable with assigned alias
func (a ModulesTable) AS(alias string) *ModulesTable {
	return newModulesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ModulesTable with assigned schema name
func (a ModulesTable) FromSchema(schemaName string) *ModulesTable {
	return newModulesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ModulesTable with assigned table prefix
func (a ModulesTable) WithPrefix(prefix string) *ModulesTable {
	return newModulesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ModulesTable with assigned table suffix
func (a ModulesTable) WithSuffix(suffix string) *ModulesTable {
	return newModulesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newModulesTable(schemaName, tableName, alias string) *ModulesTable {
	return &ModulesTable{
		modulesTable: newModulesTableImpl(schemaName, tableName, alias),
		NEW:          newModulesTableImpl("", "new", ""),
	}
}

func newModulesTableImpl(schemaName, tableName, alias string) modulesTable {
	var (
		IDColumn              = mysql.StringColumn("id")
		UserIDColumn          = mysql.StringColumn("user_id")
		ModuleNameColumn      = mysql.StringColumn("module_name")
		ModuleTypeColumn      = mysql.StringColumn("module_type")
		ModuleDescColumn      = mysql.StringColumn("module_desc")
		ModuleShortNameColumn = mysql.StringColumn("module_short_name")
		ModulePriceColumn     = mysql.StringColumn("module_price")
		PurchasedColumn       = mysql.BoolColumn("purchased")
		CreatedAtColumn       = mysql.TimestampColumn("created_at")
		UpdatedAtColumn       = mysql.TimestampColumn("updated_at")
		allColumns            = mysql.ColumnList{IDColumn, UserIDColumn, ModuleNameColumn, ModuleTypeColumn, ModuleDescColumn, ModuleShortNameColumn, ModulePriceColumn, PurchasedColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns        = mysql.ColumnList{UserIDColumn, ModuleNameColumn, ModuleTypeColumn, ModuleDescColumn, ModuleShortNameColumn, ModulePriceColumn, PurchasedColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return modulesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:              IDColumn,
		UserID:          UserIDColumn,
		ModuleName:      ModuleNameColumn,
		ModuleType:      ModuleTypeColumn,
		ModuleDesc:      ModuleDescColumn,
		ModuleShortName: ModuleShortNameColumn,
		ModulePrice:     ModulePriceColumn,
		Purchased:       PurchasedColumn,
		CreatedAt:       CreatedAtColumn,
		UpdatedAt:       UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
