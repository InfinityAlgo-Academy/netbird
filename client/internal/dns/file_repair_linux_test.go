//go:build !android

package dns

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/netbirdio/netbird/util"
)

func TestMain(m *testing.M) {
	_ = util.InitLog("debug", "console")
	code := m.Run()
	os.Exit(code)
}

func Test_newRepairtmp(t *testing.T) {
	type args struct {
		resolvConfContent  string
		touchedConfContent string
		wantChange         bool
	}
	tests := []args{
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
nameserver 8.8.8.8
searchdomain netbird.cloud something`,
			wantChange: true,
		},
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something somethingelse`,
			wantChange: false,
		},
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
nameserver 10.0.0.1
searchdomain netbird.cloud something`,
			wantChange: false,
		},
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
searchdomain something`,
			wantChange: true,
		},
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
nameserver 10.0.0.1`,
			wantChange: true,
		},
		{
			resolvConfContent: `
nameserver 10.0.0.1
nameserver 8.8.8.8
searchdomain netbird.cloud something`,

			touchedConfContent: `
nameserver 8.8.8.8`,
			wantChange: true,
		},
	}

	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			t.Parallel()
			workDir := t.TempDir()
			operationFile := workDir + "/resolv.conf"
			err := os.WriteFile(operationFile, []byte(tt.resolvConfContent), 0755)
			if err != nil {
				t.Fatalf("failed to wrtie out resolv.conf: %s", err)
			}

			var changed bool
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			updateFn := func([]string, string, *resolvConf) error {
				changed = true
				cancel()
				return nil
			}

			r := newRepair(operationFile, updateFn)
			r.watchFileChanges([]string{"netbird.cloud"}, "10.0.0.1")

			err = os.WriteFile(operationFile, []byte(tt.touchedConfContent), 0755)
			if err != nil {
				t.Fatalf("failed to wrtie out resolv.conf: %s", err)
			}

			<-ctx.Done()

			r.stopWatchFileChanges()

			if changed != tt.wantChange {
				t.Errorf("unexpected result: want: %v, get: %v", tt.wantChange, changed)
			}
		})
	}

}