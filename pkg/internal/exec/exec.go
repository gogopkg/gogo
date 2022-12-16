package exec

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

func Exec(ctx context.Context, env []string, command string, args ...string) (stdout, stderr string, err error) {
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err = cmd.Run()
	return outbuf.String(), errbuf.String(), err
}
