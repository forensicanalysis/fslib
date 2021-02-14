package fslib_test

import (
	"runtime"
	"testing"

	"github.com/forensicanalysis/fslib"
)

func TestToForensicPath(t *testing.T) {
	type args struct {
		systemPath string
	}
	tests := []struct {
		name        string
		windowsTest bool
		args        args
		wantName    string
		wantErr     bool
	}{
		{"Windows Abs Path", true, args{"C:\\Windows"}, "C/Windows", false},
		// {"Windows Rel Path", true, args{"\\Windows"}, "/C/Windows", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.windowsTest && runtime.GOOS == "windows") || !tt.windowsTest {
				gotName, err := fslib.ToFSPath(tt.args.systemPath)
				if (err != nil) != tt.wantErr {
					t.Errorf("ToFSPath() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotName != tt.wantName {
					t.Errorf("ToFSPath() gotName = %v, want %v", gotName, tt.wantName)
				}
			}
		})
	}
}
