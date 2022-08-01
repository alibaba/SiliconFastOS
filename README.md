# SifastOS: An Operating System for Silicon Validation

## Introduction

SifastOS is an operating system for OS-based silicon validation, which is small, fast, flexible in configuration, and has rich cmdlines. SifastOS can improve your development efficiency in silicon validation, ensuring the delivery of high-quality silicon.

SifastOS contains a small build system, you just need to add a small number of configurations, so you can build a specific silicon validation platform easily.

SifastOS supports OpenAnolis kernel and Upstream kernel, and we use it as the silicon validation operating system in the development of ARMv9 server chip.

## What is OS for silicon validation

Pre-silicon software validation is becoming the standard in the computing platform development lifecycle in the semiconductor industry. Enabling the software stack in pre-silicon has a huge impact on the quality of the software stack, reducing the number of hardware-software integration issues in post-silicon and improving the time-to-market for end-to-end compute platforms.

As hardware manufacturers are adding more new features, including new instruction sets, new memory technology, accelerators, I/O technology, hardware-based security features, hardware-based diagnostics, and telemetry technologies, to their CPUs or processors to meet the compute demand and improve efficiency. Operating systems for software and applications need to take advantage of these new hardware-enabled enhancements and enable them in the OS.

During pre-silicon validation, the developer needs to test for BIOS/UEFI, firmware, device drivers, operating system boot and diagnostics tools, sometimes even need to optimize the entire software stack for a specific middleware framework, workload, and application. All these processes need to run an operating system on virtual platforms (often FPGA-based hardware emulators).

There are more information in [Why SifastOS](docs/why_SifastOS.md)

## Feature

### Silicon validation framework and test toolsets

- **Bee framework**: silicon validation framework. See detail in [Beetest](docs/bee.md)
- **Munit**: asynchronous kernel unit test framework. See detail in [Munit](docs/munit.md)
- Acpi tool set（Developing）: powerful acpi analysis tool

### More commands and flexible configuration

- **More commands and more command options than busybox** and easy to install:

    SifastOS uses command-line tools provided by open-source packages, avoiding the problem of using busybox commands that are not rich enough.
- Customize the userspace toolset as needed to control the size of the rootfs:

    SifastOS uses kbuild to configure userspace tools, and developers can not only select the packages built into SifastOS but even individual commands.

    <img src=docs/interface.png width="50%">

- Default config to build your SifastOS quickly:

    SifastOS provides some default configurations in config/ to help developers build SifastOS quickly. It contains the smallest collection of command-line tools needed for silicon validation.

### OS boot time speedup

- Use SystemV instead of Systemd and simplify SystemV boot scripts to reduce OS boot time.
- Minimizes necessary command support:

    Silicon validation requires rich command-line tool support, but too many command-line tools will make the rootfs too large. SifastOS has a minimal set of command-line tools for silicon validation to balance command-line tools with file size.
- Minimal kernel driver support（Developing）:

    SifastOS removes some kernel modules that will not be used during the silicon validation and reserves key drivers such as ACPI, Smmu, PCIe, etc to reduce kernel boot time.
- PCIe init delay（Developing）:

    Initialize PCIe driver after Linux Shell terminal starts instead of in the kernel. Reduces the time required for waiting for the kernel to boot.
- Minimize share library:

    The size of the shared lib libraries also accounts for a large portion of the rootfs, and SifastOS provides only the lib libraries required by its command-line tools. This is used to control the size of the rootfs.

### More tools for chip test and validation (Planning)

- Performance benchmark
- Debug and trace tool

## Quick Start

You can use 'make menuconfig' to config your SifastOS, then use 'make' to build.

You will get a new directory named result in root directory, ./result has three subdirectory

    result/
    ├── qemu
    ├── rootfs
    └── src

./src storage the source code. ./rootfs storage a rootfs used by SifastOS. ./qemu storage a rootfs image and a kernel image, you can run start_qemu.sh to run SifastOS with qemu, then you can see

<img src=docs/SifastOS_logo.jpg>

SifastOS can build and install Bee automatically, you just need to choose bee in menuconfig, if you want to build Bee individually, please see [there](docs/bee.md)

## Licensing

[GPL-2.0 license](COPYING)

## Developer Manual

1. [Bee framework document](docs/bee.md)
2. [Munit framework document](docs/munit.md)
3. [How to contribute](CONTRIBUTE.md)

## Version control guidelines

The version number of the project is in the form of x.x.x, where x is a number.

## Roadmap

We will introduce some features for kernel boot speed up and more silicon validation tools in SifastOS 1.0.0.

Then, we plan to support silicon validation for RISC-V server chip in SifastOS 2.0.0. If you are also a developer for OS-based silicon validation, we are glad you contact us and develop SifastOS together.

## Maintainer

- Ruidong Tian \<tianruidong@linux.alibaba.com>
- Jiankang Chen \<jkchen@linux.alibaba.com>
- Shaolin Xie \<shawnless.xie@linux.alibaba.com>

## Contributor

- Zhuo Song \<zhuo.song@linux.alibaba.com>
- Baolin Wang \<baolin.wang@linux.alibaba.com>
- Liguang Zhang \<zhangliguang@linux.alibaba.com>
- Wen Cheng \<yinxuan_cw@linux.alibaba.com>
- Neng Chen \<nengchen@linux.alibaba.com>
- Hongliang Yan \<archeryan@linux.alibaba.com>
- Hongbo Yao \<yaohongbo@linux.alibaba.com>

## Golden Rule

Please keep it small and beautiful :)
