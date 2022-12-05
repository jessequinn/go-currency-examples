# Introduction to Concurrency

### CPU Basics

While the CPU is considered the most important unit in a machine,
the CPU does quite a lot of things on its own, but it can't do it all by itself.
The CPU works in a direct relationship with the **RAM** and **Storage** the user's machine has.

Thus, in many times your concurrent programs, the CPU might reach a limit not necessarily
because it is too busy, or the workload is too high, it might be caused because of **limited
resources** on RAM or Storage side, and their quality.

Of course many things will depend on your CPU, things like how many CORES it has.
Normally when picking up a machine, the number of CORES a CPU has will in many cases
dictate its performance.

CPU Cores are of 2 types: **Physical** and **Logical**. If **HYPER-THREADING** is enabled,
usually your **LOGICAL CORES = double the PHYSICAL CORES**.  

### OSX `sysctl`

Here are a couple of commands that will help you inspect
the CPU resources if you happen to own a Mac.

Similar commands are usually available on Windows and Linux as well.

```bash
# get the number of logical CPU cores
sysctl hw.logicalcpu

# get the number of physical CPU cores
sysctl hw.physicalcpu

# get the number of logical cores
sysctl hw.ncpu

# get the number of physical/logical cores
# also thread count meaning the total count of running threads in parallel
sysctl -a | grep machdep.cpu | grep count
```

### Threads

Tasks inside our machines are usually ***organized as threads***, or a list of tasks that ***belong*** to the ***same
context*** which can be picked up by any **CPU Core**.

Important to know is a **Logical CORE** can only execute **1 thread at a time**,
meaning if there is any **blocking task** part of the same thread, and the CORE
gets to execute that task, **the entire thread will block**, causing the ***rest***
of the tasks to be ***blocked*** as well.

Threads inside our machines are of 2 kinds: **User Level Threads** and **Kernel Level Threads**.
Every process that requires access to either of: ***CPU***, ***RAM***, ***Storage***, ***Network***
and other types of resources, usually goes through the middle man aka the **OS Kernel**.

***Kernel Level Threads*** are fully ***managed*** by the ***Kernel*** and there's not much or nothing
influencing them or, the way they would behave.

The control comes with ***User Level Threads***. Each user inside an Operating System
respectively runs a certain amount of ***Applications***, including Go apps.
Applications respectively can **create/destroy** their own ***threads***.

### Concurrency vs Parallelism

***Concurrency*** and ***Parallelism*** are terms which are quite confused.
Most of the times developers speak of it as if it's 1 and the same thing.

To make things as clear and simple as possible think of it this way:

"**Concurrency** is **DEALING** with a lot of things at the same time.
**Parallelism** is DOING a lot of things at the same time."

**Concurrency** is a property of the **code**. **Parallelism** is
a property of the **running program**.

### Context Switching

Under every layer of abstraction in the end all tasks **scheduled** to be executed on a **specific
thread** will end up being picked by a **CPU Core** to be executed. The reality is
in the modern concurrent world we never know ***how much*** time a task inside a thread will take.

The goal is to try, and ***complete*** all the tasks in the most performant way possible,
***without wasting unnecessary time*** on CPU. The best way modern CPUs achieve this is by
using a known process called **context switching**.

In very Layman's terms, context switching is when in the **context** of a **thread**,
the **CPU CORE** has some kind of internal **ticker**, which goes through each task available
inside the thread, ***trying to execute*** it for a ***limited burst time***, if the task is
***taking too long***, there might be a task which takes a smaller amount of time, causing
the CORE to ***context switch*** aka move to another task.

The ticker time is usually a pretty small amount of time, so that the CPU CORE
can check for tasks availability pretty fast.

As stated before, the **CORE** will **block** if a certain task contains **blocking
work**, such as **I/O** or **SysCall(s)**, which usually involves some kind of
***mutual exclusion*** or ***blocking operation***, thus causing the ***entire thread*** to be blocked.

Bottom line: **Concurrency** is pretty much **context switching** and dealing with many
tasks at the same time, while **Parallelism** is executing things in parallel using
the Concurrency principle aka **context switching in parallel**.

### Go Application Overview

The picture with the CPU and its companion resources, the Operating System, the OS Kernel
and a bunch of threads containing tasks pretty much, remains unchanged.

What has changed is the way the **binary** generated by Go works. Compared to other languages,
a ***Go binary*** is **slightly bigger** because it also contains other things, other than the code
you'll end up executing.

All Go binaries when generated they also contain the **Garbage Collector** responsible for
***memory management*** inside your Go application, and as well another very important component
which is the **Go Scheduler**, which is responsible for ***scheduling*** tasks (normally functions)
on a specific thread, which will again gets created for you before even running your binary.

And you guessed it, both the ***Go Scheduler*** and **Garbage Collector*** will actively communicate
with the **OS Kernel**, because they'll frequently make use of your machine's resources.

### Fork Join Model

While the `main` function seems only like the **entry point** for every **Go binary**. What else happens,
is the fact that every `main` function runs as a **Go Routine** aka the **Main Go Routine**,
which automatically makes your simplest binary **run concurrently**.

