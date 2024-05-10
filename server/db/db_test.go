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
	var reply ReqResponse
	dumdb := NewDumDB()
	dumdb.Set(&SetArgs{"a", "11"}, &reply)
	dumdb.Set(&SetArgs{"b", "6"}, &reply)
	dumdb.Set(&SetArgs{"c", "7"}, &reply)

	expected := []struct {
		Key           string
		Op            string
		By            int64
		ExpectedValue string
	}{
		{"a", "+", 1, "12"},
		{"a", "-", 1, "11"},
		{"b", "+", 5, "11"},
		{"b", "-", 10, "1"},
		{"c", "+", 3, "10"},
		{"c", "-", 11, "-1"},
	}

	for i, exp := range expected {
		var reply ReqResponse
		dumdb.UpdateInt(exp.Key, exp.Op, exp.By, &reply)

		if !reply.Success {
			t.Fatalf("[test #%d] - UpdateInt command failed > %s", i, reply.Msg)
		}

		value := dumdb.Data[exp.Key]
		if value.Raw != exp.ExpectedValue {
			t.Fatalf("[test #%d] - expected %s got %s ", i, exp.ExpectedValue, value.Raw)
		}

	}
}
