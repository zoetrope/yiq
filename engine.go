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
	candidates    []string
	candidatemode bool
	candidateidx  int
	contentOffset int
	cursorOffsetX int
}

func NewEngine(s io.Reader, args []string, initialquery string) *Engine {
	j, err := ioutil.ReadAll(s)
	if err != nil {
		return &Engine{}
	}
	e := &Engine{
		json:          string(j),
		term:          NewTerminal(FilterPrompt, DefaultY),
		query:         NewQuery([]rune(initialquery)),
		args:          args,
		complete:      []string{"", ""},
		candidates:    []string{},
		candidatemode: false,
		candidateidx:  0,
		contentOffset: 0,
		cursorOffsetX: len(initialquery),
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
		e.setCandidates()
		e.setCandidateData()

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
			case termbox.KeyArrowLeft:
				e.moveCursorBackward()
			case termbox.KeyArrowRight:
				e.moveCursorForward()
			case termbox.KeyHome, termbox.KeyCtrlA:
				e.moveCursorToTop()
			case termbox.KeyEnd:
				e.moveCursorToEnd()
			case termbox.KeyCtrlK, termbox.KeyPgup:
				e.scrollToAbove()
			case termbox.KeyCtrlJ, termbox.KeyPgdn:
				e.scrollToBelow()
			case termbox.KeyCtrlW:
				e.deleteWordBackward()
			case termbox.KeyEsc:
				e.candidatemode = false
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
	cc, _ := jqrun(e.query.StringGet(), e.json, e.args)
	return strings.Split(cc, "\n")
}

func (e *Engine) setCandidates() {
	filter := e.query.StringGet()
	e.candidates = []string{}
	if strings.IndexAny(strings.TrimLeft(filter, " "), "|{} @(),") == -1 {
		// try to find suggestions since it seems to be a simple jq filter
		validUntilNow, next := e.query.StringSplitLastKeyword()
		if validUntilNow == "" {
			validUntilNow = "."
		}
		keys, err := jqrun(validUntilNow+" | keys", e.json, []string{"-c"})
		if err == nil {
			candidates := strings.Split(keys[1:len(keys)-1], ",")
			if len(candidates) > 0 && candidates[0][0] == '"' {
				// only suggest if keys are strings
				// filter out
				for _, cand := range candidates {
					if strings.HasPrefix(cand, `"`+next) {
						e.candidates = append(e.candidates, cand)
					}
				}
			}
		}
	}
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
	filter, _ := e.query.StringSplitLastKeyword()
	filter += "." + e.candidates[e.candidateidx]
	e.query.StringSet(filter)
	e.cursorOffsetX = len(filter)
	e.candidatemode = false
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
