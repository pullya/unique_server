package model

import (
	"reflect"
	"testing"
)

func TestMessage_AsJson(t *testing.T) {
	type fields struct {
		Response string
		Text     string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Success #1",
			fields: fields{
				Response: "OK",
				Text:     "123456789",
			},
			want: []byte("{\"response\":\"OK\",\"text\":\"123456789\"}"),
		},
		{
			name: "Success #2",
			fields: fields{
				Response: "",
				Text:     "",
			},
			want: []byte("{\"response\":\"\",\"text\":\"\"}"),
		},
		{
			name:   "Success #3",
			fields: fields{},
			want:   []byte("{\"response\":\"\",\"text\":\"\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Response: tt.fields.Response,
				Text:     tt.fields.Text,
			}
			if got := m.AsJson(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.AsJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
