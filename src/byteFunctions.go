package main

import (
	"fmt"
	"math"
	"unsafe"
)

/*---FILE DESCRIPTION---
IKEMEN specific file.
The ikemen engine allows characters to have variables which serve as constraints, meaning certain moves can only be executed when the variables fulfill some arbitrary constraint.
To check against these constraints we use a modified function from the IKEMEN engine and only select cases where moves are executed, if the constraints for that move are fulfilled.
---FILE DESCRIPTION---*/

//goes through a byte array, and checks if the constraint function is fulfilled
func constraintCheck(constraints []byte, curGamestate CBRRawFrames_CharData, constraintNr int32) bool {
	if constraints == nil || len(constraints) <= 0 {
		return true
	}

	var be []OpCode
	for j := range constraints {
		be = append(be, OpCode(constraints[j]))
	}
	genericInt := curGamestate.GenericIntVars
	genericFloat := curGamestate.GenericFloatVars

	bcs := BytecodeStack{}
	for i := 1; i <= len(be); i++ {
		switch be[i-1] {
		case OC_jsf8:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			if bcs.Top().IsSF() {
				if be[i] == 0 {
					i = len(be)
				} else {
					i += int(uint8(be[i])) + 1
				}
			} else {
				i++
			}
		case OC_jz8, OC_jnz8:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			if bcs.Top().ToB() == (be[i-1] == OC_jz8) {
				i++
				break
			}
			fallthrough
		case OC_jmp8:
			if be[i] == 0 {
				i = len(be)
			} else {
				i += int(uint8(be[i])) + 1
			}
		case OC_jz, OC_jnz:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			if bcs.Top().ToB() == (be[i-1] == OC_jz) {
				i += 4
				break
			}
			fallthrough
		case OC_jmp:
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_parent:

		case OC_root:

		case OC_helper:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			helpBuff := false
			for _, val := range curGamestate.HelperData {
				if val.CompData.HelperID == bcs.Top().ToI() {
					genericInt = val.GenericIntVars
					genericFloat = val.GenericFloatVars
					i += 4
					helpBuff = true
				}
			}
			bcs.Pop()
			if helpBuff == true {
				continue
			}
			bcs.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_target:

		case OC_partner:

		case OC_enemy:

		case OC_enemynear:

		case OC_playerid:

		case OC_p2:

		case OC_rdreset:
			// NOP
		case OC_run:

		case OC_nordrun:

		case OC_int8:
			bcs.PushI(int32(int8(be[i])))
			i++
		case OC_int:
			bcs.PushI(*(*int32)(unsafe.Pointer(&be[i])))
			i += 4
		case OC_float:
			bcs.PushF(*(*float32)(unsafe.Pointer(&be[i])))
			i += 4
		case OC_neg:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.neg(bcs.Top())
		case OC_not:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.not(bcs.Top())
		case OC_blnot:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.blnot(bcs.Top())
		case OC_pow:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.pow(bcs.Top(), v2)
		case OC_mul:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.mul(bcs.Top(), v2)
		case OC_div:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.div(bcs.Top(), v2)
		case OC_mod:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.mod(bcs.Top(), v2)
		case OC_add:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.add(bcs.Top(), v2)
		case OC_sub:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.sub(bcs.Top(), v2)
		case OC_gt:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.gt(bcs.Top(), v2)
		case OC_ge:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.ge(bcs.Top(), v2)
		case OC_lt:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.lt(bcs.Top(), v2)
		case OC_le:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.le(bcs.Top(), v2)
		case OC_eq:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.eq(bcs.Top(), v2)
		case OC_ne:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.ne(bcs.Top(), v2)
		case OC_and:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.and(bcs.Top(), v2)
		case OC_xor:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.xor(bcs.Top(), v2)
		case OC_or:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.or(bcs.Top(), v2)
		case OC_bland:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.bland(bcs.Top(), v2)
		case OC_blxor:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.blxor(bcs.Top(), v2)
		case OC_blor:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.blor(bcs.Top(), v2)
		case OC_abs:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.abs(bcs.Top())
		case OC_exp:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.exp(bcs.Top())
		case OC_ln:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.ln(bcs.Top())
		case OC_log:
			if len(bcs) < 2 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v2 := bcs.Pop()
			bcs.log(bcs.Top(), v2)
		case OC_cos:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.cos(bcs.Top())
		case OC_sin:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.sin(bcs.Top())
		case OC_tan:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.tan(bcs.Top())
		case OC_acos:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.acos(bcs.Top())
		case OC_asin:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.asin(bcs.Top())
		case OC_atan:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.atan(bcs.Top())
		case OC_floor:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.floor(bcs.Top())
		case OC_ceil:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.ceil(bcs.Top())
		case OC_ifelse:
			if len(bcs) < 3 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			v3 := bcs.Pop()
			v2 := bcs.Pop()

			if bcs.Top().ToB() {
				*bcs.Top() = v2
			} else {
				*bcs.Top() = v3
			}
		case OC_pop:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.Pop()
		case OC_dup:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.Dup()
		case OC_swap:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			bcs.Swap()
		case OC_ailevel:

		case OC_alive:

		case OC_anim:

		case OC_animelemno:

		case OC_animelemtime:

		case OC_animexist:

		case OC_animtime:

		case OC_backedge:

		case OC_backedgebodydist:

		case OC_backedgedist:

		case OC_bottomedge:

		case OC_camerapos_x:

		case OC_camerapos_y:

		case OC_camerazoom:

		case OC_canrecover:

		case OC_command:

		case OC_ctrl:

		case OC_facing:

		case OC_frontedge:

		case OC_frontedgebodydist:

		case OC_frontedgedist:

		case OC_gameheight:

		case OC_gametime:

		case OC_gamewidth:

		case OC_hitcount:

		case OC_hitdefattr:

		case OC_hitfall:

		case OC_hitover:

		case OC_hitpausetime:

		case OC_hitshakeover:

		case OC_hitvel_x:

		case OC_hitvel_y:

		case OC_id:

		case OC_inguarddist:

		case OC_ishelper:

		case OC_leftedge:

		case OC_life:

		case OC_lifemax:

		case OC_movecontact:

		case OC_moveguarded:

		case OC_movehit:

		case OC_movereversed:

		case OC_movetype:

		case OC_numenemy:

		case OC_numexplod:

		case OC_numhelper:

		case OC_numpartner:

		case OC_numproj:

		case OC_numprojid:

		case OC_numtarget:

		case OC_palno:

		case OC_pos_x:

		case OC_pos_y:

		case OC_power:
			bcs.PushI(int32(curGamestate.MeterMax * curGamestate.MeterPercentage))
		case OC_powermax:

		case OC_playeridexist:

		case OC_prevstateno:

		case OC_projcanceltime:

		case OC_projcontacttime:

		case OC_projguardedtime:

		case OC_projhittime:

		case OC_random:
			bcs.PushI(Rand(0, 999))
		case OC_rightedge:

		case OC_roundsexisted:

		case OC_roundstate:

		case OC_screenheight:

		case OC_screenpos_x:

		case OC_screenpos_y:

		case OC_screenwidth:

		case OC_selfanimexist:

		case OC_stateno:

		case OC_statetype:

		case OC_teammode:

		case OC_teamside:

		case OC_time:

		case OC_topedge:

		case OC_uniqhitcount:

		case OC_vel_x:

		case OC_vel_y:

		case OC_st_:

		case OC_const_:

		case OC_ex_:

		case OC_var:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			*bcs.Top() = BytecodeValue{v: float64(genericInt[bcs.Top().ToI()]), t: VT_Int}
		case OC_sysvar:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			*bcs.Top() = BytecodeValue{v: float64(genericInt[bcs.Top().ToI()+int32(NumFvar)]), t: VT_Int}
		case OC_fvar:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			*bcs.Top() = BytecodeValue{v: float64(genericFloat[bcs.Top().ToI()]), t: VT_Float}
		case OC_sysfvar:
			if len(bcs) < 1 {
				print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
				return false
			}
			*bcs.Top() = BytecodeValue{v: float64(genericFloat[bcs.Top().ToI()+int32(NumFvar)]), t: VT_Float}
		case OC_localvar:

		default:
			vi := be[i-1]
			if vi < OC_sysvar0+NumSysVar {
				bcs.PushI(genericInt[vi-OC_var0])
			} else {
				bcs.PushF(genericFloat[vi-OC_fvar0])
			}
		}
		genericInt = curGamestate.GenericIntVars
		genericFloat = curGamestate.GenericFloatVars
	}
	if len(bcs) < 1 {
		print("asd")
	}
	if len(bcs) < 1 {
		print("EmptyStackError: " + fmt.Sprintf("%v", constraintNr))
		return false
	}
	return bcs.Pop().ToB()
}

