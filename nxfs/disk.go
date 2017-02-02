package nxfs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

type FileDisk struct {
	file           *os.File
	Size           uint32
	SizeOfBlock    uint32
	FirstFreeBlock uint32
	LastFreeBlock  uint32
}

func NewFileDisk(size_of_disk, size_of_block uint32, file *os.File) (*FileDisk, error) {
	var i int = 0
	for true {
		if uint32(math.Pow(2, float64(i))) == size_of_block {
			break
		} else if uint32(math.Pow(2, float64(i))) > size_of_block {
			return nil, errors.New("Size of block must be power of 2")
		}
		i++
	}

	info, err := file.Stat()
	CheckError(err)
	os.Truncate(info.Name(), 0)

	sizeB := (size_of_disk) * 1024 * 1024
	blocksCount := (sizeB / (size_of_block))
	fmt.Println("bc32: ", blocksCount)
	file.Seek(0, 0)
	var j uint32 = 0
	for j = 0; j < blocksCount; j++ {
		store := make([]byte, 4)
		if j != blocksCount-1 {
			binary.LittleEndian.PutUint32(store, j+1)
		} else {
			binary.LittleEndian.PutUint32(store, 0)
		}
		file.Seek(int64(size_of_block-4), 1)
		file.Write(store)
	}

	file.Sync()

	disk := &FileDisk{
		file:           file,
		Size:           sizeB,
		SizeOfBlock:    size_of_block,
		FirstFreeBlock: 2,
		LastFreeBlock:  blocksCount,
	}

	err = saveDiskMetadata(disk)
	if err != nil {
		return nil, err
	}
	return disk, nil
}

func saveDiskMetadata(disk *FileDisk) error {
	disk.file.Seek(int64(disk.SizeOfBlock-4), 0)
	store := make([]byte, 4)
	binary.LittleEndian.PutUint32(store, 0)
	disk.file.Write(store)

	disk.file.Seek(0, 0)

	store = make([]byte, 4)
	binary.LittleEndian.PutUint32(store, disk.Size)
	_, err := disk.file.Write(store)
	if err != nil {
		return err
	}

	store = make([]byte, 4)
	binary.LittleEndian.PutUint32(store, disk.SizeOfBlock)
	_, err = disk.file.Write(store)
	if err != nil {
		return err
	}

	store = make([]byte, 4)
	binary.LittleEndian.PutUint32(store, disk.FirstFreeBlock)
	_, err = disk.file.Write(store)
	if err != nil {
		return err
	}

	store = make([]byte, 4)
	binary.LittleEndian.PutUint32(store, disk.LastFreeBlock)
	_, err = disk.file.Write(store)
	if err != nil {
		return err
	}

	return nil
}
