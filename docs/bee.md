# Bee Test Project

## Goal of Bee

The goal of bee is to establish an important framework for silicon validation.

## Feature of Bee

  Bee is a simple silicon validation framework at **OS level**, it has some advantages：

  1. It is small to run on a low-frequency chip;
  2. It is very easy to add a new case in Bee；
  3. You can run a single case, a case set, or all cases, which means that you can use Bee to verify Unit, IP, and Kit in silicon validation.

## Bee Framework Overview

### Build system

<img src=bee_framework.png width="50%">

The Bee framework uses the recursive compilation build system, mainly involving bee_kernel_case, bee_user_case, bee_tools, etc. Users only need to edit the Makefile based on the template to add their test cases to Bee framework. Details can be found at [bee_add_case](bee_add_case.md)；

### Runtime system

<img src=runtime_framework.png width="50%">

Bee includes kernel cases and user space cases, kernel cases are executed asynchronously using the Munit module, and user-space cases are executed separately according to the shell script. Each case feature is a minimum execution unit.

## Compile and install

enter the directory of Bee:

1. Find the help:

    ./build.sh -h

2. Compile the bee:

    ./build.sh -m /lib/modules/\`uname -r\`/build

3. Install the bee:
   - Install to the DEFAULT path in /Your path/Bee/build

      ./build.sh -i ""

   - Install the bee to another path eg. /Your path/

      ./build.sh -i "/Your path"

   - Clean the bee:

      ./build.sh -c

## Run Case

After compiling, you can run the cases by Bee;

1. Export Bee root path:

    export BEETOPDIR=/Your Path/build

    cd /Your Path && ls

 you can see the files: bin  ko  lib64

2. Run the case by bt cmd:

    cd /Your Path/bin/

 - Get help:

    ./bt -h

 - Insmod a kernel case driver module in /your path/ko

    ./bt -i modules.ko

 - You can get a cases list

    ./bt -l

 - You can run a case for kernel case:

    ./bt -k case_name_01

 - You can run all cases, now you don’t need to run bt -i to insmod the case ko:

    ./bt -a

# Add Case Feature

See [how to add case in bee](bee_add_case.md)
