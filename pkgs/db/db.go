package db

import (
	"fmt"
	"strconv"
	"strings"
)

func (db *DumDB) Health(reply *Response) {
	msg := fmt.Sprintf("Server is up, %d elements in db", len(db.Data))
	*reply = Response{true, msg}
}

func (db *DumDB) Set(args *SetArgs, reply *Response) error {
	dataValue := DataValue{args.Value, inferType(args.Value)}
	db.Data[args.Key] = dataValue
	*reply = Response{true, args.Value}

	return nil
}

func (db *DumDB) Get(args *GetArgs, reply *Response) error {
	if value, ok := db.Data[args.Key]; ok {
		*reply = Response{true, value.Raw}
		return nil
	}

	msg := fmt.Errorf("No value for key '%s'", args.Key)
	*reply = Response{false, msg.Error()}
	return msg
}

func (db *DumDB) UpdateInt(key string, operator string, by int64, reply *Response) error {
	if value, ok := db.Data[key]; ok {
		if value.Type != INT {
			msg := fmt.Errorf("Cannot increment '%s' type", value.Type)
			*reply = Response{false, msg.Error()}
			return msg
		}

		intValue, _ := strconv.ParseInt(value.Raw, 10, 64)
		var updatedValue int64
		switch operator {
		case "+":
			updatedValue = intValue + by
		case "-":
			updatedValue = intValue - by
		}
		db.Data[key] = DataValue{Raw: fmt.Sprintf("%d", updatedValue), Type: value.Type}
		*reply = Response{true, db.Data[key].Raw}
		return nil
	}

	msg := fmt.Errorf("No value for key '%s'", key)
	*reply = Response{false, msg.Error()}
	return msg
}

func (db *DumDB) Request(req *ReqArgs, reply *Response) error {
	switch cmd := strings.Split(req.Request, " "); cmd[0] {
	case "HEALTH":
		db.Health(reply)
	case "SET":
		setArgs := &SetArgs{cmd[1], cmd[2]}
		db.Set(setArgs, reply)
	case "GET":
		getArgs := &GetArgs{cmd[1]}
		db.Get(getArgs, reply)
	case "INCR":
		db.UpdateInt(cmd[1], "+", 1, reply)
	case "DECR":
		db.UpdateInt(cmd[1], "-", 1, reply)
	case "INCRBY":
		by, err := strconv.ParseInt(cmd[2], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid increment value, got '%s'", cmd[2])
		}
		db.UpdateInt(cmd[1], "+", by, reply)
	case "DECRBY":
		by, err := strconv.ParseInt(cmd[2], 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid increment value, got '%s'", cmd[2])
		}
		db.UpdateInt(cmd[1], "-", by, reply)
	default:
		return fmt.Errorf("Unknown cmd '%s'\n", cmd)
	}

	return nil
}