type StateType int32

const (
	ST_S StateType = 1 << iota
	ST_C
	ST_A
	ST_L
	ST_N
	ST_U
	ST_MASK = 1<<iota - 1
	ST_D    = ST_L
	ST_F    = ST_N
	ST_P    = ST_U
	ST_SCA  = ST_S | ST_C | ST_A
)

type AttackType int32

const (
	AT_NA AttackType = 1 << (iota + 6)
	AT_NT
	AT_NP
	AT_SA
	AT_ST
	AT_SP
	AT_HA
	AT_HT
	AT_HP
	AT_AA  = AT_NA | AT_SA | AT_HA
	AT_AT  = AT_NT | AT_ST | AT_HT
	AT_AP  = AT_NP | AT_SP | AT_HP
	AT_ALL = AT_AA | AT_AT | AT_AP
	AT_AN  = AT_NA | AT_NT | AT_NP
	AT_AS  = AT_SA | AT_ST | AT_SP
	AT_AH  = AT_HA | AT_HT | AT_HP
)

type MoveType int32

const (
	MT_I MoveType = 1 << (iota + 15)
	MT_H
	MT_A
	MT_U
	MT_MNS = MT_I
	MT_PLS = MT_H
)

type ValueType int

const (
	VT_None ValueType = iota
	VT_Float
	VT_Int
	VT_Bool
	VT_SFalse
)

