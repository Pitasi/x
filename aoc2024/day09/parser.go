package main

func Parse(input []byte) *Disk {
	p := &Parser{
		input: input,
	}
	return p.parse()
}

type Parser struct {
	input []byte
}

func (p *Parser) parse() *Disk {
	blocksCount := p.countBlocks()
	disk := &Disk{
		blocks: make([]*Block, blocksCount),
	}

	var cur uint
	for i := 0; i < len(p.input); i++ {
		// 2333133121414131402
		length := p.getValue(i)
		if i%2 == 0 {
			// write file
			start := cur
			cur += length
			end := cur
			fileID := i / 2
			disk.write(fileID, start, end)
		} else {
			// skip empty blocks
			cur += length
		}
	}

	return disk
}

func (p *Parser) countBlocks() uint {
	var sum uint
	for i := range p.input {
		sum += p.getValue(i)
	}
	return sum
}

func (p *Parser) getValue(idx int) uint {
	return uint(p.input[idx] - '0')
}
