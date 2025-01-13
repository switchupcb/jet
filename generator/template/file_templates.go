package template

var autoGenWarningTemplate = `
//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior  
// and will be lost if the code is regenerated
//

`

var tableSQLBuilderTemplate = ` 
package {{package}}

import (
	"github.com/go-jet/jet/v2/{{dialect.PackageName}}"
)

var {{tableTemplate.InstanceName}} = new{{tableTemplate.TypeName}}("{{schemaName}}", "{{.Name}}", "{{tableTemplate.DefaultAlias}}")

{{golangComment .Comment}}
type {{structImplName}} struct {
	{{dialect.PackageName}}.Table
	
	// Columns
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
{{- if not $field.Skip}}
	{{$field.Name}} {{dialect.PackageName}}.Column{{$field.Type}} {{golangComment .Comment}}
{{- end}}
{{- end}}

	AllColumns     {{dialect.PackageName}}.ColumnList
	MutableColumns {{dialect.PackageName}}.ColumnList
}

type {{tableTemplate.TypeName}} struct {
	{{structImplName}}

	{{toUpper insertedRowAlias}} {{structImplName}}
}

// AS creates new {{tableTemplate.TypeName}} with assigned alias
func (a {{tableTemplate.TypeName}}) AS(alias string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new {{tableTemplate.TypeName}} with assigned schema name
func (a {{tableTemplate.TypeName}}) FromSchema(schemaName string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new {{tableTemplate.TypeName}} with assigned table prefix
func (a {{tableTemplate.TypeName}}) WithPrefix(prefix string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new {{tableTemplate.TypeName}} with assigned table suffix
func (a {{tableTemplate.TypeName}}) WithSuffix(suffix string) *{{tableTemplate.TypeName}} {
	return new{{tableTemplate.TypeName}}(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func new{{tableTemplate.TypeName}}(schemaName, tableName, alias string) *{{tableTemplate.TypeName}} {
	return &{{tableTemplate.TypeName}}{
		{{structImplName}}: new{{tableTemplate.TypeName}}Impl(schemaName, tableName, alias),
		{{toUpper insertedRowAlias}}:  new{{tableTemplate.TypeName}}Impl("", "{{insertedRowAlias}}", ""),
	}
}

func new{{tableTemplate.TypeName}}Impl(schemaName, tableName, alias string) {{structImplName}} {
	var (
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
{{- if not $field.Skip }}
		{{$field.Name}}Column = {{dialect.PackageName}}.{{$field.Type}}Column("{{$c.Name}}")
{{- end}}
{{- end}}
		allColumns     = {{dialect.PackageName}}.ColumnList{ {{columnList .Columns}} }
		mutableColumns = {{dialect.PackageName}}.ColumnList{ {{columnList .MutableColumns}} }
	)

	return {{structImplName}}{
		Table: {{dialect.PackageName}}.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
{{- range $i, $c := .Columns}}
{{- $field := columnField $c}}
{{- if not $field.Skip }}
		{{$field.Name}}: {{$field.Name}}Column,
{{- end}}
{{- end}}

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
`

var tableSqlBuilderSetSchemaTemplate = `package {{package}}

// UseSchema sets a new schema name for all generated {{type}} SQL builder types. It is recommended to invoke 
// this method only once at the beginning of the program.
func UseSchema(schema string) {
{{- range .}}
	{{ .InstanceName }} = {{ .InstanceName }}.FromSchema(schema)
{{- end}}
}
`

var tableModelFileTemplate = `package {{package}}

{{ with modelImports }}
import (
{{- range .}}
	"{{.}}"
{{- end}}
)
{{end}}

{{$modelTableTemplate := tableTemplate}}
{{golangComment .Comment}}
type {{$modelTableTemplate.TypeName}} struct {
{{- range .Columns}}
{{- $field := structField .}}
{{- if not $field.Skip}}
	{{$field.Name}} {{$field.Type.Name}} ` + "{{$field.TagsString}}" + ` {{golangComment .Comment}}
{{- end }}
{{- end}}
}

`

var enumSQLBuilderTemplate = `package {{package}}

import "github.com/go-jet/jet/v2/{{dialect.PackageName}}"

{{golangComment .Comment}}
var {{enumTemplate.InstanceName}} = &struct {
{{- range $index, $value := .Values}}
	{{enumValueName $value}} {{dialect.PackageName}}.StringExpression
{{- end}}
} {
{{- range $index, $value := .Values}}
	{{enumValueName $value}}: {{dialect.PackageName}}.NewEnumValue("{{$value}}"),
{{- end}}
}
`

var enumModelTemplate = `package {{package}}
{{- $enumTemplate := enumTemplate}}

import "errors"

{{golangComment .Comment}}
type {{$enumTemplate.TypeName}} string

const (
{{- range $_, $value := .Values}}
	{{valueName $value}} {{$enumTemplate.TypeName}} = "{{$value}}"
{{- end}}
)

var {{$enumTemplate.TypeName}}AllValues = []{{$enumTemplate.TypeName}} {
{{- range $_, $value := .Values}}
	{{valueName $value}},
{{- end}}
}

func (e *{{$enumTemplate.TypeName}}) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
{{- range $_, $value := .Values}}
	case "{{$value}}":
		*e = {{valueName $value}}
{{- end}}
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for {{$enumTemplate.TypeName}} enum")
	}

	return nil
}

func (e {{$enumTemplate.TypeName}}) String() string {
	return string(e)
}

`