type OpCode byte

const (
	OC_var OpCode = iota + 110
	OC_sysvar
	OC_fvar
	OC_sysfvar
	OC_localvar
	OC_int8
	OC_int
	OC_float
	OC_pop
	OC_dup
	OC_swap
	OC_run
	OC_nordrun
	OC_jsf8
	OC_jmp8
	OC_jz8
	OC_jnz8
	OC_jmp
	OC_jz
	OC_jnz
	OC_eq
	OC_ne
	OC_gt
	OC_le
	OC_lt
	OC_ge
	OC_neg
	OC_blnot
	OC_bland
	OC_blxor
	OC_blor
	OC_not
	OC_and
	OC_xor
	OC_or
	OC_add
	OC_sub
	OC_mul
	OC_div
	OC_mod
	OC_pow
	OC_abs
	OC_exp
	OC_ln
	OC_log
	OC_cos
	OC_sin
	OC_tan
	OC_acos
	OC_asin
	OC_atan
	OC_floor
	OC_ceil
	OC_ifelse
	OC_time
	OC_animtime
	OC_animelemtime
	OC_animelemno
	OC_statetype
	OC_movetype
	OC_ctrl
	OC_command
	OC_random
	OC_pos_x
	OC_pos_y
	OC_vel_x
	OC_vel_y
	OC_screenpos_x
	OC_screenpos_y
	OC_facing
	OC_anim
	OC_animexist
	OC_selfanimexist
	OC_alive
	OC_life
	OC_lifemax
	OC_power
	OC_powermax
	OC_canrecover
	OC_roundstate
	OC_ishelper
	OC_numhelper
	OC_numexplod
	OC_numprojid
	OC_numproj
	OC_teammode
	OC_teamside
	OC_hitdefattr
	OC_inguarddist
	OC_movecontact
	OC_movehit
	OC_moveguarded
	OC_movereversed
	OC_projcontacttime
	OC_projhittime
	OC_projguardedtime
	OC_projcanceltime
	OC_backedge
	OC_backedgedist
	OC_backedgebodydist
	OC_frontedge
	OC_frontedgedist
	OC_frontedgebodydist
	OC_leftedge
	OC_rightedge
	OC_topedge
	OC_bottomedge
	OC_camerapos_x
	OC_camerapos_y
	OC_camerazoom
	OC_gamewidth
	OC_gameheight
	OC_screenwidth
	OC_screenheight
	OC_stateno
	OC_prevstateno
	OC_id
	OC_playeridexist
	OC_gametime
	OC_numtarget
	OC_numenemy
	OC_numpartner
	OC_ailevel
	OC_palno
	OC_hitcount
	OC_uniqhitcount
	OC_hitpausetime
	OC_hitover
	OC_hitshakeover
	OC_hitfall
	OC_hitvel_x
	OC_hitvel_y
	OC_roundsexisted
	OC_parent
	OC_root
	OC_helper
	OC_target
	OC_partner
	OC_enemy
	OC_enemynear
	OC_playerid
	OC_p2
	OC_rdreset
	OC_const_
	OC_st_
	OC_ex_
	OC_var0     = 0
	OC_sysvar0  = 60
	OC_fvar0    = 65
	OC_sysfvar0 = 105
)
const (
	NumSysVar = OC_fvar0 - OC_sysvar0
	NumFvar   = OC_sysfvar0 - OC_fvar0
)

