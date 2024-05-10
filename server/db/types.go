package db

type SetArgs struct {
	Key   string
	Value string
}

type GetArgs struct {
	Key string
}

type RawType string

const (
	INT    RawType = "INT"
	BOOL           = "BOOL"
	STRING         = "STRING"
)

type DataValue struct {
	Raw  string
	Type RawType
}

type DumDB struct {
	Data map[string]DataValue
}
