package ngxlog

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"time"
)

type Scanner struct {
	// 格式必须 空格分隔， 如有空格，则必须用""、[]包裹
	format  string
	scanner *bufio.Scanner

	// header information
	nFields  int
	fieldCol map[string]int

	// current dealing Record(the parsed line)
	rec *Record
}
type Record struct {
	scanner *Scanner
	indexes [][2]int // 指明怎么从bLine中截取字段值
	bLine   []byte   // []byte line 变化（GC）
}

func (s *Scanner) Scan() bool {
	rec := s.rec
	if !s.scanner.Scan() {
		if s.scanner.Err() == bufio.ErrTooLong { // token too long
			rec.bLine = rec.bLine[:0]
			rec.indexes = rec.indexes[:0]
			return true
		}
		return false
	}

	rec.bLine = s.scanner.Bytes()
	if len(rec.bLine) == 0 { // skip empty line
		rec.indexes = rec.indexes[:0]
		return true
	}
	rec.indexes = parseLine(rec.bLine, s.nFields)

	return true
}

func (s *Scanner) parseFieldCol() {
	s.fieldCol = make(map[string]int)
	format := s2b(s.format)
	indexes := parseLine(format, 0)
	for i, index := range indexes {
		field := format[index[0]:index[1]]
		if field[0] == '$' { // strip $
			field = field[1:]
		}
		s.fieldCol[b2s(field)] = i
	}
	s.nFields = len(indexes)
}

func (s *Scanner) Record() *Record {
	return s.rec
}

func (rec *Record) Bytes() []byte {
	return rec.bLine
}

func (rec *Record) Text() string {
	return b2s(rec.bLine)
}

func (rec *Record) Map() map[string]string {
	m := make(map[string]string)
	for field, col := range rec.scanner.fieldCol {
		m[field] = rec.Col(col)
	}
	return m
}

func (rec *Record) Field(name string) string {
	return rec.Col(rec.scanner.fieldCol[name])
}

func (rec *Record) FieldTime(name string) (time.Time, error) {
	return time.Parse(NgxTimeLocalCommonLogFormat, rec.Field(name))
}

func (rec *Record) FieldInt(name string) (int, error) {
	return strconv.Atoi(rec.Field(name))
}

func (rec *Record) FieldFloat(name string) (float64, error) {
	return strconv.ParseFloat(rec.Field(name), 64)
}

// Col zero-based
func (rec *Record) Col(i int) string {
	return b2s(rec.bLine[rec.indexes[i][0]:rec.indexes[i][1]])
}

func (rec *Record) Mismatch() bool {
	return rec.scanner.nFields != len(rec.indexes)
}

// 只存 offset
func parseLine(bLine []byte, nFields int) [][2]int {
	indexes := make([][2]int, 0, nFields)
	lineLen := len(bLine)
	for i := 0; i < lineLen; i++ {
		if bLine[i] == ' ' {
			continue
		}
		//  start 第一个非空字符，end 结束字符（Space、[、"、行尾）
		start, end, sc := i, i, bLine[i]
		for end = i; end < lineLen; end++ {
			if bLine[end] == ' ' {
				if sc == '"' || sc == '[' { // 在界定符之内的空格
					continue
				}
				break // 定界符之外的空格
			} else { // 非空字符，识别是否为界定符
				if bLine[end] == '"' && sc == '"' && end > i {
					start++
					break
				} else if bLine[end] == ']' && sc == '[' {
					start++
					break
				}
			}
		}

		indexes = append(indexes, [2]int{start, end})
		//fmt.Println(string(bLine[start:end]), i, start, end)

		if end == lineLen {
			break
		}
		i = end // end < lineLen
	}

	return indexes
}

func NewScanner(format string, reader io.Reader) *Scanner {
	s := new(Scanner)
	s.format = format
	s.scanner = bufio.NewScanner(reader)
	s.parseFieldCol()
	s.rec = new(Record)
	s.rec.scanner = s

	return s
}

func LogMismatch(rec *Record) {
	log.Printf("Values len %d mismatch fields len %d, raw text: %s \n", len(rec.indexes), rec.scanner.nFields, rec.Text())
}
