package db

import (
	"testing"
)

func TestInferType(t *testing.T) {
	tests := []struct {
		raw          string
		expectedType RawType
	}{
		{"11", INT},
		{"true", BOOL},
		{"str", STRING},
	}

	for i, expected := range tests {
		infered := inferType(expected.raw)
		if infered != expected.expectedType {
			t.Fatalf("[test #%d] - Wrong infered type, expected '%s' got '%s'", i, expected.expectedType, infered)
		}
	}
}

func TestUpdateInt(t *testing.T) {
	var reply Response
	dumdb := NewDumDB()
	dumdb.Set(&SetArgs{"a", "11"}, &reply)
	dumdb.Set(&SetArgs{"b", "6"}, &reply)
	dumdb.Set(&SetArgs{"c", "7"}, &reply)

	expected := []string{"12", "7", "8"}

	i := 0
	for key := range dumdb.Data {
		var reply Response
		dumdb.UpdateInt(key, "+", 1, &reply)

		if !reply.Success {
			t.Fatalf("[test #%d] - UpdateInt command failed > %s", i, reply.Msg)
		}

		value := dumdb.Data[key]
		if value.Raw != expected[i] {
			t.Fatalf("[test #%d] - expected %s got %s ", i, expected[i], value.Raw)
		}

		i++
	}
}
