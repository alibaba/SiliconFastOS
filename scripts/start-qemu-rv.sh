
qemu-system-riscv64 -M virt -bios fw_jump.elf -kernel Image -append "rootwait root=/dev/vda ro" -drive file=rootfs.ext4,format=raw,id=hd0 -device virtio-blk-device,drive=hd0 -netdev user,id=net0 -device virtio-net-device,netdev=net0 -nographic
