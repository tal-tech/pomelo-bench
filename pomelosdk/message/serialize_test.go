package message

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Message
		wantErr bool
	}{

		{
			name: "t1",
			args: args{
				data: []byte{
					0,                                                                                               // 类型
					127,                                                                                             // id
					21,                                                                                              // route len
					99, 104, 97, 116, 46, 99, 104, 97, 116, 72, 97, 110, 100, 108, 101, 114, 46, 115, 101, 110, 100, // route
					99, 99, 99, 99, // data
				},
			},
			want: &Message{
				Type:       0,
				ID:         127,
				Route:      "chat.chatHandler.send",
				Data:       []byte{99, 99, 99, 99},
				compressed: false,
			},
			wantErr: false,
		},
		{
			name: "t2",
			args: args{
				data: []byte{
					0,   // 类型
					128, // id
					1,
					21,                                                                                              // route len
					99, 104, 97, 116, 46, 99, 104, 97, 116, 72, 97, 110, 100, 108, 101, 114, 46, 115, 101, 110, 100, // route
					99, 99, 99, 99, // data
				},
			},
			want: &Message{
				Type:       0,
				ID:         128,
				Route:      "chat.chatHandler.send",
				Data:       []byte{99, 99, 99, 99},
				compressed: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		m *Message
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				m: &Message{
					Type:       0,
					ID:         129,
					Route:      "chat.chatHandler.send",
					Data:       []byte{99, 99, 99, 99},
					compressed: false,
				},
			},
			want: []byte{
				0,                                                                                               // 类型
				129,                                                                                             // id
				1,                                                                                               // id 高位
				21,                                                                                              // route len
				99, 104, 97, 116, 46, 99, 104, 97, 116, 72, 97, 110, 100, 108, 101, 114, 46, 115, 101, 110, 100, // route
				99, 99, 99, 99, // data
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
