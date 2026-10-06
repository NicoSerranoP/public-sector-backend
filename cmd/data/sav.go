package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

const (
	recordVariable   = 2
	recordValueLabel = 3
	recordLabelVars  = 4
	recordDocument   = 6
	recordExtension  = 7
	recordEnd        = 999

	subtypeLongNames = 13

	continuationType = -1
	headerSkipBytes  = 84
	slotSize         = 8

	commandSkip      = 0
	commandEnd       = 252
	commandLiteral   = 253
	commandSpaces    = 254
	commandMissing   = 255
	compressedFormat = 1
)

var (
	missingBits = math.Float64bits(-math.MaxFloat64)
	spacesBits  = binary.LittleEndian.Uint64([]byte("        "))
)

type variable struct {
	name   string
	width  int
	labels map[float64]string
}

type savFile struct {
	variables []variable
	rows      [][]string
}

type savParser struct {
	r       *bufio.Reader
	err     error
	bias    float64
	vars    []variable
	slotVar []int
	pending map[float64]string
}

func readSav(r io.Reader) (*savFile, error) {
	p := &savParser{r: bufio.NewReader(r)}
	p.readHeader()
	p.readDictionary()
	rows := p.readCases()
	if p.err != nil {
		return nil, p.err
	}
	return &savFile{variables: p.vars, rows: rows}, nil
}

func (p *savParser) readHeader() {
	if string(p.bytes(4)) != "$FL2" {
		p.fail(errors.New("not an SPSS .sav file"))
		return
	}
	p.skip(60)
	p.skip(4 + 4) // layout code, nominal case size
	if p.int32() != compressedFormat {
		p.fail(errors.New("only bytecode-compressed .sav files are supported"))
		return
	}
	p.skip(4 + 4) // weight index, case count
	p.bias = p.float64()
	p.skip(headerSkipBytes)
}

func (p *savParser) readDictionary() {
	for p.err == nil {
		switch recordType := p.int32(); recordType {
		case recordVariable:
			p.readVariable()
		case recordValueLabel:
			p.readValueLabels()
		case recordLabelVars:
			p.applyValueLabels()
		case recordDocument:
			p.skip(80 * int(p.int32()))
		case recordExtension:
			p.readExtension()
		case recordEnd:
			p.skip(4)
			return
		default:
			p.fail(fmt.Errorf("unknown record type %d", recordType))
		}
	}
}

func (p *savParser) readVariable() {
	width := int(p.int32())
	hasLabel := p.int32()
	missingCount := int(p.int32())
	p.skip(4 + 4) // print and write formats
	name := string(p.bytes(8))

	if hasLabel != 0 {
		p.skip(padTo4(int(p.int32())))
	}
	if missingCount != 0 {
		p.skip(slotSize * abs(missingCount))
	}

	if width == continuationType {
		p.slotVar = append(p.slotVar, continuationType)
		return
	}
	p.slotVar = append(p.slotVar, len(p.vars))
	p.vars = append(p.vars, variable{name: strings.TrimSpace(name), width: width})
}

func (p *savParser) readValueLabels() {
	count := int(p.int32())
	p.pending = make(map[float64]string, count)
	for range count {
		value := math.Float64frombits(binary.LittleEndian.Uint64(p.bytes(slotSize)))
		length := int(p.bytes(1)[0])
		label := string(p.bytes(length))
		p.skip(padTo8(1+length) - (1 + length))
		p.pending[value] = label
	}
}

func (p *savParser) applyValueLabels() {
	count := int(p.int32())
	for range count {
		index := p.slotVar[p.int32()-1]
		if p.vars[index].width == 0 {
			p.vars[index].labels = p.pending
		}
	}
}

func (p *savParser) readExtension() {
	subtype := p.int32()
	size := int(p.int32())
	count := int(p.int32())
	data := p.bytes(size * count)
	if subtype == subtypeLongNames {
		p.applyLongNames(string(data))
	}
}

