package exporter

import (
	"reflect"
	"testing"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

func TestBuildSheetRows(t *testing.T) {
	rows := buildSheetRows(app.ExtractedData{
		Summary:      "Weekly sync",
		Participants: []string{"Alice", "Bob"},
		Decisions:    []string{"Ship the release", "Use webhook export"},
		ActionItems: []app.ActionItems{
			{Owner: "Alice", Task: "Prepare release notes", Due: "Friday"},
			{Owner: "Bob", Task: "Send follow-up", Due: ""},
		},
	})

	expected := [][]interface{}{
		{"", "", "", "", "", ""},
		{"Summary", "Weekly sync", "", "", "", ""},
		{"Participants", "Alice", "", "", "", ""},
		{"", "Bob", "", "", "", ""},
		{"Decisions", "Ship the release", "", "", "", ""},
		{"", "Use webhook export", "", "", "", ""},
		{"Action Items", "Owner", "Task", "Due", "", ""},
		{"", "Alice", "Prepare release notes", "Friday", "", ""},
		{"", "Bob", "Send follow-up", "", "", ""},
	}

	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("unexpected rows:\nexpected: %#v\nactual:   %#v", expected, rows)
	}
}
