package beaker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
)

type OptFn func(*Options) error

func WithWriter(w io.Writer) OptFn {
	return func(o *Options) error {
		if w == nil {
			return errors.New("passed writer is nil")
		}
		o.Writer = w
		return nil
	}
}

func New(funcs ...OptFn) (*Launcher, error) {
	opts := Options{
		Writer: os.Stdout,
	}
	for _, f := range funcs {
		if err := f(&opts); err != nil {
			return nil, err
		}
	}
	return &Launcher{
		Options: opts,
	}, nil
}

type Options struct {
	Writer io.Writer
}

type Launcher struct {
	Options Options
}

// Test launches the test suite defined in the launch pack
func (l *Launcher) Test(ctx context.Context, pack *LaunchPack) error {
	output, _, err := pack.Runner.Run(ctx)
	if err != nil {
		return fmt.Errorf("runner error: %w", err)
	}

	fmt.Printf("OUTPUT:\n%s\n", string(output))

	att, err := pack.Parser.ParseResults(ctx, output)
	if err != nil {
		return fmt.Errorf("parsing results: %w", err)
	}

	if l.Options.Writer == nil {
		return fmt.Errorf("tests ran successfully but no writer was configured")
	}

	encoder := protojson.MarshalOptions{
		Multiline:         true,
		Indent:            "  ",
		EmitDefaultValues: true,
	}

	jdata, err := encoder.Marshal(att)
	if err != nil {
		return fmt.Errorf("marshalling attestation: %w", err)
	}

	if _, err := l.Options.Writer.Write(jdata); err != nil {
		return fmt.Errorf("wiriting attestation data: %w", err)
	}

	return nil
}