func (p *savParser) applyLongNames(data string) {
	longNames := make(map[string]string)
	for pair := range strings.SplitSeq(data, "\t") {
		short, long, found := strings.Cut(pair, "=")
		if found {
			longNames[short] = long
		}
	}
	for i := range p.vars {
		if long, found := longNames[p.vars[i].name]; found {
			p.vars[i].name = long
		}
	}
}

func (p *savParser) readCases() [][]string {
	var rows [][]string
	slots := &slotReader{r: p.r, bias: p.bias, cursor: slotSize}
	for p.err == nil {
		row, ok := p.readCase(slots)
		if !ok {
			break
		}
		rows = append(rows, row)
	}
	return rows
}

func (p *savParser) readCase(slots *slotReader) ([]string, bool) {
	row := make([]string, len(p.vars))
	for i, v := range p.vars {
		raw := make([]uint64, 0, 1)
		for range slotCount(v.width) {
			bits, ok := slots.next()
			if !ok {
				if i == 0 && len(raw) == 0 {
					p.fail(slots.err)
					return nil, false
				}
				p.fail(io.ErrUnexpectedEOF)
				return nil, false
			}
			raw = append(raw, bits)
		}
		row[i] = formatValue(v, raw)
	}
	return row, true
}

func formatValue(v variable, raw []uint64) string {
	if v.width > 0 {
		buf := make([]byte, 0, len(raw)*slotSize)
		for _, bits := range raw {
			buf = binary.LittleEndian.AppendUint64(buf, bits)
		}
		return strings.TrimRight(string(buf[:v.width]), " ")
	}
	if raw[0] == missingBits {
		return ""
	}
	number := math.Float64frombits(raw[0])
	if label, found := v.labels[number]; found {
		return label
	}
	return strconv.FormatFloat(number, 'f', -1, 64)
}

func slotCount(width int) int {
	if width == 0 {
		return 1
	}
	return (width + slotSize - 1) / slotSize
}

// slotReader decodes the bytecode-compressed data section one 8-byte slot at a time.
type slotReader struct {
	r        *bufio.Reader
	bias     float64
	commands [slotSize]byte
	cursor   int
	err      error
	finished bool
}

//nolint:cyclop
func (s *slotReader) next() (uint64, bool) {
	for !s.finished {
		if s.cursor >= slotSize {
			if _, err := io.ReadFull(s.r, s.commands[:]); err != nil {
				s.finished = true
				s.err = err
				return 0, false
			}
			s.cursor = 0
		}
		command := s.commands[s.cursor]
		s.cursor++
		switch command {
		case commandSkip:
		case commandEnd:
			s.finished = true
			s.err = io.EOF
		case commandLiteral:
			var literal [slotSize]byte
			if _, err := io.ReadFull(s.r, literal[:]); err != nil {
				s.finished = true
				s.err = err
				return 0, false
			}
			return binary.LittleEndian.Uint64(literal[:]), true
		case commandSpaces:
			return spacesBits, true
		case commandMissing:
			return missingBits, true
		default:
			return math.Float64bits(float64(command) - s.bias), true
		}
	}
	return 0, false
}

func (p *savParser) fail(err error) {
	if p.err == nil && !errors.Is(err, io.EOF) {
		p.err = err
	}
}

func (p *savParser) bytes(n int) []byte {
	buf := make([]byte, n)
	if p.err == nil {
		_, p.err = io.ReadFull(p.r, buf)
	}
	return buf
}

func (p *savParser) skip(n int) {
	if p.err == nil {
		_, p.err = p.r.Discard(n)
	}
}

func (p *savParser) int32() int32 {
	return int32(binary.LittleEndian.Uint32(p.bytes(4)))
}

func (p *savParser) float64() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(p.bytes(slotSize)))
}

func padTo4(n int) int { return (n + 3) &^ 3 }

func padTo8(n int) int { return (n + 7) &^ 7 }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (f *savFile) header() []string {
	names := make([]string, len(f.variables))
	for i, v := range f.variables {
		names[i] = v.name
	}
	return names
}
