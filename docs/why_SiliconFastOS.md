# Why need silicon validation at OS level?

1. Help you integrate and iterate the discrete test cases together.
2. There are some differences between UEFI and OS. At the OS level, we can find the problems that can't be found in UEFI mode.
3. Testing on OS is more similar to the actual production environment.
4. If we want to guarantee the function quality of our self-developed peripherals and modules, we can only fully pass the detailed test under OS to ensure the quality of hardware modules.

## Differences in silicon validation between hardware firmware and OS

|                                   | Hardware | Firmware |  OS |
|  ----                             | ----     |  ----    | --- |
| HW/SW co-design                   | &#10006; |  &#10006; | &#10004; |
| Flexibility                       | &#10006; |  &#10006; | &#10004; |
| Mutli-IP co-design                |   kit    |  IP       | Unit     |
| Similarity to real e2e environment| &#10006; |  &#10006; | &#10004; |
| Ecosystem                         | &#10006; |  &#10006; | &#10004; |
| Chip Simulation                   | &#10004; |  &#10006; | &#10006; |
| Chip debug                        | &#10004; |  &#10006; | &#10006; |
| HW/FW Interaction                 | &#10006; |  &#10004; | &#10006; |
| Firmware Debug                    | &#10006; |  &#10004; | &#10006; |

## Why we need a new OS for silicon validation

The main difference between FPGA-based hardware emulators and real hardware is that the emulator just run at MHz speed, so there are some limitations to selecting the OS which will run on the emulator.

1. OS needs to be small, otherwise, it will run slowly. So we can not use the Linux distribution like Ubuntu, Fedora, etc.
2. OS needs rich cmdline tools, especially when you are developing a server chip. In this respect, busybox and Embedded kernel is not a good choice either.
3. OS needs easy configuration, because your system may be a simple embedded chip but also a complex server chip. The OS needs to be config with different kernel/rootfs sizes and cmdline tools set.

We also need some frameworks and debug tools to speed up pre-silicon silicon validation process.

In summary, we need a new OS for silicon validation that is small, fast, and easy to customize.

## Difference between SiliconFastOS and busbox buildroot

Both Busybox and Buildroot are designed for embedded scenarios, and can't completely cover silicon validation scenarios.

In addition to this, silicon validation and embedded development have different requirements for OS. The Embedded system has simpler hardware. But the hardware for silicon validation is more complex and will change during the different stages frequently. Therefore, embedded development and silicon validation are two different domains, and software developed for embedded scenarios cannot be well applied in the silicon validation domain.

SiliconFastOS provides richer command-line tools and a proprietary testing framework

|               | Busybox | BuildRoot etc | SiliconFastOS |
|  ----         | ----    |  ----  | --- |
| OS speed up   | &#10004; | &#10006;  | &#10004; |
| Rich cmd      | &#10006;     | &#10004; | &#10004; |
| Native env    | &#10006;    | &#10006; | &#10004; |
| Flexibility   | &#10006;     | &#10004; | &#10004; |
| Out of box    | &#10006; | &#10004;| &#10004; |

