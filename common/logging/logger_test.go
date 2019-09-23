package logging_test

import (
	"reflect"
	"testing"

	"github.com/fernandoocampo/thepingthepong/common/logging"
	"github.com/sirupsen/logrus"
)

func TestLogHandle_SetFormat(t *testing.T) {
	lf := logrus.Fields{"foo": "bar"}
	logger := logrus.StandardLogger().WithFields(lf)

	type fields struct {
		Entry *logrus.Entry
	}
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "format is empty", fields: fields{Entry: logger}, args: args{format: ""}, wantErr: true},
		{name: "format is text", fields: fields{Entry: logger}, args: args{format: "text"}, wantErr: false},
		{name: "format is json", fields: fields{Entry: logger}, args: args{format: "json"}, wantErr: false},
		{name: "format is invalid", fields: fields{Entry: logger}, args: args{format: "invalid"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &logging.Handle{
				Entry: tt.fields.Entry,
			}
			if err := lh.SetFormat(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("Handle.SetFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	lf := logrus.Fields{"foo": "bar"}
	logger := logrus.StandardLogger().WithFields(lf)
	lh := &logging.Handle{logger}

	tests := []struct {
		name    string
		args    logging.Options
		want    *logging.Handle
		wantErr bool
	}{
		{name: "LogLevel is empty", args: logging.Options{LogLevel: "", LogFormat: "text", LogFields: lf}, want: nil, wantErr: true},
		{name: "LogFormat is empty", args: logging.Options{LogLevel: "info", LogFormat: "", LogFields: lf}, want: nil, wantErr: true},
		{name: "LogFormat and LogLevel empty", args: logging.Options{LogLevel: "", LogFormat: "", LogFields: lf}, want: nil, wantErr: true},
		{name: "LogLevel is wrong", args: logging.Options{LogLevel: "pepe", LogFormat: "text", LogFields: lf}, want: nil, wantErr: true},
		{name: "LogFormat is wrong", args: logging.Options{LogLevel: "info", LogFormat: "pepe", LogFields: lf}, want: nil, wantErr: true},
		{name: "LogFormat is json", args: logging.Options{LogLevel: "info", LogFormat: "json", LogFields: lf}, want: lh, wantErr: false},
		{name: "LogLevel is warn", args: logging.Options{LogLevel: "warn", LogFormat: "text", LogFields: lf}, want: lh, wantErr: false},
		{name: "LogLevel is info", args: logging.Options{LogLevel: "info", LogFormat: "text", LogFields: lf}, want: lh, wantErr: false},
		{name: "LogLevel is debug", args: logging.Options{LogLevel: "debug", LogFormat: "text", LogFields: lf}, want: lh, wantErr: false},
		{name: "Options is empty", args: logging.Options{}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logging.NewLogger(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogHandle_SetLevel(t *testing.T) {
	lf := logrus.Fields{"foo": "bar"}
	logger := logrus.StandardLogger().WithFields(lf)

	type fields struct {
		Entry *logrus.Entry
	}
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "level is empty", fields: fields{Entry: logger}, args: args{level: ""}, wantErr: true},
		{name: "level is warn", fields: fields{Entry: logger}, args: args{level: "warn"}, wantErr: false},
		{name: "level is info", fields: fields{Entry: logger}, args: args{level: "info"}, wantErr: false},
		{name: "level is debug", fields: fields{Entry: logger}, args: args{level: "debug"}, wantErr: false},
		{name: "level is invalid", fields: fields{Entry: logger}, args: args{level: "invalid"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &logging.Handle{
				Entry: tt.fields.Entry,
			}
			if err := lh.SetLevel(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Handle.SetLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
