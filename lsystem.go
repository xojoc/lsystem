// This package was written by xojoc (http://xojoc.pw)
// and is in the Public Domain do what you want with it.

/*Package lsystem implements the L-system rewriting system.

 L-systems are comonly used to draw fractals for an explanation see:
https://en.wikipedia.org/wiki/L-system
*/
package lsystem

import (
	"bytes"
	"github.com/xojoc/turtle"
	"image/color"
	"log"
	"strconv"
	"strings"
)

// LSystem keeps track of the state of the L-system.
type LSystem struct {
	t          *turtle.Turtle
	rules      map[rune]string
	operations map[rune]string
	stack      [][3]float64
}

// New generates a new L-system. Rules are the rewriting rules.
// Operations is the set of operations to perform for each symbol.
//
// List of operations:
//    push - Save x,y coordinates and angle on the stack.
//    pop - Load x,y coordinates and angle from the stack.
//    rotate N - Change the direction of the next drawing operation by N degrees.
//    move N - Move by N pixels without drawing.
//    draw C W L - Draw a line with color C (in #rrggbbaa notation), width W and long L pixels
func New(rules map[rune]string, operations map[rune]string) *LSystem {
	l := &LSystem{}
	l.t = turtle.New()
	l.rules = rules
	l.operations = operations
	l.stack = [][3]float64{}
	return l
}

func (l *LSystem) push() {
	x := l.t.X
	y := l.t.Y
	a := l.t.A
	l.stack = append(l.stack, [3]float64{x, y, a})
}
func (l *LSystem) pop() {
	x := l.stack[len(l.stack)-1][0]
	y := l.stack[len(l.stack)-1][1]
	a := l.stack[len(l.stack)-1][2]
	l.stack = l.stack[:len(l.stack)-1]
	l.t.X = x
	l.t.Y = y
	l.t.A = a
}

func parsef64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
func hex(str string) uint8 {
	i, err := strconv.ParseUint(str, 16, 8)
	if err != nil {
		log.Fatal(err)
	}
	return uint8(i)
}
func parseColor(str string) color.Color {
	str = str[1:]
	return color.RGBA{hex(str[:2]), hex(str[2:4]), hex(str[4:6]), hex(str[6:8])}
}

// Run applies the L-system rules i times starting from axiom.
func (l *LSystem) Run(axiom string, i int) {
	s := axiom
	for j := 0; j < i; j++ {
		var buf bytes.Buffer
		for k := 0; k < len(s); k++ {
			if v, ok := l.rules[rune(s[k])]; ok {
				buf.WriteString(v)
			} else {
				buf.WriteByte(s[k])
			}
		}
		s = buf.String()
	}
	for j := 0; j < len(s); j++ {
		if o, ok := l.operations[rune(s[j])]; ok {
			fields := strings.Fields(o)
			for q := 0; q < len(fields); q++ {
				switch fields[q] {
				case "push":
					l.push()
				case "pop":
					l.pop()
				case "rotate":
					q++
					l.t.Rotate(parsef64(fields[q]))
				case "move":
					q++
					l.t.PenUp()
					l.t.Move(parsef64(fields[q]))
					l.t.PenDown()
				case "draw":
					q++
					l.t.SetColor(parseColor(fields[q]))
					q++
					l.t.SetWidth(parsef64(fields[q]))
					q++
					l.t.Move(parsef64(fields[q]))
				default:
					log.Fatal("unknown operation: " + fields[q])
				}
			}
		}
	}
}

// Save saves the image produced after executing Run in the given file name.
// The file format is based on the extension. Currently only PNG is supported,
// with extension .png.
func (l *LSystem) Save(name string) error {
	return l.t.Save(name)
}
