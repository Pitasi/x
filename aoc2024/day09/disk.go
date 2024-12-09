package main

import (
	"fmt"
	"strings"
)

type Disk struct {
	blocks []*Block
}

func (d *Disk) write(v int, start, end uint) {
	for i := start; i < end; i++ {
		if d.blocks[i] == nil {
			d.blocks[i] = &Block{
				data: v,
			}
		} else {
			panic("block already written, can't overwrite")
		}
	}
}

func (d *Disk) checksum() int {
	total := 0
	for i, block := range d.blocks {
		if block == nil {
			continue
		}
		total += block.data * i
	}
	return total
}

func (d *Disk) String() string {
	var sb strings.Builder
	for _, block := range d.blocks {
		if block == nil {
			sb.WriteRune('.')
		} else {
			sb.WriteString(fmt.Sprintf("%d", block.data))
		}
	}
	return sb.String()
}

type Block struct {
	data int
}

func CompactByBlock(d *Disk) {
	lCur := 0
	rCur := len(d.blocks) - 1
	for {
		if lCur > rCur {
			break
		}

		emptyBlock, ok := findNextEmptyBlock(d, lCur)
		if !ok {
			break
		}
		nonEmptyBlock, ok := findLastNonEmptyBlock(d, rCur)
		if !ok {
			break
		}

		if emptyBlock >= nonEmptyBlock {
			break
		}

		// swap blocks
		d.blocks[emptyBlock], d.blocks[nonEmptyBlock] = d.blocks[nonEmptyBlock], d.blocks[emptyBlock]

		lCur = emptyBlock + 1
		rCur = nonEmptyBlock - 1
	}
}

func CompactByFile(d *Disk) {
	rCur := len(d.blocks) - 1
	for {
		fileStart, fileEnd, ok := findLastFile(d, rCur)
		if !ok {
			// no more files
			break
		}
		rCur = fileStart - 1
		size := fileEnd - fileStart + 1

		emptyStart, ok := findEmptyContigBlocks(d, size)
		if !ok {
			// this file can't be fitted anywhere else
			continue
		}

		if emptyStart >= fileStart {
			// this file can't be fitted anywhere else
			continue
		}

		// write file blocks
		for i := range size {
			d.blocks[emptyStart+i] = d.blocks[fileStart+i]
			d.blocks[fileStart+i] = nil
		}
	}
}

func findNextEmptyBlock(d *Disk, start int) (int, bool) {
	for i := start; i < len(d.blocks); i++ {
		if d.blocks[i] == nil {
			return i, true
		}
	}

	return 0, false
}

func findLastNonEmptyBlock(d *Disk, end int) (int, bool) {
	for i := end; i >= 0; i-- {
		if d.blocks[i] != nil {
			return i, true
		}
	}
	return 0, false
}

func findLastFile(d *Disk, end int) (int, int, bool) {
	fileEnd, ok := findLastNonEmptyBlock(d, end)
	if !ok {
		return 0, 0, false
	}

	fileStart := fileEnd
	fileID := d.blocks[fileEnd].data
	for fileStart >= 0 && d.blocks[fileStart] != nil && d.blocks[fileStart].data == fileID {
		fileStart--
	}

	return fileStart + 1, fileEnd, true
}

func findEmptyContigBlocks(d *Disk, size int) (int, bool) {
	var (
		start     int
		ongoing   bool
		foundSize int
	)

	for i := 0; i < len(d.blocks); i++ {
		if foundSize == size {
			return start, true
		}

		if d.blocks[i] != nil {
			ongoing = false
			foundSize = 0
		}

		if d.blocks[i] == nil && !ongoing {
			start = i
			foundSize = 1
			ongoing = true
			continue
		}

		if d.blocks[i] == nil && ongoing {
			foundSize++
		}
	}

	return 0, false
}