In Go, you're not limited to **how many** go routines you can run or how many **nested** go routines you can have.
Thus, every go routine which **derives** from the main go routine is considered a **child go routine**.

With child go routines, comes management of these **child processes**. Normally the process that
creates a child process is responsible for **destroying/waiting** on that process.

If you ever ran a go routine inside main without any kind of **waiting mechanism** other than `time.Sleep`,
you may have probably wondered why does it exit without waiting on the child process to finish?

This is all due to how **Concurrency** is done in Go, namely in Go Concurrency is done in
conformance with something called **FORK JOIN** model.

Every time your Go program encounters the `go` keyword it automatically creates
a **FORK** from the ***parent process***, thus causing the ***Go Scheduler*** to go ahead and schedule that task.
For ***every FORK point***, there must be a **JOIN** point,
meaning every time a **parent process** created a **child process**,
it either **must wait** on it to finish or have some kind of **cancellation/cleanup process**
for all its **child processes**, in order to avoid any kind of **leaks**.

TLDR; **FORK** points are automatically created when using the `go` keyword,
**JOIN** points can be created using a `sync.WaitGroup` or a `channel`.

***Note***: Not to be confused, the term **process** here does not refer to
an **OS process**, which is way more expensive.

### Common Concurrency Issues

Here are couple of common issues, which Concurrent Code usually faces:

𐄂 **RACE CONDITIONS**

When **multiple concurrent** operations try to **read/write shared data** **at the same time**,
thus causing **non-deterministic** results.

𐄂 **DEADLOCKS**

When concurrent operations are **protected** by some kind of **lock** / **mutual exclusion**,
making each process involved in **waiting forever** on one another, causing a dead lock.

𐄂 **LIVELOCKS**

When multiple **concurrent processes**, **pretend** they **modify shared data**, which is **protected**
by some kind of **lock**, when in reality they just end up acquiring a lock **without changing the state**.

𐄂 **STARVATION**

When **1 of multiple concurrent processes** involved **abuses the CPU**, thus causing other processes
**waiting** for their turn to **starve**, hence the name.

𐄂 **CODE COMPLEXITY**

Writing concurrent code does not always look as normal synchronous code. Sometimes complexity
**grows naturally** just because of the **concurrent complex nature**, thus requiring **shifting the code design**.

### Zip Archives

- [Concurrency in Go #1 - Introduction to Concurrency](https://youtu.be/_uQgGS_VIXM) - [[Download Zip]](https://github.com/golang-basics/concurrency/raw/master/archives/concurrency-1.tar.gz)

### Presentations

- [Concurrency in Go #1 - Introduction to Concurrency](https://github.com/golang-basics/concurrency/raw/master/presentations/1_introduction-to-concurrency)

### Examples

- [Synchronous Tasks](https://github.com/golang-basics/concurrency/blob/master/intro/sync-tasks/main.go)
- [Asynchronous Tasks](https://github.com/golang-basics/concurrency/blob/master/intro/async-tasks/main.go)
- [Asynchronous Tasks - Fixed](https://github.com/golang-basics/concurrency/blob/master/intro/fixed-async-tasks/main.go)
- [Fork Join - No Join Point](https://github.com/golang-basics/concurrency/blob/master/intro/fork-join/no-join-point/main.go)
- [Fork Join - WaitGroup Join Point](https://github.com/golang-basics/concurrency/blob/master/intro/fork-join/wg-join-point/main.go)
- [Fork Join - Channel Join Point](https://github.com/golang-basics/concurrency/blob/master/intro/fork-join/channel-join-point/main.go)

### Resources

- [Fork-Join Model - Wiki](https://en.wikipedia.org/wiki/Fork%E2%80%93join_model)
- [Hyper-Threading - Wiki](https://en.wikipedia.org/wiki/Hyper-threading)
- [Threads - Wiki](https://en.wikipedia.org/wiki/Thread_(computing))
- [Concurrency is not Parallelism](https://blog.golang.org/waza-talk)
- [Concurrent Computing - Wiki](https://en.wikipedia.org/wiki/Concurrent_computing)
- [Parallel Computing - Wiki](https://en.wikipedia.org/wiki/Parallel_computing)
- [Context Switching - Wiki](https://en.wikipedia.org/wiki/Context_switch#:~:text=In%20computing%2C%20a%20context%20switch,of%20a%20multitasking%20operating%20system.)
- [Race Condition - Wiki](https://en.wikipedia.org/wiki/Race_condition#:~:text=A%20race%20condition%20or%20race,the%20possible%20behaviors%20is%20undesirable.)
- [Deadlock - Wiki](https://en.wikipedia.org/wiki/Deadlock)
- [Livelock - Wiki](https://en.wikipedia.org/wiki/Deadlock#Livelock)
- [Livelock - StackOverflow](https://stackoverflow.com/questions/6155951/whats-the-difference-between-deadlock-and-livelock)
- [Starvation - Wiki](https://en.wikipedia.org/wiki/Starvation_(computer_science))

[Home](https://github.com/golang-basics/concurrency)