package main

import (
	"./nxfs"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("./sample.dsk")
	nxfs.CheckError(err)
	disk, err := nxfs.NewFileDisk(1, 512, file)
	nxfs.CheckError(err)
	fmt.Println(disk.Size, disk.SizeOfBlock, disk.FirstFreeBlock, disk.LastFreeBlock)
}
