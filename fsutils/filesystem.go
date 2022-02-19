package fsutils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kensh1ro/willie/winapi"
)

func Cd(dir string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}

func Rm(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}
	return nil
}

func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return dir
}

func ListDrives() string {
	n := winapi.GetLogicalDriveStrings(0, nil)
	a := make([]byte, n)
	winapi.GetLogicalDriveStrings(n, &a[0])
	drives := strings.Split(strings.TrimRight(string(a), "\x00"), "\x00")
	var output strings.Builder
	for _, d := range drives {
		t := []byte(d)
		switch winapi.GetDriveType(&t[0]) {
		case 0:
			output.WriteString(fmt.Sprintf("%s\tDrive Unknown\n", d))
		case 1:
			output.WriteString(fmt.Sprintf("%s\tDrive not mounted\n", d))
		case 2:
			output.WriteString(fmt.Sprintf("%s\tRemovable Media\n", d))
		case 3:
			output.WriteString(fmt.Sprintf("%s\tHard Disk\n", d))
		case 4:
			output.WriteString(fmt.Sprintf("%s\tNetwork Drive\n", d))
		case 5:
			output.WriteString(fmt.Sprintf("%s\tCD-ROM\n", d))
		case 6:
			output.WriteString(fmt.Sprintf("%s\tRAM Disk\n", d))
		}
	}
	return output.String()
}

func ListDir(dir string) string {
	var output strings.Builder
	var s, _ = os.ReadDir(dir)
	//var _type string
	for _, file := range s {
		info, _ := file.Info()
		t := info.ModTime()

		/*		if info.IsDir() {
					_type = "dir"
				} else {
					_type = "fil"
				}*/

		output.WriteString(info.Mode().String())
		output.WriteByte(9)
		output.WriteString(strconv.FormatInt(info.Size(), 10))
		output.WriteByte(9)
		//output.WriteString(_type)
		output.WriteByte(9)
		output.WriteString(t.Format("2006-01-02 15:04:05"))
		output.WriteByte(9)
		output.WriteString(info.Name())
		output.WriteByte(10)
	}
	return output.String()
}
