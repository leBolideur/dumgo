package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type ReqArgs struct {
	Request string
}

type ReqResponse struct {
	Success bool
	Msg     string
}

func (db *DumDB) Health(reply *ReqResponse) {
	msg := fmt.Sprintf("Server is up, %d elements in db", len(db.Data))
	*reply = ReqResponse{Success: true, Msg: msg}
}

func (db *DumDB) Set(args *SetArgs, reply *ReqResponse) error {
	dataValue := DataValue{args.Value, inferType(args.Value)}
	db.Data[args.Key] = dataValue
	*reply = ReqResponse{Success: true, Msg: args.Value}

	return nil
}

func (db *DumDB) Get(args *GetArgs, reply *ReqResponse) error {
	if value, ok := db.Data[args.Key]; ok {
		*reply = ReqResponse{Success: true, Msg: value.Raw}
		return nil
	}

	msg := fmt.Errorf("No value for key '%s'", args.Key)
	*reply = ReqResponse{Success: false, Msg: msg.Error()}
	return msg
}

func (db *DumDB) UpdateInt(key string, operator string, by int64, reply *ReqResponse) error {
	if value, ok := db.Data[key]; ok {
		if value.Type != INT {
			msg := fmt.Errorf("Cannot increment '%s' type", value.Type)
			*reply = ReqResponse{Success: false, Msg: msg.Error()}
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
		*reply = ReqResponse{Success: true, Msg: db.Data[key].Raw}
		return nil
	}

	msg := fmt.Errorf("No value for key '%s'", key)
	*reply = ReqResponse{Success: false, Msg: msg.Error()}
	return msg
}

func (db *DumDB) Export(reply *ReqResponse) error {
	json, err := json.Marshal(db.Data)
	if err != nil {
		msg := fmt.Errorf("Export failed > %s", err.Error())
		*reply = ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}

	file, err := os.Create("export")
	if err != nil {
		msg := fmt.Errorf("Export failed > %s", err.Error())
		*reply = ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}
	defer file.Close()

	file.Write(json)
	fmt.Printf("json >> %s\n", json)
	*reply = ReqResponse{Success: true, Msg: "Export with success"}
	return nil
}

func (db *DumDB) Restore(reply *ReqResponse) error {
	file, err := os.Open("export")
	if err != nil {
		msg := fmt.Errorf("Restore failed > %s", err.Error())
		*reply = ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Errorf("No export found > %s", err.Error())
		*reply = ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}
	err = json.Unmarshal(content, &db.Data)
	if err != nil {
		msg := fmt.Errorf("Unable to restore the export > %s", err.Error())
		*reply = ReqResponse{Success: false, Msg: msg.Error()}
		return msg
	}
	*reply = ReqResponse{Success: true, Msg: "Restore with success"}
	return nil
}
