<img src="https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExdzF2OWRmNnJzem9wYmptNGw2bWMzZXNvdzNjYXh4dDY3Yjg5b3V4dyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/l2QEfG0Zk1Fl5Shck/giphy.gif" width="100%" height="256"/>

<h1 align="center">Goch8</h1>

A CHIP-8 emulator written in Go, built to explore the fundamentals of CPU emulation, virtual machine design, and low-level systems concepts. Rather than relying on existing emulator frameworks, this implementation models the CHIP-8 architecture directly including memory, registers, stack, program counter, timers, input handling, and a pixel-based display.

The goal of this project is educational: to gain hands-on experience with instruction fetch–decode–execute loops, opcode implementation, sprite rendering, input mapping, and timing synchronization, all within a clean and minimal Go codebase.

---

## What is CHIP-8?

[CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) is an interpreted language from the mid-1970s, originally designed to make game development easier on early microcomputers. It has a 64×32 monochrome display, 16 registers, 4096 bytes of memory, and 35 opcodes. It's the "hello world" of emulator development.

---

## File structure

```
Goch8/
├── main.go
├── go.mod
├── core/
│   ├── cpu.go          — registers, memory, opcode execution
│   ├── display.go      — 64×32 pixel framebuffer
│   └── keyboard.go     — 16-key hex pad state
└── renderer/
    ├── renderer.go     — renderer interface
    └── ebiten/
        └── ebiten.go   — ebitengine window, draw loop, input
```

The `core` package knows nothing about rendering. It owns the CPU state, framebuffer, and key state. The renderer reads from those and drives the loop. Swapping the renderer out for a different backend requires no changes to `core`.

---

## Tools

- [Ebitengine](https://ebitengine.org): A simple 2D game engine for Go that handles the window, draw loop, and input.

---

## Installation

You need Go 1.21+ installed.

Clone and run:

```bash
git clone https://github.com/9lynd/Goch8
cd Goch8
go mod tidy
go run . path/to/rom.ch8
```

Test ROMs to try first: [Chip-8 test suite by Timendus](https://github.com/Timendus/chip8-test-suite) — it tests opcodes individually and shows pass/fail on screen.

If you don't have any ROMs, download from this [repository](https://github.com/kripod/chip8-roms/tree/master/games), or from this [website](https://johnearnest.github.io/chip8Archive/?sort=platform#chip8)

> ⚠️
Goch8 targets the original CHIP-8 specification only. SUPER-CHIP, XO-CHIP, and other extensions are not supported. Sound is not implemented in the current version.
---

## Key mapping

CHIP-8 uses a 16-key hex pad. It maps to your keyboard like this:

```
CHIP-8      Keyboard
---------   ---------
1 2 3 C     1 2 3 4
4 5 6 D     Q W E R
7 8 9 E     A S D F
A 0 B F     Z X C V
```

---

## Guide

This [Video](https://www.youtube.com/watch?v=YtSgV3gY3fs) explains a lot about chip-8 emulator and how to build it, exploring its concepts briefly; Instructions, Sprites, Opcodes, and more.

> 📖 *[Manual Guide — a step-by-step guide to writing a CHIP-8 emulator](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#memmap)*  
