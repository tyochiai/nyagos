package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"../history"
	"../interpreter"
	"../readline"
)

type ICmdStream interface {
	ReadLine(*context.Context) (string, error)
}

type CmdStreamConsole struct {
	editor         *readline.LineEditor
	historyWrapper *THistory
	histPath       string
}

func NewCmdStreamConsole(it *interpreter.Interpreter) *CmdStreamConsole {
	editor := readline.NewLineEditor()
	editor.Prompt = printPrompt
	editor.Tag = it

	histPath := filepath.Join(AppDataDir(), "nyagos.history")
	historyWrapper := &THistory{editor}
	history.Load(histPath, historyWrapper)
	history.Save(histPath, historyWrapper)

	readline.DefaultEditor = editor

	return &CmdStreamConsole{
		editor:         editor,
		historyWrapper: historyWrapper,
		histPath:       histPath,
	}
}

func (this *CmdStreamConsole) ReadLine(ctx *context.Context) (string, error) {
	history_count := this.editor.HistoryLen()
	*ctx = context.WithValue(*ctx, "history", this.historyWrapper)
	var line string
	var err error
	for {
		line, err = this.editor.ReadLine()
		if err != nil {
			return line, err
		}
		var isReplaced bool
		line, isReplaced = history.Replace(this.historyWrapper, line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	if this.editor.HistoryLen() > history_count {
		fd, err := os.OpenFile(this.histPath, os.O_APPEND, 0600)
		if err != nil && os.IsNotExist(err) {
			// print("create ", this.histPath, "\n")
			fd, err = os.Create(this.histPath)
		}
		if err == nil {
			fmt.Fprintln(fd, line)
			fd.Close()
		} else {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	} else {
		this.editor.HistoryResetPointer()
	}
	return line, err
}

type CmdStreamFile struct {
	breader *bufio.Reader
}

func NewCmdStreamFile(r io.Reader) *CmdStreamFile {
	return &CmdStreamFile{breader: bufio.NewReader(r)}
}

func (this *CmdStreamFile) ReadLine(ctx *context.Context) (string, error) {
	line, err := this.breader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimRight(line, "\r\n")
	return line, nil
}

type UnCmdStream struct {
	body  ICmdStream
	queue []string
}

func NewUnCmdStream(body ICmdStream) *UnCmdStream {
	return &UnCmdStream{body: body, queue: nil}
}

func (this *UnCmdStream) ReadLine(ctx *context.Context) (string, error) {
	if this.queue == nil || len(this.queue) <= 0 {
		return this.body.ReadLine(ctx)
	} else {
		line := this.queue[0]
		this.queue = this.queue[1:]
		return line, nil
	}
}

func (this *UnCmdStream) UnreadLine(line string) {
	if this.queue == nil {
		this.queue = []string{line}
	} else {
		this.queue = append(this.queue, line)
	}
}