type BytecodeValue struct {
	t ValueType
	v float64
}

func (bv BytecodeValue) IsNone() bool { return bv.t == VT_None }
func (bv BytecodeValue) IsSF() bool   { return bv.t == VT_SFalse }
func (bv BytecodeValue) ToF() float32 {
	if bv.IsSF() {
		return 0
	}
	return float32(bv.v)
}
func (bv BytecodeValue) ToI() int32 {
	if bv.IsSF() {
		return 0
	}
	return int32(bv.v)
}
func (bv BytecodeValue) ToB() bool {
	if bv.IsSF() || bv.v == 0 {
		return false
	}
	return true
}
func (bv *BytecodeValue) SetF(f float32) {
	if math.IsNaN(float64(f)) {
		*bv = BytecodeSF()
	} else {
		*bv = BytecodeValue{VT_Float, float64(f)}
	}
}
func (bv *BytecodeValue) SetI(i int32) {
	*bv = BytecodeValue{VT_Int, float64(i)}
}
func (bv *BytecodeValue) SetB(b bool) {
	bv.t = VT_Bool
	if b {
		bv.v = 1
	} else {
		bv.v = 0
	}
}

func bvNone() BytecodeValue {
	return BytecodeValue{VT_None, 0}
}
func BytecodeSF() BytecodeValue {
	return BytecodeValue{VT_SFalse, math.NaN()}
}
func BytecodeFloat(f float32) BytecodeValue {
	return BytecodeValue{VT_Float, float64(f)}
}
func BytecodeInt(i int32) BytecodeValue {
	return BytecodeValue{VT_Int, float64(i)}
}
func BytecodeBool(b bool) BytecodeValue {
	return BytecodeValue{VT_Bool, float64(Btoi(b))}
}

type BytecodeStack []BytecodeValue

func (bs *BytecodeStack) Clear()                { *bs = (*bs)[:0] }
func (bs *BytecodeStack) Push(bv BytecodeValue) { *bs = append(*bs, bv) }
func (bs *BytecodeStack) PushI(i int32)         { bs.Push(BytecodeInt(i)) }
func (bs *BytecodeStack) PushF(f float32)       { bs.Push(BytecodeFloat(f)) }
func (bs *BytecodeStack) PushB(b bool)          { bs.Push(BytecodeBool(b)) }
func (bs BytecodeStack) Top() *BytecodeValue {
	return &bs[len(bs)-1]
}

func (bs *BytecodeStack) Pop() (bv BytecodeValue) {
	bv, *bs = *bs.Top(), (*bs)[:len(*bs)-1]
	return
}
func (bs *BytecodeStack) Dup() {
	bs.Push(*bs.Top())
}
func (bs *BytecodeStack) Swap() {
	*bs.Top(), (*bs)[len(*bs)-2] = (*bs)[len(*bs)-2], *bs.Top()
}
func (bs *BytecodeStack) Alloc(size int) []BytecodeValue {
	if len(*bs)+size > cap(*bs) {
		tmp := *bs
		*bs = make(BytecodeStack, len(*bs)+size)
		copy(*bs, tmp)
	} else {
		*bs = (*bs)[:len(*bs)+size]
		for i := len(*bs) - size; i < len(*bs); i++ {
			(*bs)[i] = bvNone()
		}
	}
	return (*bs)[len(*bs)-size:]
}

type BytecodeExp []OpCode

