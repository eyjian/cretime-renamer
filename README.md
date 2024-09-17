# cretime-renamer
批量将文件名修改为文件创建时间的工具，适用于相机、手机导出的照片管理。Windows 请在 WSL 中运行。重命名会保持原文件的修改时间。如果文件名已经是 YYYYMMDDhhmmss 格式（不含后缀），则不会做重名。

# 参数说明

* **-dirs**

指定文件（如照片或视频）所在目录，多个目录间以逗号分割，如：-dirs="/mnt/e/Phone,/mnt/f/Phone,/mnt/g/Phone" 。

* **-create-year-dir**

为 true 表示根据文件的修改时间自动创建年份目录，默认为 false

* **-create-month-dir**

为 true 表示根据文件的修改时间自动创建月份目录，默认为 false。注意如果 -create-month-dir 为 true，则会强制 -create-year-dir 也为 true，即使已设置 -create-year-dir 为 false 。

* **-skip-date-dir**

是否跳过日期目录，仅当 -create-year-dir 或 -create-month-dir 为 true 时有作用。

* **-ignore-dirs**

用于指定忽略不处理的目录。

* **-suffixes**

用于指定需要处理的文件名后缀，多个后缀以逗号分割。默认为空，表示处理所有文件。如指定只处理 png 和 jpg 文件：-suffixes="png,jpg" 。

* **sibling-dir**

用于指定是否创建同级的年份目录，默认为 false。为 true 表示在文件的上一级目录创建年份目录。仅当 -create-year-dir 或 -create-month-dir 为 true 时有效。

# WSL

* **创建挂载目录**

```shell
sudo mkdir /mnt/e
```

* **挂载U盘或者移动硬盘**

```shell
sudo mount -t drvfs E: /mnt/e
```

* **卸载U盘或者移动硬盘**

```shell
sudo umount /mnt/e
```

`drvfs` 是 Windows Subsystem for Linux (WSL) 中的一个虚拟文件系统，用于将 Windows 驱动器（如 C:、D:、E: 等）挂载到 WSL 的文件系统中。`drvfs` 允许在 WSL 中直接访问 Windows 文件系统中的文件和目录。

- `-t drvfs` 指定了文件系统的类型为 `drvfs`。
- `G:` 是 Windows 中的一个驱动器字母，表示要挂载的 Windows 驱动器。
- `/mnt/e` 是在 WSL 中挂载该驱动器的目标路径。