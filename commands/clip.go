package commands

import (
	"context"
	"io/ioutil"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/zetamatta/go-mbcs"

	"github.com/zetamatta/nyagos/shell"
)

func cmd_clip(ctx context.Context, cmd *shell.Cmd) (int, error) {
	data, err := ioutil.ReadAll(cmd.Stdin)
	if err != nil {
		return 1, err
	}
	if utf8.Valid(data) {
		clipboard.WriteAll(string(data))
	} else {
		str, err := mbcs.AtoU(data)
		if err == nil {
			clipboard.WriteAll(str)
		} else {
			return 2, err
		}
	}
	return 0, nil
}
