
package chip8

type Opcode byte

const (
	SYS Opcode = iota
	CLS
	RET
	JP
	CALL
	SE  // Skip if equal to byte
	SNE // Skip in not equal to byte
	SRE // Skip If registers equal
	LD  // Load to register
	OR	// Bitwise OR
	AND
	XOR
	SHR
	SUBN
	SHL
	SRNE	// Skip uf registers not equal
	LDI 	// Load value to Instruction Pointer register
	JP0		// Jump to value + V0
	RND
	DRW
	SKP
	SKNP
	LDDT 	// Load delay timer to register
	LGKP	// Load keypress to register
	LDTDT	// Load register to delay time
	LDST	// Load register to sound timer
	ADDI	// Add register to Instruction Pointer
	LDF		// Set Instruction Pointer to location of sprite
	LDB		// Load BCD to I, I+1, I+2
	LDR		// Load registers to memory starting at I
	LDV		// Read memory into registers, starting at I

)