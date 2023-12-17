package model

import (
	"reflect"
	"testing"
)

func TestMessage_AsJson(t *testing.T) {
	type fields struct {
		MessageType   string
		MessageString string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Success #1",
			fields: fields{
				MessageType:   "Response",
				MessageString: "123456789",
			},
			want: []byte("{\"messageType\":\"Response\",\"messageString\":\"123456789\"}"),
		},
		{
			name: "Success #2",
			fields: fields{
				MessageType:   "",
				MessageString: "",
			},
			want: []byte("{\"messageType\":\"\",\"messageString\":\"\"}"),
		},
		{
			name:   "Success #3",
			fields: fields{},
			want:   []byte("{\"messageType\":\"\",\"messageString\":\"\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				MessageType:   tt.fields.MessageType,
				MessageString: tt.fields.MessageString,
			}
			if got := m.AsJson(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.AsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
