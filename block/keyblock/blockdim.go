package keyblock

import "fmt"
import . "file-structures/block/byteslice"

const (
	RECORDS = 1 << iota
	POINTERS
	EXTRAPTR
	EQUAPTRS
	NODUP
)

type BlockDimensions struct {
	Mode         uint8
	BlockSize    uint32
	KeySize      uint32
	PointerSize  uint32
	RecordFields []uint32
	record_size  uint32
}

func calcRecordSize(fields []uint32) uint32 {
	sum := uint32(0)
	for _, v := range fields {
		sum += v
	}
	return sum
}

func NewBlockDimensions(Mode uint8, BlockSize, KeySize, PointerSize uint32, RecordFields []uint32) (*BlockDimensions, bool) {
	dim := &BlockDimensions{
		Mode, BlockSize, KeySize, PointerSize, RecordFields,
		calcRecordSize(RecordFields)}
	if !dim.Valid() {
		return nil, false
	}
	return dim, true
}

func (self *BlockDimensions) NewRecord(key ByteSlice) *Record {
	return newRecord(key, self)
}

func (self *BlockDimensions) KeysPerBlock() int {
	var n int
	if self.Mode&(POINTERS|EQUAPTRS) == (POINTERS | EQUAPTRS) {
		n = int((self.BlockSize - BLOCKHEADER) /
			(self.KeySize + self.PointerSize))
	} else if self.Mode&EXTRAPTR == (EXTRAPTR) {
		n = int((self.BlockSize - self.PointerSize - BLOCKHEADER) /
			(self.RecordSize() + self.KeySize))
	} else {
		n = int((self.BlockSize - self.PointerSize - BLOCKHEADER) /
			(self.RecordSize() + self.KeySize + self.PointerSize))
	}
	return n
}

func (self *BlockDimensions) RecordSize() uint32 {
	return self.record_size
}

func (self *BlockDimensions) Valid() bool {
	if self.KeySize <= 0 {
		return false
	}
	switch self.Mode {
	case RECORDS, RECORDS | NODUP:
		if self.RecordSize() > 0 && self.PointerSize == 0 &&
			self.BlockSize >= self.RecordSize()+self.KeySize+BLOCKHEADER {
			return true
		} else {
			return false
		}
	case POINTERS, POINTERS | NODUP:
		if self.PointerSize > 0 && self.RecordSize() == 0 &&
			self.BlockSize >= (2*self.PointerSize)+self.KeySize+BLOCKHEADER {
			return true
		} else {
			return false
		}
	case POINTERS | EQUAPTRS, POINTERS | EQUAPTRS | NODUP:
		if self.PointerSize > 0 && self.RecordSize() == 0 &&
			self.BlockSize >= self.PointerSize+self.KeySize+BLOCKHEADER {
			return true
		} else {
			return false
		}
	case RECORDS | EXTRAPTR, RECORDS | EXTRAPTR | NODUP:
		if self.RecordSize() > 0 && self.PointerSize > 0 &&
			self.BlockSize >= self.PointerSize+self.RecordSize()+self.KeySize+BLOCKHEADER {
			return true
		} else {
			return false
		}
	case RECORDS | POINTERS, RECORDS | POINTERS | NODUP:
		if self.RecordSize() > 0 && self.PointerSize > 0 &&
			self.BlockSize >= (2*self.PointerSize)+self.RecordSize()+self.KeySize+BLOCKHEADER {
			return true
		} else {
			return false
		}
	case EXTRAPTR | RECORDS | POINTERS, EXTRAPTR, POINTERS | EXTRAPTR:
		return false
	}
	return false
}

func (self *BlockDimensions) String() string {
	return fmt.Sprintf(
		"Dimensions{Mode = %v, BlockSize = %v, KeySize = %v, PointerSize = %v, RecordFields = %v, KeysPerBlock=%v}",
		self.Mode, self.BlockSize, self.KeySize, self.PointerSize, self.RecordFields, self.KeysPerBlock())
}
