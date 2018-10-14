// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

type opcode struct {
	value  byte
	name   string
	length int
	opfunc func(*Engine) error
}

// Instruction -- opcode that has been parsed and includes any potential data associated with it.
type Instruction struct {
	op   *opcode
	data []byte
}

// OpCode -- returns the opcode value.
func (instr *Instruction) OpCode() byte {
	return instr.op.value
}

// Data -- returns the instruce data.
func (instr *Instruction) Data() []byte {
	return instr.data
}

// bytes --
// returns any data associated with the opcode encoded as it would be in a script.
// This is used for unparsing scripts from parsed opcodes.
func (instr *Instruction) bytes() ([]byte, error) {
	var retbytes []byte
	length := instr.op.length

	if length > 0 {
		retbytes = make([]byte, 1, length)
	} else {
		retbytes = make([]byte, 1, 1+len(instr.data)-length)
	}

	retbytes[0] = instr.op.value
	if length == 1 {
		return retbytes, nil
	}
	if length < 0 {
		l := len(instr.data)
		// tempting just to hardcode to avoid the complexity here.
		switch length {
		case -1:
			retbytes = append(retbytes, byte(l))
		case -2:
			retbytes = append(retbytes, byte(l&0xff),
				byte(l>>8&0xff))
		case -4:
			retbytes = append(retbytes, byte(l&0xff),
				byte((l>>8)&0xff), byte((l>>16)&0xff),
				byte((l>>24)&0xff))
		}
	}
	retbytes = append(retbytes, instr.data...)
	return retbytes, nil
}

// Conditional execution constants.
const (
	OpCondFalse = 0
	OpCondTrue  = 1
	OpCondSkip  = 2
)

// isConditional -- returns whether or not the opcode is a conditional opcode
// which changes the conditional execution stack when executed.
func (instr *Instruction) isConditional() bool {
	switch instr.op.value {
	case OP_IF:
		return true
	case OP_NOTIF:
		return true
	case OP_ELSE:
		return true
	case OP_ENDIF:
		return true
	default:
		return false
	}
}

