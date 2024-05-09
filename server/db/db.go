package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func (db *DumDB) Export(reply *Response) error {
	json, err := json.Marshal(db.Data)
	if err != nil {
		msg := fmt.Errorf("Export failed > %s", err.Error())
		*reply = Response{false, msg.Error()}
		return msg
	}

	file, err := os.Create("export")
	if err != nil {
		msg := fmt.Errorf("Export failed > %s", err.Error())
		*reply = Response{false, msg.Error()}
		return msg
	}
	defer file.Close()

	file.Write(json)
	fmt.Printf("json >> %s\n", json)
	*reply = Response{true, "Export with success"}
	return nil
}

func (db *DumDB) Restore(reply *Response) error {
	file, err := os.Open("export")
	if err != nil {
		msg := fmt.Errorf("Restore failed > %s", err.Error())
		*reply = Response{false, msg.Error()}
		return msg
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Errorf("No export found", err.Error())
		*reply = Response{false, msg.Error()}
		return msg
	}
	err = json.Unmarshal(content, &db.Data)
	if err != nil {
		msg := fmt.Errorf("Unable to restore the export", err.Error())
		*reply = Response{false, msg.Error()}
		return msg
	}
	*reply = Response{true, "Restore with success"}
	return nil
}

func (db *DumDB) Request(req *ReqArgs, reply *Response) error {
	switch cmd := strings.Split(req.Request, " "); cmd[0] {
	case "HEALTH":
		db.Health(reply)
	case "EXPORT":
		db.Export(reply)
	case "RESTORE":
		db.Restore(reply)
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
