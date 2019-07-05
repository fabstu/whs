package helper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"

	"aduu.dev/tools/aduu/helper/testhelper"
)

func TestRun(t *testing.T) {
	type args struct {
		cmd  *exec.Cmd
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.cmd, tt.args.name)
		})
	}
}

func TestRunE(t *testing.T) {
	type args struct {
		cmd  *exec.Cmd
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunE(tt.args.cmd, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("RunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunWithOutput(t *testing.T) {
	type args struct {
		cmd  *exec.Cmd
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunWithOutput(tt.args.cmd, tt.args.name); got != tt.want {
				t.Errorf("RunWithOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunWithOutputE(t *testing.T) {
	type args struct {
		cmd  *exec.Cmd
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunWithOutputE(tt.args.cmd, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunWithOutputE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunWithOutputE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitCommand(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"kubectl get external ips", args{`kubectl get svc --all-namespaces -o jsonpath='{range .items[?(@.spec.type=="LoadBalancer")]}{.metadata.name}:{.status.loadBalancer.ingress[0].ip}{"\n"}{end}'`},
			[]string{"kubectl", "get", "svc", "--all-namespaces", "-o", `jsonpath={range .items[?(@.spec.type=="LoadBalancer")]}{.metadata.name}:{.status.loadBalancer.ingress[0].ip}{"\n"}{end}`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitCommand(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, got, tt.want)
				t.Errorf("SplitCommand() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	type args struct {
		s   string
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute(tt.args.s, tt.args.obj)
		})
	}
}

func TestExecuteWithOutput(t *testing.T) {
	type args struct {
		s   string
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExecuteWithOutput(tt.args.s, tt.args.obj); got != tt.want {
				t.Errorf("ExecuteWithOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecuteWithOutputE(t *testing.T) {
	type args struct {
		s   string
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test expansion.", args{`echo $HOME`, ""}, os.ExpandEnv("$HOME"), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteWithOutputE(tt.args.s, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithOutputE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExecuteWithOutputE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Execute_UsesWithDir(t *testing.T) {
	tempDir := testhelper.MakeTempDir(t, "execute_UseWithDir")
	defer testhelper.DeleteTempDir(t, tempDir)

	// Creating a so i can try to touch a/b and so won't create a if it did not switch directories correctly
	a := filepath.Join(tempDir, "a")
	b := filepath.Join(tempDir, "a", "b")

	if err := os.Mkdir(a, 0777); err != nil {
		t.Fatalf("failed to create test dir a: %v", err)
	}

	Execute("touch a/b", "", WithDir(tempDir))

	if !DoesPathExist(b) {
		t.Errorf("failed to change directory into %s", tempDir)
		return
	}
	fmt.Println("b does exist:", b)
}
