package parts

import "github.com/hellodoge/bonk/bonk"

type Assembler struct {
	buffer  []byte
	written []Bound
}

type Bound struct {
	Start uint32
	End   uint32 // excluding
}

func NewAssembler(capacity uint) *Assembler {
	return &Assembler{
		buffer:  make([]byte, capacity),
		written: nil,
	}
}

func (c *Assembler) AddPart(part bonk.Part) {
	copy(c.buffer[part.Offset:], part.Block)
	c.addBound(Bound{
		Start: part.Offset,
		End:   part.Offset + uint32(len(part.Block)),
	})
}

func (c *Assembler) TryToAssemble() []byte {
	if c.gotEntirePiece() {
		return c.buffer
	} else {
		return nil
	}
}

func (c *Assembler) addBound(bound Bound) {
	for i := range c.written {
		if c.written[i].End <= bound.Start || bound.End <= c.written[i].Start {
			if i+1 == len(c.written) {
				c.written[i] = bound
				return
			} else {
				c.written[i] = c.written[len(c.written)-1]
				c.written = c.written[:len(c.written)-1]
			}
		}
	}
	c.written = append(c.written, bound)
}

func (c *Assembler) gotEntirePiece() bool {
	return len(c.written) == 1 &&
		c.written[0].Start == 0 &&
		c.written[0].End == uint32(len(c.buffer))
}
