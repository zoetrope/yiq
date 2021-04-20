package yiq

import (
	"regexp"
	"strings"
)

type Query struct {
	query *[]rune
}

func NewQuery(query []rune) *Query {
	q := &Query{
		query: &[]rune{},
	}
	_ = q.Set(query)
	return q
}
func (q *Query) Get() []rune {
	return *q.query
}

func (q *Query) Set(query []rune) []rune {
	q.query = &query
	return q.Get()
}

func (q *Query) Insert(query []rune, idx int) []rune {
	qq := q.Get()
	if idx == 0 {
		qq = append(query, qq...)
	} else if idx > 0 && len(qq) >= idx {
		_q := make([]rune, idx+len(query)-1)
		copy(_q, qq[:idx])
		qq = append(append(_q, query...), qq[idx:]...)
	}
	return q.Set(qq)
}

func (q *Query) StringInsert(query string, idx int) string {
	return string(q.Insert([]rune(query), idx))
}

func (q *Query) Add(query []rune) []rune {
	return q.Set(append(q.Get(), query...))
}

func (q *Query) Delete(i int) []rune {
	qq := q.Get()
	lastIdx := len(qq)
	if i < 0 {
		if lastIdx+i >= 0 {
			qq = qq[0 : lastIdx+i]
		} else {
			qq = qq[0:0]
		}
	} else if i == 0 {
		qq = qq[1:]
	} else if i > 0 && i < lastIdx {
		qq = append(qq[:i], qq[i+1:]...)
	}
	return q.Set(qq)
}

func (q *Query) Clear() []rune {
	return q.Set([]rune(""))
}

func (q *Query) GetKeywords() [][]rune {
	query := string(*q.query)

	if query == "" {
		return [][]rune{}
	}

	splitQuery := strings.Split(query, ".")
	lastIdx := len(splitQuery) - 1

	keywords := [][]rune{}
	for i, keyword := range splitQuery {
		if keyword != "" || i == lastIdx {
			re := regexp.MustCompile(`\[[0-9]*\]?`)
			matchIndexes := re.FindAllStringIndex(keyword, -1)
			if len(matchIndexes) < 1 {
				keywords = append(keywords, []rune(keyword))
			} else {
				if matchIndexes[0][0] > 0 {
					keywords = append(keywords, []rune(keyword[0:matchIndexes[0][0]]))
				}
				for _, matchIndex := range matchIndexes {
					k := keyword[matchIndex[0]:matchIndex[1]]
					keywords = append(keywords, []rune(k))
				}
			}
		}
	}
	return keywords
}

func (q *Query) GetLastKeyword() []rune {
	keywords := q.GetKeywords()
	if l := len(keywords); l > 0 {
		return keywords[l-1]
	}
	return []rune("")
}

func (q *Query) StringGetLastKeyword() string {
	return string(q.GetLastKeyword())
}

func (q *Query) PopKeyword() ([]rune, []rune) {
	var keyword []rune
	var lastSepIdx int
	var lastBracketIdx int
	qq := q.Get()
	for i, e := range qq {
		if e == '.' {
			lastSepIdx = i
		} else if e == '[' {
			lastBracketIdx = i
		}
	}

	if lastBracketIdx > lastSepIdx {
		lastSepIdx = lastBracketIdx
	}

	keywords := q.GetKeywords()
	if l := len(keywords); l > 0 {
		keyword = keywords[l-1]
	}
	query := q.Set(qq[0:lastSepIdx])
	return keyword, query
}

func (q *Query) StringGet() string {
	return string(q.Get())
}

func (q *Query) StringSet(query string) string {
	return string(q.Set([]rune(query)))
}

func (q *Query) StringAdd(query string) string {
	return string(q.Add([]rune(query)))
}

func (q *Query) StringGetKeywords() []string {
	var keywords []string
	for _, keyword := range q.GetKeywords() {
		keywords = append(keywords, string(keyword))
	}
	return keywords
}

func (q *Query) StringPopKeyword() (string, []rune) {
	keyword, query := q.PopKeyword()
	return string(keyword), query
}

func (q *Query) StringSplitLastKeyword() (string, string) {
	filter := q.StringGet()
	last := strings.LastIndexAny(filter, "].")
	if last == -1 {
		return strings.TrimSpace(filter), ""
	}
	switch filter[last] {
	case '.':
		return filter[:last], filter[last+1:]
	case ']':
		return filter[:last+1], filter[last+1:]
	}
	return filter, ""
}
