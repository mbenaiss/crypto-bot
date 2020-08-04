package csv

import (
	"testing"
)

func TestRead(t *testing.T) {
	expected := []map[string]string{
		{
			"time":    "2017-12-22 09:56:08",
			"type":    "deposit",
			"asset":   "ZEUR",
			"amount":  "2.0000",
			"fee":     "0.0000",
			"balance": "2.0000",
		},
	}
	columns := []string{
		"time", "type", "asset", "amount", "fee", "balance",
	}
	i := New(',', columns)
	actual, err := i.Read("./testdata/ledgers.csv")
	if err != nil {
		t.Fatal(err)
	}

	if (actual[0]["time"] != expected[0]["time"]) ||
		(actual[0]["type"] != expected[0]["type"]) ||
		(actual[0]["asset"] != expected[0]["asset"]) ||
		(actual[0]["amount"] != expected[0]["amount"]) ||
		(actual[0]["fee"] != expected[0]["fee"]) ||
		(actual[0]["balance"] != expected[0]["balance"]) {
		t.Fatalf("Expected \n%+v, got \n%+v", expected, actual)
	}
}
