package db

import "strconv"

func NewDumDB() *DumDB {
	return &DumDB{
		Data: make(map[string]DataValue),
	}
}

func (db *DumDB) NewArgs(key, value string) SetArgs {
	return SetArgs{
		Key:   key,
		Value: value,
	}
}

func inferType(raw string) RawType {
	if _, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return INT
	} else if _, err := strconv.ParseBool(raw); err == nil {
		return BOOL
	}

	return STRING
}