func (be *BytecodeExp) append(op ...OpCode) {
	*be = append(*be, op...)
}
func (be *BytecodeExp) appendValue(bv BytecodeValue) (ok bool) {
	switch bv.t {
	case VT_Float:
		be.append(OC_float)
		f := float32(bv.v)
		be.append((*(*[4]OpCode)(unsafe.Pointer(&f)))[:]...)
	case VT_Int:
		if bv.v >= -128 && bv.v <= 127 {
			be.append(OC_int8, OpCode(bv.v))
		} else {
			be.append(OC_int)
			i := int32(bv.v)
			be.append((*(*[4]OpCode)(unsafe.Pointer(&i)))[:]...)
		}
	case VT_Bool:
		if bv.v != 0 {
			be.append(OC_int8, 1)
		} else {
			be.append(OC_int8, 0)
		}
	case VT_SFalse:
		be.append(OC_int8, 0)
	default:
		return false
	}
	return true
}
func (be *BytecodeExp) appendI32Op(op OpCode, addr int32) {
	be.append(op)
	be.append((*(*[4]OpCode)(unsafe.Pointer(&addr)))[:]...)
}
func (bs *BytecodeStack) neg(v *BytecodeValue) {
	if v.t == VT_Float {
		v.v *= -1
	} else {
		v.SetI(-v.ToI())
	}
}
func (bs *BytecodeStack) not(v *BytecodeValue) {
	v.SetI(^v.ToI())
}
func (bs *BytecodeStack) blnot(v *BytecodeValue) {
	v.SetB(!v.ToB())
}
func (bs *BytecodeStack) pow(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(Pow(v1.ToF(), v2.ToF()))
	} else if v2.ToF() < 0 {
		v1.SetF(Pow(v1.ToF(), v2.ToF()))
	} else {
		i1, i2, hb := v1.ToI(), v2.ToI(), int32(-1)
		for uint32(i2)>>uint(hb+1) != 0 {
			hb++
		}
		var i, bit, tmp int32 = 1, 0, i1
		for ; bit <= hb; bit++ {
			var shift uint
			if bit == hb {
				shift = uint(bit)
			} else {
				shift = uint((hb - 1) - bit)
			}
			if i2&(1<<shift) != 0 {
				i *= tmp
			}
			tmp *= tmp
		}
		v1.SetI(i)
	}
}
func (bs *BytecodeStack) mul(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() * v2.ToF())
	} else {
		v1.SetI(v1.ToI() * v2.ToI())
	}
}
func (bs *BytecodeStack) div(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() / v2.ToF())
	} else if v2.ToI() == 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetI(v1.ToI() / v2.ToI())
	}
}
func (bs *BytecodeStack) mod(v1 *BytecodeValue, v2 BytecodeValue) {
	if v2.ToI() == 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetI(v1.ToI() % v2.ToI())
	}
}
func (bs *BytecodeStack) add(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() + v2.ToF())
	} else {
		v1.SetI(v1.ToI() + v2.ToI())
	}
}
func (bs *BytecodeStack) sub(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() - v2.ToF())
	} else {
		v1.SetI(v1.ToI() - v2.ToI())
	}
}
func (bs *BytecodeStack) gt(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() > v2.ToF())
	} else {
		v1.SetB(v1.ToI() > v2.ToI())
	}
}
func (bs *BytecodeStack) ge(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() >= v2.ToF())
	} else {
		v1.SetB(v1.ToI() >= v2.ToI())
	}
}
func (bs *BytecodeStack) lt(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() < v2.ToF())
	} else {
		v1.SetB(v1.ToI() < v2.ToI())
	}
}
func (bs *BytecodeStack) le(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() <= v2.ToF())
	} else {
		v1.SetB(v1.ToI() <= v2.ToI())
	}
}
func (bs *BytecodeStack) eq(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() == v2.ToF())
	} else {
		v1.SetB(v1.ToI() == v2.ToI())
	}
}
func (bs *BytecodeStack) ne(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() != v2.ToF())
	} else {
		v1.SetB(v1.ToI() != v2.ToI())
	}
}
func (bs *BytecodeStack) and(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() & v2.ToI())
}
func (bs *BytecodeStack) xor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() ^ v2.ToI())
}
func (bs *BytecodeStack) or(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() | v2.ToI())
}
func (bs *BytecodeStack) bland(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() && v2.ToB())
}
func (bs *BytecodeStack) blxor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() != v2.ToB())
}
func (bs *BytecodeStack) blor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() || v2.ToB())
}
func (bs *BytecodeStack) abs(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		v1.v = math.Abs(v1.v)
	} else {
		v1.SetI(Abs(v1.ToI()))
	}
}
func (bs *BytecodeStack) exp(v1 *BytecodeValue) {
	v1.SetF(float32(math.Exp(v1.v)))
}
func (bs *BytecodeStack) ln(v1 *BytecodeValue) {
	if v1.v <= 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetF(float32(math.Log(v1.v)))
	}
}
func (bs *BytecodeStack) log(v1 *BytecodeValue, v2 BytecodeValue) {
	if v1.v <= 0 || v2.v <= 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetF(float32(math.Log(v2.v) / math.Log(v1.v)))
	}
}
func (bs *BytecodeStack) cos(v1 *BytecodeValue) {
	v1.SetF(float32(math.Cos(v1.v)))
}
func (bs *BytecodeStack) sin(v1 *BytecodeValue) {
	v1.SetF(float32(math.Sin(v1.v)))
}
func (bs *BytecodeStack) tan(v1 *BytecodeValue) {
	v1.SetF(float32(math.Tan(v1.v)))
}
func (bs *BytecodeStack) acos(v1 *BytecodeValue) {
	v1.SetF(float32(math.Acos(v1.v)))
}
func (bs *BytecodeStack) asin(v1 *BytecodeValue) {
	v1.SetF(float32(math.Asin(v1.v)))
}
func (bs *BytecodeStack) atan(v1 *BytecodeValue) {
	v1.SetF(float32(math.Atan(v1.v)))
}
func (bs *BytecodeStack) floor(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		f := math.Floor(v1.v)
		if math.IsNaN(f) {
			*v1 = BytecodeSF()
		} else {
			v1.SetI(int32(f))
		}
	}
}
func (bs *BytecodeStack) ceil(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		f := math.Ceil(v1.v)
		if math.IsNaN(f) {
			*v1 = BytecodeSF()
		} else {
			v1.SetI(int32(f))
		}
	}
}
func (bs *BytecodeStack) max(v1 *BytecodeValue, v2 BytecodeValue) {
	if v1.v >= v2.v {
		v1.SetF(float32(v1.v))
	} else {
		v1.SetF(float32(v2.v))
	}
}
func (bs *BytecodeStack) min(v1 *BytecodeValue, v2 BytecodeValue) {
	if v1.v <= v2.v {
		v1.SetF(float32(v1.v))
	} else {
		v1.SetF(float32(v2.v))
	}
}
func (bs *BytecodeStack) random(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(RandI(int32(v1.v), int32(v2.v)))
}
func (bs *BytecodeStack) round(v1 *BytecodeValue, v2 BytecodeValue) {
	shift := math.Pow(10, v2.v)
	v1.SetF(float32(math.Floor((v1.v*shift)+0.5) / shift))
}

//common

func Min(arg ...int32) (min int32) {
	if len(arg) > 0 {
		min = arg[0]
		for i := 1; i < len(arg); i++ {
			if arg[i] < min {
				min = arg[i]
			}
		}
	}
	return
}

func Abs(i int32) int32 {
	if i < 0 {
		return -i
	}
	return i
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

var randInt int32

func Random() int32 {
	w := randInt / 127773
	randInt = (randInt-w*127773)*16807 - w*2836
	if randInt <= 0 {
		randInt += IMax - Btoi(randInt == 0)
	}
	return randInt
}

func Rand(min, max int32) int32 { return min + Random()/(IMax/(max-min+1)+1) }

func RandI(x, y int32) int32 {
	if y < x {
		if uint32(x-y) > uint32(IMax) {
			return int32(int64(y) + int64(Random())*(int64(x)-int64(y))/int64(IMax))
		}
		return Rand(y, x)
	}
	if uint32(y-x) > uint32(IMax) {
		return int32(int64(x) + int64(Random())*(int64(y)-int64(x))/int64(IMax))
	}
	return Rand(x, y)
}

const (
	IMax = int32(^uint32(0) >> 1)
	IErr = ^IMax
)

func Btoi(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
