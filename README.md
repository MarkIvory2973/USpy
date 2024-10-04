# USpy

Steal files from usb disks.

## Installation

```bash
git clone https://github.com/MarkIvory2973/USpy.git
```

## Usage

```bash
cd ProxyTest/src
go run ./main.go --help
go run ./main.go --scan-rules xls,xlsx,doc,docx,pdf --scan-level 15
```

## Parameters

|Parameter|Required|Default|Description|
|:-|:-:|:-|:-|
|--host|-|127.0.0.1|HTTP Server host|
|--port|-|6789|HTTP Server port|
|--scan-rules|-|ppt,pptx,xls,xlsx,doc,docx,pdf,txt,jpg,jpeg,png,bmp,gif|Scan rules|
|--scan-level|-|20|Scan level|
|--admin-name|-|Admin_USpy|Admin USB volume name|
|--temp-path|-|D:/USpy/|Temporary folder path|
