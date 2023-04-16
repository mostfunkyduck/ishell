package ishell

// Jack: see documentation in fileCompleter function for main shell for why and how this was hacked
import (
	"strings"

	"github.com/flynn-archive/go-shlex"
)

type iCompleter struct {
	cmd      *Cmd
	disabled func() bool
}

func (ic iCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if ic.disabled != nil && ic.disabled() {
		return nil, len(line)
	}
	var words []string
	if w, err := shlex.Split(string(line)); err == nil {
		words = w
	} else {
		// fall back
		words = strings.Fields(string(line))
	}

	var wordsWithoutFlags []string
	for _, w := range words {
		if !strings.HasPrefix(w, "-") {
			wordsWithoutFlags = append(wordsWithoutFlags, w)
		}
	}

	var cWords []string
	prefix := ""
	if len(wordsWithoutFlags) > 0 && pos > 0 && line[pos-1] != ' ' {
		prefix = wordsWithoutFlags[len(wordsWithoutFlags)-1]
		cWords = ic.getWords(prefix, words[:len(wordsWithoutFlags)-1])
	} else {
		cWords = ic.getWords(prefix, wordsWithoutFlags)
	}

	var suggestions [][]rune
	for _, w := range cWords {
		suggestions = append(suggestions, []rune(w))
	}
	return suggestions, len(prefix)
}

func (ic iCompleter) getWords(prefix string, w []string) (s []string) {
	cmd, args := ic.cmd.FindCmd(w)
	if cmd == nil {
		cmd, args = ic.cmd, w
	}
	if cmd.CompleterWithPrefix != nil {
		return cmd.CompleterWithPrefix(prefix, args)
	}
	if cmd.Completer != nil {
		return cmd.Completer(args)
	}
	for k := range cmd.children {
		s = append(s, k)
	}
	return
}
