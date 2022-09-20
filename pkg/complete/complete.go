package complete

import (
	"bytes"

	"github.com/chzyer/readline"
)

var runes = readline.Runes{}

type Completer struct {
	Name     []rune
	Dynamic  bool
	HasArgs  bool
	Callback readline.DynamicCompleteFunc
	Children []readline.PrefixCompleterInterface
}

func NewCompleter(args []string, pc ...readline.PrefixCompleterInterface) *Completer {
	return CompleterItem("", args, pc...)
}

func CompleterItem(name string, args []string, pc ...readline.PrefixCompleterInterface) *Completer {
	name += " "
	children := pc
	c := &Completer{
		Name:     []rune(name),
		Dynamic:  false,
		Children: children,
	}
	for _, v := range args {
		children = append(children, &Completer{
			Name:     []rune(v + "="),
			Dynamic:  false,
			Children: nil,
		})
	}
	c.HasArgs = len(args) > 0
	c.Children = children
	return c
}

func (p *Completer) Tree(prefix string) string {
	buf := bytes.NewBuffer(nil)
	p.Print(prefix, 0, buf)
	return buf.String()
}

func (p *Completer) Print(prefix string, level int, buf *bytes.Buffer) {
	readline.Print(p, prefix, level, buf)
}

func (p *Completer) IsDynamic() bool {
	return p.Dynamic
}

func (p *Completer) GetName() []rune {
	return p.Name
}

func (p *Completer) GetDynamicNames(line []rune) [][]rune {
	var names = [][]rune{}
	for _, name := range p.Callback(string(line)) {
		names = append(names, []rune(name+" "))
	}
	return names
}

func (p *Completer) GetChildren() []readline.PrefixCompleterInterface {
	return p.Children
}

func (p *Completer) SetChildren(children []readline.PrefixCompleterInterface) {
	p.Children = children
}

func (p *Completer) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	return doInternal(p, line, pos, line)
}

func isChildArg(name []rune) bool {
	return runes.Index('=', name) > -1
}

func processLine(childName []rune, line []rune) (newLine []rune, offset int, goNext bool) {
	if len(line) >= len(childName) {
		if runes.HasPrefix(line, childName) {
			if len(line) == len(childName) {
				return []rune{' '}, len(childName), true
			} else {
				return childName, len(childName), true
			}
		}
	} else {
		if runes.HasPrefix(childName, line) {
			return childName[len(line):], len(line), false
		}
	}
	return
}

func processArgLine(childName []rune, line []rune) (newLine []rune, offset int, goNext bool) {
	if len(line) >= len(childName) {
		lastRune := line[len(line)-1]
		if runes.HasPrefix(line, childName) {
			if lastRune != ' ' && lastRune != '=' {
				return []rune{' '}, len(line), false
			} else if lastRune == ' ' {
				return childName, len(line), true
			}
		}
	} else {
		if runes.HasPrefix(childName, line) {
			return childName[len(line):], len(line), false
		}
	}
	return
}

func doInternal(p *Completer, line []rune, pos int, origLine []rune) (newLine [][]rune, offset int) {
	line = runes.TrimSpaceLeft(line[:pos])
	goNext := false
	var lineCompleter *Completer
	var (
		tmpLine   []rune
		tmpOffset int
		tmpNext   bool
	)

	for _, child := range p.GetChildren() {
		childNames := make([][]rune, 1)

		childDynamic, ok := child.(readline.DynamicPrefixCompleterInterface)
		if ok && childDynamic.IsDynamic() {
			childNames = childDynamic.GetDynamicNames(origLine)
		} else {
			childNames[0] = child.GetName()
		}
		for _, childName := range childNames {
			if !isChildArg(childName) {
				tmpLine, tmpOffset, tmpNext = processLine(childName, line)
			} else {
				tmpLine, tmpOffset, tmpNext = processArgLine(childName, line)
			}
			if len(tmpLine) > 0 {
				newLine = append(newLine, tmpLine)
				offset = tmpOffset
				goNext = tmpNext
				if isChildArg(childName) {
					lineCompleter = p
				} else {
					lineCompleter = child.(*Completer)
				}
			}
		}
	}

	if len(newLine) != 1 {
		return
	}

	tmpLine = make([]rune, 0, len(line))
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}

		tmpLine = append(tmpLine, line[i:]...)
		return doInternal(lineCompleter, tmpLine, len(tmpLine), origLine)
	}

	if goNext {
		return doInternal(lineCompleter, nil, 0, origLine)
	}
	return
}
