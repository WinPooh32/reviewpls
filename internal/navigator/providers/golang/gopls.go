package golang

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/sourcegraph/jsonrpc2"
)

type JSON = map[string]any

type stdio struct {
	in  io.ReadCloser
	out io.WriteCloser
}

func (s *stdio) Read(p []byte) (n int, err error) {
	//nolint
	return s.in.Read(p)
}

func (s *stdio) Write(p []byte) (n int, err error) {
	//nolint
	return s.out.Write(p)
}

func (s *stdio) Close() error {
	return errors.Join(s.in.Close(), s.out.Close())
}

type handler struct{}

func (h *handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {}

type gopls struct {
	proc *exec.Cmd
	conn *jsonrpc2.Conn
}

func newGopls(ctx context.Context, workspaceFolders []string) (*gopls, error) {
	cli := &gopls{
		proc: exec.Command("gopls"),
		conn: nil,
	}

	go func() {
		in, err := cli.proc.StderrPipe()
		if err != nil {
			return
		}

		r := bufio.NewReader(in)

		for {
			str, err := r.ReadString('\n')
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				panic(err)
			}

			fmt.Fprintln(os.Stderr, str)
		}
	}()

	in, err := cli.proc.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("proc stdout pipe: %w", err)
	}

	out, err := cli.proc.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("proc stdin pipe: %w", err)
	}

	if err := cli.proc.Start(); err != nil {
		return nil, fmt.Errorf("start gopls: %w", err)
	}

	cli.conn = jsonrpc2.NewConn(
		ctx,
		jsonrpc2.NewBufferedStream(&stdio{in, out}, jsonrpc2.VSCodeObjectCodec{}),
		&handler{},
		jsonrpc2.LogMessages(log.Default()),
	)

	if err := cli.initWorkspace(ctx, workspaceFolders); err != nil {
		return nil, fmt.Errorf("init workspace: %w", err)
	}

	return cli, nil
}

func (g *gopls) initWorkspace(ctx context.Context, workspaceFolders []string) (err error) {
	workspace := make([]JSON, 0, len(workspaceFolders))

	for i, folder := range workspaceFolders {
		absDir, err := filepath.Abs(folder)
		if err != nil {
			return fmt.Errorf("make absolute path of the workspace folder %q: %w", folder, err)
		}

		basename := filepath.Base(absDir)
		if basename == "." || basename == "/" {
			basename = ""
		}

		name := strconv.Itoa(i) + "-" + basename

		workspace = append(workspace, JSON{
			"name": name,
			"uri":  "file://" + absDir,
		})
	}

	var result any
	if err := g.conn.Call(ctx,
		"initialize",
		JSON{
			"clientInfo":       JSON{"name": "reviewpls"},
			"workspaceFolders": workspace,
		},
		&result,
	); err != nil {
		return fmt.Errorf("request initialize: %w", err)
	}

	var resinit JSON
	if err := g.conn.Call(ctx, "initialized", JSON{}, &resinit); err != nil {
		return fmt.Errorf("req initialized: %w", err)
	}

	return nil
}

type documentSymbol struct {
	Name           string   `json:"name"`
	Kind           int      `json:"kind"`
	SelectionRange locRange `json:"selection_range"`
}

type locRange struct {
	Start location `json:"start"`
	End   location `json:"end"`
}

type location struct {
	Line      uint `json:"line"`
	Character uint `json:"character"`
}

func (g *gopls) documentSymbols(ctx context.Context, document string) ([]documentSymbol, error) {
	absPath, err := filepath.Abs(document)
	if err != nil {
		return nil, fmt.Errorf("make absolute path of the document %q: %w", document, err)
	}

	var respSymbols []documentSymbol

	if err := g.conn.Call(ctx,
		"textDocument/documentSymbol",
		JSON{"textDocument": JSON{"uri": "file://" + absPath}},
		&respSymbols,
	); err != nil {
		return nil, fmt.Errorf("request document symbols: %w", err)
	}

	if respSymbols == nil {
		return nil, nil
	}

	return respSymbols, nil
}
