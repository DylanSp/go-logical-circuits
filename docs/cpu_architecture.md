# CPU Architecture

Overall architecture is based on MIPS.

## Word Size

32 bits.

- All registers are 32 bits
- All instructions are 32 bits
- Each memory address stores a 32-bit word.

## Registers

32 registers (`r0` - `r31`). Some of these are reserved for special purposes:

- `r0` is a dedicated `0` register - reading from it always produces `0`, writes to it are ignored.
- `pc` - program counter
  - (is this one of the 32 general-purpose registers, or is it a separate register that's not accessible to user programs?)
- `sp` - stack pointer
- (frame pointer?)
- (return address?)
- (condition flags register for reading results of comparison instructions and/or detecting overflow?)

## Instruction Set

The first 6 bits of an instruction (bits 31-26) are the opcode; the remaining bits differ between the different instruction formats.

3 formats:

- R-type (register) - operations on two registers, outputting to a third register.
  - 6 bits for opcode (bits 31-26) - always `000 000`
  - 5 bits to specify destination register (bits 25-21)
  - 5 bits to specify first operand register (bits 20-16)
  - 5 bits to specify second operand register (bits 15-11)
  - 5 bits to specify shift amount (bits 10-6)
    - (do I want to use this?)
  - 6 bits to specify operation (bits 5-0)
- I-mode (immediate) - operations using a constant value.
  - 6 bits for opcode (bits 31-26)
  - 5 bits to specify destination register (bits 25-21)
  - 5 bits to specify operand register (bits 20-16)
  - 16 bits to specify immediate value (bits 15-0)
- J-mode (jump) - operations that jump to an address.
  - 6 bits for opcode (bits 31-26)
  - 26 bits for destination address (bits 25-0)

## Outstanding Design Questions

- Do I want to have the `shamt` shift amount field in R-type instructions?
- Do I want to use the MIPS style of using I-type instructions for branching, or should I put them in the J-type, using a separate condition flags register and a dedicated `cmp` register that sets those flags?
