package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
)

// op is a 3-bit opcode
type op byte

const (
	adv op = 0b000
	bxl op = 0b001
	bst op = 0b010
	jnz op = 0b011
	bxc op = 0b100
	out op = 0b101
	bdv op = 0b110
	cdv op = 0b111
)

func (o op) String() string {
	switch o {
	case adv:
		return "adv"
	case bxl:
		return "bxl"
	case bst:
		return "bst"
	case jnz:
		return "jnz"
	case bxc:
		return "bxc"
	case out:
		return "out"
	case bdv:
		return "bdv"
	case cdv:
		return "cdv"
	default:
		return "???"
	}
}

// M is the machine with registers A, B, and C and that can run a program.
type M struct {
	A, B, C int

	program []op
	ip      int

	Out io.ByteWriter
	Err string
}

func (m *M) SetRegs(a, b, c int) {
	m.A = a
	m.B = b
	m.C = c
}

func (m *M) Load(program []byte) {
	m.program = make([]op, len(program))
	for i, b := range program {
		m.program[i] = op(b)
	}
	m.ip = 0
	m.Err = ""
}

func (m *M) Run() {
	for m.ip < len(m.program) {
		op := m.program[m.ip]

		// load operand
		operand := m.program[m.ip+1]

		var combo int
		switch operand {
		case 0, 1, 2, 3:
			combo = int(operand)
		case 4:
			combo = m.A
		case 5:
			combo = m.B
		case 6:
			combo = m.C
		case 7:
			combo = -1
		}

		switch op {
		case adv:
			num := m.A
			if combo-1 < 0 {
				m.Err = "adv: combo-1 < 0"
				return
			}
			m.A = num >> combo
			m.ip += 2
		case bxl:
			m.B = m.B ^ int(operand)
			m.ip += 2
		case bst:
			m.B = combo & 0b111
			m.ip += 2
		case jnz:
			if m.A != 0 {
				m.ip = int(operand)
			} else {
				m.ip += 2
			}
		case bxc:
			m.B = m.B ^ m.C
			m.ip += 2
		case out:
			m.B = combo & 0b111
			err := m.Out.WriteByte(byte(m.B))
			if err != nil {
				m.Err = err.Error()
				return
			}
			m.ip += 2
		case bdv:
			num := m.A
			if combo-1 < 0 {
				m.Err = "bdv: combo-1 < 0"
				return
			}
			m.B = num >> combo
			m.ip += 2
		case cdv:
			num := m.A
			if combo-1 < 0 {
				m.Err = "cdv: combo-1 < 0"
				return
			}
			m.C = num >> combo
			m.ip += 2
		default:
			panic("invalid opcode")
		}
	}
}

//go:embed input.txt
var input []byte

type w struct {
	expected []byte
	got      []byte
	cur      int
}

func (w *w) Done() bool {
	return w.cur == len(w.expected)
}

func (w *w) Count() int {
	return w.cur
}

func (w *w) Reset() {
	w.cur = 0
	w.got = nil
}

func (w *w) WriteByte(p byte) (err error) {
	if w.expected[w.cur] != p {
		return errors.New("expected " + string(w.expected[w.cur]) + ", got " + string(p))
	}
	w.cur++
	w.got = append(w.got, p)
	return
}

func (w *w) Bytes() []byte {
	return w.got
}

func main() {
	_, b, c, program := Parse(input)
	//writer := &w{expected: program}
	buf := bytes.Buffer{}
	m := M{Out: &buf}

	potentials := make(map[int]struct{})
	for i := 0; i < 0b111; i++ {
		potentials[i] = struct{}{}
	}

	for cur := 0; cur < len(program); cur++ {
		newPotentials := make(map[int]struct{})

		for pot := range potentials {
			for i := 0; i < 0b1111111111; i++ {
				a := (i << (3 * cur)) // | (pot & (0b111 << (i * cur)))
				for j := 0; j < cur; j++ {
					a |= (pot & (0b111 << (j * 3)))
				}
				buf.Reset()
				m.SetRegs(a, b, c)
				m.Load(program)
				m.Run()

				output := buf.Bytes()
				if m.Err != "" || cur >= len(output) {
					continue
				}

				ok := true
				for idx := range cur + 1 {
					if output[idx] != program[idx] {
						ok = false
						break
					}
				}
				if ok {
					newPotentials[a] = struct{}{}
				}
			}
		}

		potentials = newPotentials
		if len(potentials) == 0 {
			break
		}
	}

	var min int
	for pot := range potentials {
		if min == 0 || pot < min {
			min = pot
		}
	}
	fmt.Println(min)
}

func star1() {
	var out bytes.Buffer
	a, b, c, program := Parse(input)
	m := M{Out: &out}
	m.Load(program)
	m.SetRegs(a, b, c)
	m.Run()

	printOut(out.Bytes())
}

func printOut(bytes []byte) {
	for i, b := range bytes {
		fmt.Printf("%c", b+'0')
		if i != len(bytes)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println()
}
