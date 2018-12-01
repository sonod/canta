package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./canta -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("canta version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

func TestParseConsulEvents(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *ConsulEvent
		wantErr bool
	}{
		{
			name: "Success: eventInPayload",
			args: args{
				in: bytes.NewReader([]byte(`
				[
				  {
				      "ID": "f83fa997-91c5-02b7-f3b3-70949f3e9f55",
				      "Name": "hello",
				      "Payload": "SGV5ISBZb3UgYXJlIGluIHRoZXJlISBZb3VyIGhvdXNlIGlzIGhhdW50ZWQh",
				      "NodeFilter": "",
				      "ServiceFilter": "",
				      "TagFilter": "",
				      "Version": 1,
				      "LTime": 52494
				  }
				]
				`)),
			},
			want: &ConsulEvent{
				ID:      "f83fa997-91c5-02b7-f3b3-70949f3e9f55",
				Name:    "hello",
				Payload: []byte("Hey! You are in there! Your house is haunted!"),
				LTime:   52494,
			},
		},
		{
			name: "Failed: eventEmptyPayload",
			args: args{
				in: bytes.NewReader([]byte(`[]`)),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConsulEvents(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConsulEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConsulEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
