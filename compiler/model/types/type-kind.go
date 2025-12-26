package types

type TypeKind int

const (
	INT TypeKind = iota
	BYTE
	FLOAT
	DECIMAL
	STRING
	BOOLEAN
	BLOB
	TYPEDESC
	TYPEREFDESC
	STREAM
	TABLE
	JSON
	XML
	ANY
	ANYDATA
	MAP
	FUTURE
	PACKAGE
	SERVICE
	CONNECTOR
	ENDPOINT
	FUNCTION
	ANNOTATION
	ARRAY
	UNION
	INTERSECTION
	VOID
	NIL
	NEVER
	NONE
	OTHER
	ERROR
	TUPLE
	OBJECT
	RECORD
	FINITE
	CHANNEL
	HANDLE
	READONLY
	TYPEPARAM
	PARAMETERIZED
	REGEXP
)

var typeKindNames = map[TypeKind]string{
	INT:           "int",
	BYTE:          "byte",
	FLOAT:         "float",
	DECIMAL:       "decimal",
	STRING:        "string",
	BOOLEAN:       "boolean",
	BLOB:          "blob",
	TYPEDESC:      "typedesc",
	TYPEREFDESC:   "typerefdesc",
	STREAM:        "stream",
	TABLE:         "table",
	JSON:          "json",
	XML:           "xml",
	ANY:           "any",
	ANYDATA:       "anydata",
	MAP:           "map",
	FUTURE:        "future",
	PACKAGE:       "package",
	SERVICE:       "service",
	CONNECTOR:     "connector",
	ENDPOINT:      "endpoint",
	FUNCTION:      "function",
	ANNOTATION:    "annotation",
	ARRAY:         "[]",
	UNION:         "|",
	INTERSECTION:  "&",
	VOID:          "",
	NIL:           "null",
	NEVER:         "never",
	NONE:          "",
	OTHER:         "other",
	ERROR:         "error",
	TUPLE:         "tuple",
	OBJECT:        "object",
	RECORD:        "record",
	FINITE:        "finite",
	CHANNEL:       "channel",
	HANDLE:        "handle",
	READONLY:      "readonly",
	TYPEPARAM:     "typeparam",
	PARAMETERIZED: "parameterized",
	REGEXP:        "regexp",
}

func (t TypeKind) TypeName() string {
	if name, ok := typeKindNames[t]; ok {
		return name
	}
	return ""
}
