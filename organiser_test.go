package main

import (
	"os"
	"testing"
)

func Test_organise(t *testing.T) {
	type args struct {
		inputFolder  string
		outputFolder string
		copyOrMove   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				inputFolder:  "./input",
				outputFolder: "./output",
				copyOrMove:   "copy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := organise(tt.args.inputFolder, tt.args.outputFolder, tt.args.copyOrMove, nil, nil); (err != nil) != tt.wantErr {
				t.Errorf("organise() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.RemoveAll(tt.args.outputFolder)
		})
	}
}
