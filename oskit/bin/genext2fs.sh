#!/bin/bash

gunzip ${1}-initrd.cpio.gz
mkdir -p .tmp-siliconfast
pushd .tmp-siliconfast
sudo cpio -idum < ../${1}-initrd.cpio
popd

SIZE=`sudo du -s .tmp-siliconfast | cut --fields=1`
SIZE=`expr '(' '(' $SIZE / 1024 ')' + 4 ')' '*' 1024 `
#inode_counti=`expr '(' $SIZE / 4 ')'`
echo "SIZE = $SIZE"
echo "genext2fs -b $SIZE -d .tmp-siliconfast ${1}-initrd.ramdisk"
./bin/genext2fs_$(uname -m) -b $SIZE -d .tmp-siliconfast ${1}-initrd.ramdisk
#e2fsck -fy $DST