// opcodes -- holds details about all possible opcodes such as how many bytes
// the opcode and any associated data should take, its human-readable name, and
// the handler function.
var opcodes = map[byte]opcode{
	// Data push opcodes.
	OP_FALSE:     {OP_FALSE, "OP_0", 1, opFalse},
	OP_DATA_1:    {OP_DATA_1, "OP_DATA_1", 2, opPushData},
	OP_DATA_2:    {OP_DATA_2, "OP_DATA_2", 3, opPushData},
	OP_DATA_3:    {OP_DATA_3, "OP_DATA_3", 4, opPushData},
	OP_DATA_4:    {OP_DATA_4, "OP_DATA_4", 5, opPushData},
	OP_DATA_5:    {OP_DATA_5, "OP_DATA_5", 6, opPushData},
	OP_DATA_6:    {OP_DATA_6, "OP_DATA_6", 7, opPushData},
	OP_DATA_7:    {OP_DATA_7, "OP_DATA_7", 8, opPushData},
	OP_DATA_8:    {OP_DATA_8, "OP_DATA_8", 9, opPushData},
	OP_DATA_9:    {OP_DATA_9, "OP_DATA_9", 10, opPushData},
	OP_DATA_10:   {OP_DATA_10, "OP_DATA_10", 11, opPushData},
	OP_DATA_11:   {OP_DATA_11, "OP_DATA_11", 12, opPushData},
	OP_DATA_12:   {OP_DATA_12, "OP_DATA_12", 13, opPushData},
	OP_DATA_13:   {OP_DATA_13, "OP_DATA_13", 14, opPushData},
	OP_DATA_14:   {OP_DATA_14, "OP_DATA_14", 15, opPushData},
	OP_DATA_15:   {OP_DATA_15, "OP_DATA_15", 16, opPushData},
	OP_DATA_16:   {OP_DATA_16, "OP_DATA_16", 17, opPushData},
	OP_DATA_17:   {OP_DATA_17, "OP_DATA_17", 18, opPushData},
	OP_DATA_18:   {OP_DATA_18, "OP_DATA_18", 19, opPushData},
	OP_DATA_19:   {OP_DATA_19, "OP_DATA_19", 20, opPushData},
	OP_DATA_20:   {OP_DATA_20, "OP_DATA_20", 21, opPushData},
	OP_DATA_21:   {OP_DATA_21, "OP_DATA_21", 22, opPushData},
	OP_DATA_22:   {OP_DATA_22, "OP_DATA_22", 23, opPushData},
	OP_DATA_23:   {OP_DATA_23, "OP_DATA_23", 24, opPushData},
	OP_DATA_24:   {OP_DATA_24, "OP_DATA_24", 25, opPushData},
	OP_DATA_25:   {OP_DATA_25, "OP_DATA_25", 26, opPushData},
	OP_DATA_26:   {OP_DATA_26, "OP_DATA_26", 27, opPushData},
	OP_DATA_27:   {OP_DATA_27, "OP_DATA_27", 28, opPushData},
	OP_DATA_28:   {OP_DATA_28, "OP_DATA_28", 29, opPushData},
	OP_DATA_29:   {OP_DATA_29, "OP_DATA_29", 30, opPushData},
	OP_DATA_30:   {OP_DATA_30, "OP_DATA_30", 31, opPushData},
	OP_DATA_31:   {OP_DATA_31, "OP_DATA_31", 32, opPushData},
	OP_DATA_32:   {OP_DATA_32, "OP_DATA_32", 33, opPushData},
	OP_DATA_33:   {OP_DATA_33, "OP_DATA_33", 34, opPushData},
	OP_DATA_34:   {OP_DATA_34, "OP_DATA_34", 35, opPushData},
	OP_DATA_35:   {OP_DATA_35, "OP_DATA_35", 36, opPushData},
	OP_DATA_36:   {OP_DATA_36, "OP_DATA_36", 37, opPushData},
	OP_DATA_37:   {OP_DATA_37, "OP_DATA_37", 38, opPushData},
	OP_DATA_38:   {OP_DATA_38, "OP_DATA_38", 39, opPushData},
	OP_DATA_39:   {OP_DATA_39, "OP_DATA_39", 40, opPushData},
	OP_DATA_40:   {OP_DATA_40, "OP_DATA_40", 41, opPushData},
	OP_DATA_41:   {OP_DATA_41, "OP_DATA_41", 42, opPushData},
	OP_DATA_42:   {OP_DATA_42, "OP_DATA_42", 43, opPushData},
	OP_DATA_43:   {OP_DATA_43, "OP_DATA_43", 44, opPushData},
	OP_DATA_44:   {OP_DATA_44, "OP_DATA_44", 45, opPushData},
	OP_DATA_45:   {OP_DATA_45, "OP_DATA_45", 46, opPushData},
	OP_DATA_46:   {OP_DATA_46, "OP_DATA_46", 47, opPushData},
	OP_DATA_47:   {OP_DATA_47, "OP_DATA_47", 48, opPushData},
	OP_DATA_48:   {OP_DATA_48, "OP_DATA_48", 49, opPushData},
	OP_DATA_49:   {OP_DATA_49, "OP_DATA_49", 50, opPushData},
	OP_DATA_50:   {OP_DATA_50, "OP_DATA_50", 51, opPushData},
	OP_DATA_51:   {OP_DATA_51, "OP_DATA_51", 52, opPushData},
	OP_DATA_52:   {OP_DATA_52, "OP_DATA_52", 53, opPushData},
	OP_DATA_53:   {OP_DATA_53, "OP_DATA_53", 54, opPushData},
	OP_DATA_54:   {OP_DATA_54, "OP_DATA_54", 55, opPushData},
	OP_DATA_55:   {OP_DATA_55, "OP_DATA_55", 56, opPushData},
	OP_DATA_56:   {OP_DATA_56, "OP_DATA_56", 57, opPushData},
	OP_DATA_57:   {OP_DATA_57, "OP_DATA_57", 58, opPushData},
	OP_DATA_58:   {OP_DATA_58, "OP_DATA_58", 59, opPushData},
	OP_DATA_59:   {OP_DATA_59, "OP_DATA_59", 60, opPushData},
	OP_DATA_60:   {OP_DATA_60, "OP_DATA_60", 61, opPushData},
	OP_DATA_61:   {OP_DATA_61, "OP_DATA_61", 62, opPushData},
	OP_DATA_62:   {OP_DATA_62, "OP_DATA_62", 63, opPushData},
	OP_DATA_63:   {OP_DATA_63, "OP_DATA_63", 64, opPushData},
	OP_DATA_64:   {OP_DATA_64, "OP_DATA_64", 65, opPushData},
	OP_DATA_65:   {OP_DATA_65, "OP_DATA_65", 66, opPushData},
	OP_DATA_66:   {OP_DATA_66, "OP_DATA_66", 67, opPushData},
	OP_DATA_67:   {OP_DATA_67, "OP_DATA_67", 68, opPushData},
	OP_DATA_68:   {OP_DATA_68, "OP_DATA_68", 69, opPushData},
	OP_DATA_69:   {OP_DATA_69, "OP_DATA_69", 70, opPushData},
	OP_DATA_70:   {OP_DATA_70, "OP_DATA_70", 71, opPushData},
	OP_DATA_71:   {OP_DATA_71, "OP_DATA_71", 72, opPushData},
	OP_DATA_72:   {OP_DATA_72, "OP_DATA_72", 73, opPushData},
	OP_DATA_73:   {OP_DATA_73, "OP_DATA_73", 74, opPushData},
	OP_DATA_74:   {OP_DATA_74, "OP_DATA_74", 75, opPushData},
	OP_DATA_75:   {OP_DATA_75, "OP_DATA_75", 76, opPushData},
	OP_PUSHDATA1: {OP_PUSHDATA1, "OP_PUSHDATA1", -1, opPushData},
	OP_PUSHDATA2: {OP_PUSHDATA2, "OP_PUSHDATA2", -2, opPushData},
	OP_PUSHDATA4: {OP_PUSHDATA4, "OP_PUSHDATA4", -4, opPushData},
	OP_TRUE:      {OP_TRUE, "OP_1", 1, opN},
	OP_2:         {OP_2, "OP_2", 1, opN},
	OP_3:         {OP_3, "OP_3", 1, opN},
	OP_4:         {OP_4, "OP_4", 1, opN},
	OP_5:         {OP_5, "OP_5", 1, opN},
	OP_6:         {OP_6, "OP_6", 1, opN},
	OP_7:         {OP_7, "OP_7", 1, opN},
	OP_8:         {OP_8, "OP_8", 1, opN},
	OP_9:         {OP_9, "OP_9", 1, opN},
	OP_10:        {OP_10, "OP_10", 1, opN},
	OP_11:        {OP_11, "OP_11", 1, opN},
	OP_12:        {OP_12, "OP_12", 1, opN},
	OP_13:        {OP_13, "OP_13", 1, opN},
	OP_14:        {OP_14, "OP_14", 1, opN},
	OP_15:        {OP_15, "OP_15", 1, opN},
	OP_16:        {OP_16, "OP_16", 1, opN},

	// Control opcodes.
	OP_IF:     {OP_IF, "OP_IF", 1, opIf},
	OP_NOTIF:  {OP_NOTIF, "OP_NOTIF", 1, opNotIf},
	OP_ELSE:   {OP_ELSE, "OP_ELSE", 1, opElse},
	OP_ENDIF:  {OP_ENDIF, "OP_ENDIF", 1, opEndIf},
	OP_RETURN: {OP_RETURN, "OP_RETURN", 1, opReturn},

	// Stack opcodes.
	OP_2DUP: {OP_2DUP, "OP_2DUP", 1, op2Dup},
	OP_3DUP: {OP_3DUP, "OP_3DUP", 1, op3Dup},
	OP_DUP:  {OP_DUP, "OP_DUP", 1, opDup},

	// Logic opcodes.
	OP_EQUAL:       {OP_EQUAL, "OP_EQUAL", 1, opEqual},
	OP_EQUALVERIFY: {OP_EQUALVERIFY, "OP_EQUALVERIFY", 1, opEqualVerify},

	// Numeric opcodes.
	OP_ADD: {OP_ADD, "OP_ADD", 1, opAdd},

	// Crypto opcodes.
	OP_HASH160:             {OP_HASH160, "OP_HASH160", 1, opHash160},
	OP_CHECKSIG:            {OP_CHECKSIG, "OP_CHECKSIG", 1, opCheckSig},
	OP_CHECKSIGVERIFY:      {OP_CHECKSIGVERIFY, "OP_CHECKSIGVERIFY", 1, opCheckSigVerify},
	OP_CHECKMULTISIG:       {OP_CHECKMULTISIG, "OP_CHECKMULTISIG", 1, opCheckMultiSig},
	OP_CHECKMULTISIGVERIFY: {OP_CHECKMULTISIGVERIFY, "OP_CHECKMULTISIGVERIFY", 1, opCheckMultiSigVerify},
}

var opcodesByName = make(map[string]byte)

func init() {
	for val, opcode := range opcodes {
		opcodesByName[opcode.name] = val
	}
	opcodesByName["OP_FALSE"] = OP_FALSE
	opcodesByName["OP_TRUE"] = OP_TRUE
	opcodesByName["OP_NOP2"] = OP_CHECKLOCKTIMEVERIFY
	opcodesByName["OP_NOP3"] = OP_CHECKSEQUENCEVERIFY
}
