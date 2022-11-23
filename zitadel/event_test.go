package zitadel

import (
	"reflect"
	"testing"
)

func TestEventFromCommand(t *testing.T) {
	type args struct {
		cmd Command
	}
	tests := []struct {
		name string
		args args
		want *Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EventFromCommand(tt.args.cmd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EventFromCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
