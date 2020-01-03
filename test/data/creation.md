# evidence
Sample data for forensics processing

Forensics software need to be able to parse and process many different file formats. This repository contains samples of different file formats that can be used to test forensics software.

## Create filesystems

``` sh
dd if=/dev/zero bs=1024 count=10000 of=fat16.dd
mkfs.fat -F 16 fat16.dd
mkdir -p /mnt/fat16
mount fat16.dd /mnt/fat16
cp -r filesystem_content/. /mnt/fat16
umount /mnt/fat16

dd if=/dev/zero bs=1024 count=40000 of=fat32.dd
mkfs.fat -F 32 fat32.dd
mkdir -p /mnt/fat32
mount fat32.dd /mnt/fat32
cp -r filesystem_content/. /mnt/fat32
umount /mnt/fat32

dd if=/dev/zero bs=1024 count=10000 of=ext2.dd
mkfs.ext2 ext2.dd
mkdir -p /mnt/ext2
mount ext2.dd /mnt/ext2
cp -r filesystem_content/. /mnt/ext2
umount /mnt/ext2

dd if=/dev/zero bs=1024 count=20000 of=ext3.dd
mkfs.ext3 ext3.dd
mkdir -p /mnt/ext3
mount ext3.dd /mnt/ext3
cp -r filesystem_content/. /mnt/ext3
umount /mnt/ext3

dd if=/dev/zero bs=1024 count=20000 of=ext4.dd
mkfs.ext4 ext4.dd
mkdir -p /mnt/ext4
mount ext4.dd /mnt/ext4
cp -r filesystem_content/. /mnt/ext4
umount /mnt/ext4

dd if=/dev/zero bs=1024 count=20000 of=ntfs.dd
mkfs.ntfs -F ntfs.dd
mkdir -p /mnt/ntfs
mount ntfs.dd /mnt/ntfs
cp -r filesystem_content/. /mnt/ntfs
umount /mnt/ntfs

7z a zip.zip filesystem_content/. -xr!*.DS_Store
7z a 7z.7z filesystem_content/. -xr!*.DS_Store
7z a tar.tar filesystem_content/. -xr!*.DS_Store
```