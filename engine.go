package jiq

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	DefaultY     int    = 1
	FilterPrompt string = "[Filter]> "
)

type Engine struct {
	json          string
	query         *Query
	args          []string
	term          *Terminal
	complete      []string
	keymode       bool
	candidates    []string
	candidatemode bool
	candidateidx  int
	contentOffset int
	queryConfirm  bool
	cursorOffsetX int
}

func NewEngine(s io.Reader, args []string) *Engine {
	j, err := ioutil.ReadAll(s)
	if err != nil {
		return &Engine{}
	}
	e := &Engine{
		json:          string(j),
		term:          NewTerminal(FilterPrompt, DefaultY),
		query:         NewQuery([]rune("")),
		args:          args,
		complete:      []string{"", ""},
		keymode:       false,
		candidates:    []string{},
		candidatemode: false,
		candidateidx:  0,
		contentOffset: 0,
		queryConfirm:  false,
		cursorOffsetX: 0,
	}
	return e
}

type EngineResult struct {
	Content string
	Qs      string
	Err     error
}

func (e *Engine) Run() *EngineResult {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var contents []string

	for {
		contents = e.getContents()
		e.setCandidateData()
		e.queryConfirm = false

		ta := &TerminalDrawAttributes{
			Query:           e.query.StringGet(),
			CursorOffsetX:   e.cursorOffsetX,
			Contents:        contents,
			CandidateIndex:  e.candidateidx,
			ContentsOffsetY: e.contentOffset,
			Complete:        e.complete[0],
			Candidates:      e.candidates,
		}

		e.term.draw(ta)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 0:
				e.inputChar(ev.Ch)
			case termbox.KeySpace:
				e.inputChar(32)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				e.deleteChar()
			case termbox.KeyDelete:
				e.deleteNextChar()
			case termbox.KeyTab:
				e.tabAction()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				e.moveCursorBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				e.moveCursorForward()
			case termbox.KeyHome, termbox.KeyCtrlA:
				e.moveCursorToTop()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				e.moveCursorToEnd()
			case termbox.KeyCtrlK:
				e.scrollToAbove()
			case termbox.KeyCtrlJ:
				e.scrollToBelow()
			case termbox.KeyCtrlL:
				e.toggleKeymode()
			case termbox.KeyCtrlW:
				e.deleteWordBackward()
			case termbox.KeyEsc:
				e.escapeCandidateMode()
			case termbox.KeyEnter:
				if !e.candidatemode {
					cc, err := jqrun(e.query.StringGet(), e.json, e.args)

					return &EngineResult{
						Content: cc,
						Qs:      e.query.StringGet(),
						Err:     err,
					}
				}
				e.confirmCandidate()
			case termbox.KeyCtrlC:
				return &EngineResult{}
			default:
			}
		case termbox.EventError:
			panic(ev.Err)
			break
		default:
		}
	}
}

func (e *Engine) getContents() []string {
	var contents []string

	cc, _ := jqrun(e.query.StringGet(), e.json, e.args)

	if e.keymode {
		contents = e.candidates
	} else {
		contents = strings.Split(cc, "\n")
	}
	return contents
}

func (e *Engine) setCandidateData() {
	if l := len(e.candidates); e.complete[0] == "" && l > 1 {
		if e.candidateidx >= l {
			e.candidateidx = 0
		}
	} else {
		e.candidatemode = false
	}
	if !e.candidatemode {
		e.candidateidx = 0
		e.candidates = []string{}
	}
}

func (e *Engine) confirmCandidate() {
	_, _ = e.query.PopKeyword()
	_ = e.query.StringAdd(".")
	q := e.query.StringAdd(e.candidates[e.candidateidx])
	e.cursorOffsetX = len(q)
	e.queryConfirm = true
}

func (e *Engine) deleteChar() {
	if e.cursorOffsetX > 0 {
		_ = e.query.Delete(e.cursorOffsetX - 1)
		e.cursorOffsetX -= 1
	}
}
func (e *Engine) deleteNextChar() {
	e.query.Delete(e.cursorOffsetX)
}
func (e *Engine) scrollToBelow() {
	e.contentOffset++
}
func (e *Engine) scrollToAbove() {
	if o := e.contentOffset - 1; o >= 0 {
		e.contentOffset = o
	}
}
func (e *Engine) toggleKeymode() {
	e.keymode = !e.keymode
}
func (e *Engine) deleteWordBackward() {
	if k, _ := e.query.StringPopKeyword(); k != "" && !strings.Contains(k, "[") {
		_ = e.query.StringAdd(".")
	}
	e.cursorOffsetX = len(e.query.Get())
}
func (e *Engine) tabAction() {
	if !e.candidatemode {
		e.candidatemode = true
		if e.query.StringGet() == "" {
			_ = e.query.StringAdd(".")
		} else if e.complete[0] != e.complete[1] && e.complete[0] != "" {
			if k, _ := e.query.StringPopKeyword(); !strings.Contains(k, "[") {
				_ = e.query.StringAdd(".")
			}
			_ = e.query.StringAdd(e.complete[1])
		} else {
			_ = e.query.StringAdd(e.complete[0])
		}
	} else {
		e.candidateidx = e.candidateidx + 1
	}
	e.cursorOffsetX = len(e.query.Get())
}
func (e *Engine) escapeCandidateMode() {
	e.candidatemode = false
}
func (e *Engine) inputChar(ch rune) {
	b := len(e.query.Get())
	q := e.query.StringInsert(string(ch), e.cursorOffsetX)
	if b < len(q) {
		e.cursorOffsetX += 1
	}
}

func (e *Engine) moveCursorBackward() {
	if e.cursorOffsetX > 0 {
		e.cursorOffsetX -= 1
	}
}
func (e *Engine) moveCursorForward() {
	if len(e.query.Get()) > e.cursorOffsetX {
		e.cursorOffsetX += 1
	}
}
func (e *Engine) moveCursorWordBackwark() {
}
func (e *Engine) moveCursorWordForward() {
}
func (e *Engine) moveCursorToTop() {
	e.cursorOffsetX = 0
}
func (e *Engine) moveCursorToEnd() {
	e.cursorOffsetX = len(e.query.Get())
}
