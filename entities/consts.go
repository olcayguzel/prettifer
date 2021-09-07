package entities

type Alignment uint16
type DataType string

const (
	AL_LEFT   Alignment = 1
	AL_CENTER Alignment = 2
	AL_RIGTH  Alignment = 3
)

const (
	DT_STRING   = "string"
	DT_INTEGER  = "integer"
	DT_FLOAT    = "float"
	DT_DATETIME = "datetime"
	DT_COMPLEX  = "complex"
	DT_STRUCT   = "struct"
	DT_BOOLEAN  = "boolean"
)
