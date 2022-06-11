package main

import (
	"fmt"
	"time"
)

/*---FILE DESCRIPTION---
The parameter file contains structs with parameters that changes how the CBR AI interprets or works with the gamestate and its cases.
This file alo contains the aiData struct which is where general CBRAI data is stored.
---FILE DESCRIPTION---*/

type debugStorageStringFloat32 struct {
	s string
	v float32
}
type debugStorageStringX2Float64 struct {
	s  string
	s2 string
	v  float64
}

var debugUtil = debugComparisons{
	debugActive: false,
	maxCases:    20,
	comparisons: []comparisonInstance{},
	debugText:   "",
}

type debugComparisons struct {
	debugActive bool
	maxCases    int
	comparisons []comparisonInstance
	inWorkCase  *caseDebugData
	debugText   string
}

type comparisonInstance struct {
	curGamestateFrame       *CBRRawFrames_Frame
	curGamestateCaseOutputs []*debugStorageStringX2Float64
	chosenCaseIndex         int
	chosenCaseReplayIndex   int
	chosenCaseCompValue     float32
	topCases                []*caseDebugData
	nextCaseIndex           int
	nextCaseReplayIndex     int
	replayFramecount        int64
	caseSelectionReason     string
	nextCase                *caseDebugData
}

type caseDebugData struct {
	replayIndex      int
	caseIndex        int
	caseData         *CBRData_Case
	compValue        float32
	compOutputs      []*debugStorageStringFloat32
	caseOutputs      []*debugStorageStringX2Float64
	debugDescriptors []*DebugDescriptor
}

type NavMapString_Float32 struct {
	m    map[string]float32
	keys []string
}

func (n *NavMapString_Float32) Add(k string, v float32) {
	_, ok := n.m[k]
	n.m[k] = v
	if !ok {
		n.keys = append(n.keys, k)
	}
}

type NavMapString_Float64 struct {
	m    map[string]float64
	keys []string
}

func (n *NavMapString_Float64) Add(k *string, v *float64) {
	_, ok := n.m[*k]
	n.m[*k] = *v
	if !ok {
		n.keys = append(n.keys, *k)
	}
}

type NavMapString_NavMapString_Float64 struct {
	m    map[string]NavMapString_Float64
	keys []string
}

func (n *NavMapString_NavMapString_Float64) Add(k string, v NavMapString_Float64) {
	_, ok := n.m[k]
	n.m[k] = v
	if !ok {
		n.keys = append(n.keys, k)
	}

}

func (x *debugComparisons) addReplayIndex(index int) {
	if !x.debugActive {
		return
	}
	x.inWorkCase.replayIndex = index
}
func (x *debugComparisons) addCaseIndex(index int) {
	if !x.debugActive {
		return
	}
	x.inWorkCase.caseIndex = index
}
func (x *debugComparisons) addCaseData(caseData *CBRData_Case) {
	if !x.debugActive {
		return
	}
	x.inWorkCase.caseData = caseData
}
func (x *debugComparisons) addCompValue(compValue float32) {
	if !x.debugActive {
		return
	}
	x.inWorkCase.compValue = compValue
}
func (x *debugComparisons) addCompOutputs() {
	if !x.debugActive {
		return
	}
	x.inWorkCase.compOutputs = []*debugStorageStringFloat32{} //&NavMapString_Float32{map[string]float32{}, []string{}}
}
func (x *debugComparisons) addCompOutputValue(key string, value float32) {
	if !x.debugActive {
		return
	}
	x.inWorkCase.compOutputs = append(x.inWorkCase.compOutputs, &debugStorageStringFloat32{s: key, v: value}) //.Add(key, value)
}
func (x *caseDebugData) addCaseOutputValue(key string, key2 string, value float64) {

	/*if x.caseOutputs.m == nil{
		x.caseOutputs.m = map[string]NavMapString_Float64{}
		x.caseOutputs.keys = []string{}
	}

	val, ok := x.caseOutputs.m[key]
	if !ok{
		val = NavMapString_Float64{map[string]float64{}, []string{}}
	}
	val.Add(&key2, &value)
	x.caseOutputs.m[key] = val

	*/
	x.caseOutputs = append(x.caseOutputs, &debugStorageStringX2Float64{s: key, s2: key2, v: value})
}
func (x *debugComparisons) addCurCaseOutputValue(key string, key2 string, value float64) {
	if !x.debugActive {
		return
	}
	/*
		if x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs.m == nil{
			x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs.m = map[string]NavMapString_Float64{}
			x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs.keys = []string{}
		}

		val, ok := x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs.m[key]
		if !ok{
			val = NavMapString_Float64{map[string]float64{}, []string{}}
		}
		val.Add(&key2, &value)

	*/
	x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs = append(x.comparisons[len(x.comparisons)-1].curGamestateCaseOutputs,
		&debugStorageStringX2Float64{s: key, s2: key2, v: value}) //.Add(key, val)
}
func (x *debugComparisons) addChosenCase(caseIndex int, replayIndex int, chosenCaseCompValue float32) {
	if !x.debugActive {
		return
	}
	x.comparisons[len(x.comparisons)-1].chosenCaseIndex = caseIndex
	x.comparisons[len(x.comparisons)-1].chosenCaseReplayIndex = replayIndex
	x.comparisons[len(x.comparisons)-1].chosenCaseCompValue = chosenCaseCompValue
}
func (x *debugComparisons) addNextCase(caseIndex int, replayIndex int) {
	if !x.debugActive {
		return
	}
	x.comparisons[len(x.comparisons)-1].nextCaseIndex = caseIndex
	x.comparisons[len(x.comparisons)-1].nextCaseReplayIndex = replayIndex
}
func (x *debugComparisons) setSelectionReason(reason string) {
	if !x.debugActive {
		return
	}
	x.comparisons[len(x.comparisons)-1].caseSelectionReason = reason

}
func (x *debugComparisons) debugCheck() bool {
	return x.debugActive == true
}

func (x *debugComparisons) addFrame(frameCount int64) {
	if !x.debugActive {
		return
	}
	//if len(x.comparisons) >= x.maxFrames{
	//	x.comparisons = x.comparisons[1:]
	//}
	x.comparisons = append(x.comparisons, comparisonInstance{replayFramecount: frameCount})
}

func (x *debugComparisons) resetFrame(frameCount int64) {
	if !x.debugActive {
		return
	}
	//if len(x.comparisons) >= x.maxFrames{
	//	x.comparisons = x.comparisons[1:]
	//}
	x.comparisons[len(x.comparisons)-1] = comparisonInstance{replayFramecount: frameCount}
}

func (x *debugComparisons) addCurGamestateFrame(curGamestate *CBRRawFrames) {
	if !x.debugActive {
		return
	}
	x.comparisons[len(x.comparisons)-1].curGamestateFrame = curGamestate.ReplayFile[len(curGamestate.ReplayFile)-1].Frame[len(curGamestate.ReplayFile[len(curGamestate.ReplayFile)-1].Frame)-1]
}
func (x *debugComparisons) addComparison() {
	if !x.debugActive {
		return
	}
	lastCompIndex := len(x.comparisons) - 1
	if x.comparisons[lastCompIndex].topCases == nil {
		x.comparisons[lastCompIndex].topCases = []*caseDebugData{}
	}
	x.comparisons[lastCompIndex].topCases = append(x.comparisons[lastCompIndex].topCases, x.inWorkCase)
}
func (x *debugComparisons) sortComparisons(comparisonFrameIndex int) (topCases []*caseDebugData, nextCase *caseDebugData) {

	sortedTopCases := []*caseDebugData{}
	nextCase = nil

	lastCompIndex := comparisonFrameIndex
	for j, tCase := range x.comparisons[lastCompIndex].topCases {
		if tCase.caseIndex == x.comparisons[lastCompIndex].nextCaseIndex && tCase.replayIndex == x.comparisons[lastCompIndex].nextCaseReplayIndex {
			nextCase = x.comparisons[lastCompIndex].topCases[j]
		}

		if sortedTopCases == nil {
			sortedTopCases = []*caseDebugData{}
		}
		index := 0
		for i, val := range sortedTopCases {
			if val.compValue <= tCase.compValue {
				index = i
				break
			}
			if i == len(sortedTopCases)-1 {
				index = i + 1
			}
		}
		if len(sortedTopCases) >= x.maxCases && index != 0 {
			sortedTopCases = sortedTopCases[1:]
			index--
		}

		if len(sortedTopCases) < x.maxCases {

			sortedTopCases = append(sortedTopCases, &caseDebugData{})
			copy(sortedTopCases[index+1:], sortedTopCases[index:])
			sortedTopCases[index] = tCase
		}
	}

	return sortedTopCases, nextCase

}

func (x *debugComparisons) resetWorkingCase() {
	if !x.debugActive {
		return
	}
	x.inWorkCase = &caseDebugData{}
}

func (x *debugComparisons) prepareDebugData(curGamestate *CBRRawFrames_ReplayFile, curGameFrame *CBRRawFrames_Frame, caseData *CBRData, comparisonFrameIndex int) {
	if !x.debugActive {
		return
	}
	//curFrame := x.comparisons[len(x.comparisons)-1]
	sortedTopCases, nextCase := x.sortComparisons(comparisonFrameIndex)

	for i := len(sortedTopCases) - 1; i >= 0; i-- {
		sortedTopCases[i].debugGetCaseInfo(sortedTopCases[i].caseData, curGamestate, curGameFrame, caseData, sortedTopCases[i].replayIndex)
		sortedTopCases[i].debugDescriptors = sortedTopCases[i].caseData.DebugDescriptors
	}
	nextCase.debugGetCaseInfo(nextCase.caseData, curGamestate, curGameFrame, caseData, nextCase.replayIndex)
	nextCase.debugDescriptors = nextCase.caseData.DebugDescriptors

	x.comparisons[comparisonFrameIndex].topCases = sortedTopCases
	x.comparisons[comparisonFrameIndex].nextCase = nextCase
}

func (x *debugComparisons) printDebugData() {
	if !x.debugActive {
		return
	}
	//curFrame := x.comparisons[len(x.comparisons)-1]
	for _, curFrame := range x.comparisons {

		nextCase := curFrame.nextCase
		sortedTopCases := curFrame.topCases
		/*
			for i, tCase := range curFrame.topCases {
				if tCase.caseIndex == curFrame.nextCaseIndex &&  tCase.replayIndex == curFrame.nextCaseReplayIndex {
					nextCase = curFrame.topCases[i]
				}
			}
		*/

		printString := fmt.Sprintf("repFrame: %v ----------------------------------------------------------", curFrame.replayFramecount)
		printString += fmt.Sprintf("\nCaseSelectionReason: %v\nChosen Case: %v - %v - Final Cost: %v\nCurFrame:", curFrame.caseSelectionReason, curFrame.chosenCaseReplayIndex, curFrame.chosenCaseIndex, curFrame.chosenCaseCompValue)
		/*
			counter := 0
			for _, id := range curFrame.curGamestateCaseOutputs.keys {
				val, _ := curFrame.curGamestateCaseOutputs.m[id]
				printString += fmt.Sprintf("\n%s:",id)

				for _, id2 := range val.keys{
					val2, _ := val.m[id2]
					printString += fmt.Sprintf("%s: %.2f, ",id2, val2)
				}
				counter++
			}
		*/

		lastKey := ""
		for i, val := range curFrame.curGamestateCaseOutputs {
			if lastKey != val.s {
				lastKey = val.s
				key := val.s
				printString += fmt.Sprintf("\n%s:", key)
				for j := i; j < len(curFrame.curGamestateCaseOutputs); j++ {
					val2 := curFrame.curGamestateCaseOutputs[j]
					if key == val2.s {
						printString += fmt.Sprintf("%s: %.2f, ", val2.s2, val2.v)
					}
				}
			}
		}

		if nextCase != nil {
			printString += fmt.Sprintf("\n\nNextCaseinReplay:\nCaseNr: %v - %v", nextCase.replayIndex, nextCase.caseIndex)
			printString += fmt.Sprintf(" - Ctrl: %v - Final Cost: %v ", nextCase.caseData.Controllable, nextCase.compValue)
			var frameStart float64 = -1
			var frameEnd float64 = -1
			for _, val := range nextCase.caseOutputs {
				key := val.s
				if key == "CaseLength" {
					if val.s2 == "FrameStart" {
						frameStart = val.v
					}
					if val.s2 == "FrameEnd" {
						frameEnd = val.v
					}
				}

			}
			printString += fmt.Sprintf(" \nFrameStart: %v - FrameEnd: %v", frameStart, frameEnd)

			lastKey = ""
			for _, val := range nextCase.compOutputs {
				if lastKey != val.s {
					lastKey = val.s
					if val.v != 0 {
						printString += fmt.Sprintf("\n%s: %.3f", val.s, val.v)

						for _, val2 := range nextCase.caseOutputs {
							if val2.s != val.s {
								continue
							} else {
								printString += fmt.Sprintf(" - %s: %.2f", val2.s2, val2.v)
							}
						}
					}
				}

			}

			printString += "\nDebugDescriptors:\n"
			for i, desc := range nextCase.caseData.DebugDescriptors {
				printString += fmt.Sprintf("Descriptor %v: %v - %v - %v\n", i, desc.Primary, desc.Secondary, desc.Tertiary)
			}
		}

		for i := len(sortedTopCases) - 1; i >= 0; i-- {
			printString += fmt.Sprintf("\n\nCaseNr: %v - %v", sortedTopCases[i].replayIndex, sortedTopCases[i].caseIndex)
			printString += fmt.Sprintf(" - Ctrl: %v - Final Cost: %v ", sortedTopCases[i].caseData.Controllable, sortedTopCases[i].compValue)
			var frameStart float64 = -1
			var frameEnd float64 = -1
			for _, val := range sortedTopCases[i].caseOutputs {
				key := val.s
				if key == "CaseLength" {
					if val.s2 == "FrameStart" {
						frameStart = val.v
					}
					if val.s2 == "FrameEnd" {
						frameEnd = val.v
					}
				}

			}
			printString += fmt.Sprintf("\nFrameStart: %v - FrameEnd: %v", frameStart, frameEnd)

			lastKey = ""
			for _, val := range sortedTopCases[i].compOutputs {
				if lastKey != val.s {
					lastKey = val.s
					if val.v != 0 {
						printString += fmt.Sprintf("\n%s: %.3f", val.s, val.v)

						for _, val2 := range sortedTopCases[i].caseOutputs {
							if val2.s != val.s {
								continue
							} else {
								printString += fmt.Sprintf(" - %s: %.2f", val2.s2, val2.v)
							}
						}
					}
				}
			}

			printString += "\nDebugDescriptors:\n"
			for j, desc := range sortedTopCases[i].caseData.DebugDescriptors {
				printString += fmt.Sprintf("Descriptor %v: %v - %v - %v\n", j, desc.Primary, desc.Secondary, desc.Tertiary)
			}

		}

		printString += "\n\n\n"

		debugUtil.debugText += printString
	}

	//screen.Clear()
	//screen.MoveTopLeft()
	//fmt.Printf("%s", printString)

}

func (debugUtil debugComparisons) debugGetFrameInfo(rawFrames *CBRRawFrames, refIds CharReferenceIDs) {
	if debugUtil.debugActive == false {
		return
	}
	if rawFrames.ReplayFile == nil || len(rawFrames.ReplayFile) <= 0 {
		return
	}
	if rawFrames.ReplayFile[len(rawFrames.ReplayFile)-1].Frame == nil || len(rawFrames.ReplayFile[len(rawFrames.ReplayFile)-1].Frame) <= 0 {
		return
	}

	workingFrame := rawFrames.ReplayFile[len(rawFrames.ReplayFile)-1].Frame[len(rawFrames.ReplayFile[len(rawFrames.ReplayFile)-1].Frame)-1]
	debugUtil.addCurCaseOutputValue("xRelativePosition", "Focus", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos))
	debugUtil.addCurCaseOutputValue("xRelativePosition", "Enemy", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.CharPos.XPos))
	debugUtil.addCurCaseOutputValue("yRelativePosition", "Focus", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.YPos))
	debugUtil.addCurCaseOutputValue("yRelativePosition", "Enemy", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.CharPos.YPos))
	debugUtil.addCurCaseOutputValue("airborneState", "air", boolToFloat64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.MStateAir))
	debugUtil.addCurCaseOutputValue("lyingDownState", "Lying", boolToFloat64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.MStateLying))
	debugUtil.addCurCaseOutputValue("hitState", "hit", boolToFloat64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.AStateHit))
	debugUtil.addCurCaseOutputValue("hitState", "hitStun", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.HitStun))
	debugUtil.addCurCaseOutputValue("blockState", "blockStun", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.Blockstun))
	debugUtil.addCurCaseOutputValue("attackState", "frame", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.CurrentMoveFrame))
	debugUtil.addCurCaseOutputValue("attackState", "atk", boolToFloat64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.AStateAttack))
	debugUtil.addCurCaseOutputValue("nearWall", "LWall", float64(workingFrame.WorldCBRComparisonData.StageData.LeftWallPos))
	debugUtil.addCurCaseOutputValue("nearWall", "RWall", float64(workingFrame.WorldCBRComparisonData.StageData.RightWallPos))
	debugUtil.addCurCaseOutputValue("moveID", "MID", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.CurrentMoveReferenceID))
	debugUtil.addCurCaseOutputValue("roundState", "", float64(workingFrame.WorldCBRComparisonData.RoundState))
	debugUtil.addCurCaseOutputValue("getHit", "hit", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.SelfHit))
	debugUtil.addCurCaseOutputValue("getHit", "guard", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.SelfGuard))
	debugUtil.addCurCaseOutputValue("didHit", "hit", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.MoveHit))
	debugUtil.addCurCaseOutputValue("didHit", "guard", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.MoveGuarded))
	debugUtil.addCurCaseOutputValue("enemyAirborneState", "eAir", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.MStateAir))
	debugUtil.addCurCaseOutputValue("enemyHitState", "eHit", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.AStateHit))
	debugUtil.addCurCaseOutputValue("enemyHitState", "hitStun", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.HitStun))
	debugUtil.addCurCaseOutputValue("enemyBlockState", "blockStun", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.Blockstun))
	debugUtil.addCurCaseOutputValue("enemyAttackState", "eAtk", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.AStateAttack))
	debugUtil.addCurCaseOutputValue("enemyAttackState", "Frame", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.CurrentMoveFrame))
	debugUtil.addCurCaseOutputValue("enemyMoveID", "eMID", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.CurrentMoveReferenceID))
	debugUtil.addCurCaseOutputValue("enemyLyingDownState", "eLying", boolToFloat64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.MStateLying))
	debugUtil.addCurCaseOutputValue("yVelocity", "", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.Velocity.YVel))
	debugUtil.addCurCaseOutputValue("xVelocity", "", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.Velocity.XVel))
	debugUtil.addCurCaseOutputValue("enemyXVelocity", "", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.Velocity.XVel))
	debugUtil.addCurCaseOutputValue("enemyYVelocity", "", float64(workingFrame.CharData[refIds.curEnemyCharNr].ComparisonData.Velocity.YVel))
	debugUtil.addCurCaseOutputValue("frameAdv", "", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.FrameAdv))
	debugUtil.addCurCaseOutputValue("frameAdvInitiator", "", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.FrameAdv))
	debugUtil.addCurCaseOutputValue("comboSimilarity", "", float64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.ComboMovesUsed))
	debugUtil.addCurCaseOutputValue("comboSimilarity", "", boolToFloat64(workingFrame.CharData[refIds.curFocusCharNr].ComparisonData.Pressure))

	for _, val := range ikemenVarImport[refIds.curFocusCharNr] {
		if val.HelperID >= 0 {
			continue
		}
		if val.Float {
			debugUtil.addCurCaseOutputValue("ikemenVar", fmt.Sprintf("FVar%v", val.VarNr), float64(workingFrame.CharData[refIds.curFocusCharNr].GenericFloatVars[val.VarNr]))
		} else {
			debugUtil.addCurCaseOutputValue("ikemenVar", fmt.Sprintf("IVar%v", val.VarNr), float64(workingFrame.CharData[refIds.curFocusCharNr].GenericIntVars[val.VarNr]))
		}
	}

	counter := 0
	for i, val := range workingFrame.CharData[refIds.curEnemyCharNr].HelperData {

		if val.CompData.HitboxOut || val.CompData.HurtboxOut || val.CompData.AStateAttack {
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("eHelper_%v_Position", i), "PosX", float64(val.CompData.PositionX))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("eHelper_%v_Position", i), "PosY", float64(val.CompData.PositionY))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("eHelper_%v_Velocity", i), "VelX", float64(val.CompData.Velocity.XVel))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("eHelper_%v_Velocity", i), "VelY", float64(val.CompData.Velocity.YVel))
			counter++
		}
	}
	debugUtil.addCurCaseOutputValue("enemyHelperAmount", "", float64(counter))
	counter = 0

	for _, val2 := range ikemenVarImport[refIds.curFocusCharNr] {
		if val2.VarNr == -1 {
			match := false
			for _, val := range workingFrame.CharData[refIds.curFocusCharNr].HelperData {
				if val2.HelperID == val.CompData.HelperID {
					match = true
				}
			}
			debugUtil.addCurCaseOutputValue("ikemenVar", fmt.Sprintf("H%v", val2.HelperID), boolToFloat64(match))
		}
	}

	for i, val := range workingFrame.CharData[refIds.curFocusCharNr].HelperData {

		for _, val2 := range ikemenVarImport[refIds.curFocusCharNr] {
			if val2.HelperID < 0 || val2.HelperID != val.CompData.HelperID {
				continue
			}

			if val2.VarNr != -1 {
				if val2.Float {
					debugUtil.addCurCaseOutputValue("ikemenVar", fmt.Sprintf("H%vFVar%v", val2.HelperID, val2.VarNr), float64(val.GenericFloatVars[val2.VarNr]))
				} else {
					debugUtil.addCurCaseOutputValue("ikemenVar", fmt.Sprintf("H%vIVar%v", val2.HelperID, val2.VarNr), float64(val.GenericIntVars[val2.VarNr]))
				}
			}

		}

		if val.CompData.HitboxOut || val.CompData.HurtboxOut || val.CompData.AStateAttack {
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("Helper_%v_Position", i), "PosX", float64(val.CompData.PositionX))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("Helper_%v_Position", i), "PosY", float64(val.CompData.PositionY))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("Helper_%v_Velocity", i), "VelX", float64(val.CompData.Velocity.XVel))
			debugUtil.addCurCaseOutputValue(fmt.Sprintf("Helper_%v_Velocity", i), "VelY", float64(val.CompData.Velocity.YVel))
			counter++
		}
	}
	debugUtil.addCurCaseOutputValue("HelperAmount", "", float64(counter))

}

func (x *caseDebugData) debugGetCaseInfo(workingCase *CBRData_Case, curGamestate *CBRRawFrames_ReplayFile, curFrame *CBRRawFrames_Frame, caseData *CBRData, replayIndex int) {
	refIds := refIdsMapping(curGamestate, curFrame, caseData, replayIndex)
	x.addCaseOutputValue("CaseLength", "FrameStart", float64(workingCase.FrameStartId))
	x.addCaseOutputValue("CaseLength", "FrameEnd", float64(workingCase.FrameEndId))
	x.addCaseOutputValue("xRelativePosition", "Focus", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos))
	x.addCaseOutputValue("xRelativePosition", "Enemy", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CharPos.XPos))
	x.addCaseOutputValue("yRelativePosition", "Focus", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.YPos))
	x.addCaseOutputValue("yRelativePosition", "Enemy", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CharPos.YPos))
	x.addCaseOutputValue("airborneState", "Air", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MStateAir))
	x.addCaseOutputValue("lyingDownState", "Lying", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MStateLying))
	x.addCaseOutputValue("hitState", "Hit", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateHit))
	x.addCaseOutputValue("hitState", "HitStun", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.HitStun))
	x.addCaseOutputValue("blockState", "Block", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateHit))
	x.addCaseOutputValue("blockState", "BlockStun", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Blockstun))
	x.addCaseOutputValue("attackState", "Atk", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateAttack))
	x.addCaseOutputValue("nearWall", "Focus", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos))
	x.addCaseOutputValue("nearWall", "LWall", float64(workingCase.WorldCBRComparisonData.StageData.LeftWallPos))
	x.addCaseOutputValue("nearWall", "RWall", float64(workingCase.WorldCBRComparisonData.StageData.RightWallPos))
	x.addCaseOutputValue("moveID", "MID", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CurrentMoveReferenceID))
	x.addCaseOutputValue("roundState", "", float64(workingCase.WorldCBRComparisonData.RoundState))
	x.addCaseOutputValue("pressureMoveID", "MID", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CurrentMoveReferenceID))
	x.addCaseOutputValue("getHit", "hit", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.SelfHit))
	x.addCaseOutputValue("getHit", "guard", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.SelfGuard))
	x.addCaseOutputValue("didHit", "hit", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MoveHit))
	x.addCaseOutputValue("didHit", "guard", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MoveGuarded))
	x.addCaseOutputValue("frameAdv", "", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.FrameAdv))
	x.addCaseOutputValue("enemyAirborneState", "eAir", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.MStateAir))
	x.addCaseOutputValue("enemyLyingDownState", "eLying", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.MStateLying))
	x.addCaseOutputValue("enemyHitState", "eHit", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateHit))
	x.addCaseOutputValue("enemyHitState", "eHitStun", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.HitStun))
	x.addCaseOutputValue("enemyBlockState", "eBlock", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateHit))
	x.addCaseOutputValue("enemyBlockState", "eBlockStun", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Blockstun))
	x.addCaseOutputValue("enemyAttackState", "eAtk", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateAttack))
	x.addCaseOutputValue("enemyMoveID", "eMID", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CurrentMoveReferenceID))
	x.addCaseOutputValue("pressureEnemyMoveID", "eMID", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CurrentMoveReferenceID))
	x.addCaseOutputValue("yVelocity", "", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Velocity.YVel))
	x.addCaseOutputValue("xVelocity", "", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Velocity.XVel))
	x.addCaseOutputValue("enemyXVelocity", "", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Velocity.XVel))
	x.addCaseOutputValue("enemyYVelocity", "", float64(workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Velocity.YVel))
	x.addCaseOutputValue("frameAdvInitiator", "", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.FrameAdv))
	x.addCaseOutputValue("comboSimilarity", "", float64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.ComboMovesUsed))
	x.addCaseOutputValue("comboSimilarity", "", boolToFloat64(workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Pressure))

	if workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.GenericVars != nil {
		for key, val := range workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.GenericVars.GenericFloatVars {
			x.addCaseOutputValue("ikemenVar", fmt.Sprintf("FVar%v", key), float64(val))
		}
		for key, val := range workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.GenericVars.GenericIntVars {
			x.addCaseOutputValue("ikemenVar", fmt.Sprintf("IVar%v", key), float64(val))
		}
	}

	for _, val2 := range ikemenVarImport[refIds.curFocusCharNr] {
		if val2.VarNr == -1 {
			match := false
			for _, val := range workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {
				if val2.HelperID == val.HelperID {
					match = true
				}
			}
			x.addCaseOutputValue("ikemenVar", fmt.Sprintf("H%v", val2.HelperID), boolToFloat64(match))
		}
	}

	for i, val := range workingCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {

		if val.GenericVars != nil {
			for key, val2 := range val.GenericVars.GenericFloatVars {
				x.addCaseOutputValue("ikemenVar", fmt.Sprintf("H%vFVar%v", val.HelperID, key), float64(val2))
			}
			for key, val2 := range val.GenericVars.GenericIntVars {
				x.addCaseOutputValue("ikemenVar", fmt.Sprintf("H%vIVar%v", val.HelperID, key), float64(val2))
			}
		}

		x.addCaseOutputValue(fmt.Sprintf("Helper_%v_RelativePositionX", i), "PosX", float64(val.PositionX))
		x.addCaseOutputValue(fmt.Sprintf("Helper_%v_RelativePositionY", i), "PosY", float64(val.PositionY))
		x.addCaseOutputValue(fmt.Sprintf("Helper_%v_VelocityX", i), "VelX", float64(val.Velocity.XVel))
		x.addCaseOutputValue(fmt.Sprintf("Helper_%v_VelocityY", i), "VelY", float64(val.Velocity.YVel))
	}
	for i, val := range workingCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData {
		x.addCaseOutputValue(fmt.Sprintf("eHelper_%v_RelativePositionX", i), "PosX", float64(val.PositionX))
		x.addCaseOutputValue(fmt.Sprintf("eHelper_%v_RelativePositionY", i), "PosY", float64(val.PositionY))
		x.addCaseOutputValue(fmt.Sprintf("eHelper_%v_VelocityX", i), "VelX", float64(val.Velocity.XVel))
		x.addCaseOutputValue(fmt.Sprintf("eHelper_%v_VelocityY", i), "VelY", float64(val.Velocity.YVel))
	}

}

func GetDebugText() string {
	return debugUtil.debugText
}

func resetDebug(saveLocation string) {
	if debugUtil.debugActive == true {
		dt := time.Now()
		debugUtil.printDebugData()
		saveDebugData(debugUtil.debugText, saveLocation, dt.Format("01-02-2006_15-04"))
		debugUtil.debugText = ""
		debugUtil.comparisons = []comparisonInstance{}
		debugUtil.inWorkCase = &caseDebugData{}
	}
}
func GetFramesSinceStartRecording() string {
	return fmt.Sprintf("Rec: %v", CBRinter.replayFrames)
}
func GetFramesSinceStartReplaying() string {
	return fmt.Sprintf("Rep: %v", cbrProcess.framesSinceStart)
}
func GetFramesSinceStart() string {
	retString := ""
	if aiData.replaying == true {
		retString += GetFramesSinceStartReplaying() + " "
	}
	if aiData.recording == true {
		retString += GetFramesSinceStartRecording()
	}
	return retString
}

func (x *NavMapString_Float32) debugMapAdd(key string, value float32) {
	return
	if !debugUtil.debugActive {
		return
	}
	x.Add(key, value)
}
func (debugUtil *debugComparisons) addHelperMapping(helperMaps []HelperMapping) {
	if !debugUtil.debugActive {
		return
	}
	for _, val := range helperMaps {
		for _, key := range val.debugMap.keys {
			val2 := val.debugMap.m[key]
			debugUtil.addCompOutputValue(key, val2)
		}
	}

}

type stateMemory struct {
	moveID      int64
	attackState bool
}

var stateMem = stateMemory{moveID: -1, attackState: false}

func (x *CBRRawFrames_Frame) addCaseDebugInfo(frame *CBRRawFrames_Frame, focusCharID int) {
	descriptor := []*DebugDescriptor{}

	if stateMem.moveID == -1 {
		stateMem.moveID = frame.CharData[focusCharID].ComparisonData.CurrentMoveReferenceID
		stateMem.attackState = frame.CharData[focusCharID].ComparisonData.AStateAttack

	} else if stateMem.moveID != frame.CharData[focusCharID].ComparisonData.CurrentMoveReferenceID || frame.CharData[focusCharID].IkemenMoveFrame == 1 {
		stateMem.moveID = frame.CharData[focusCharID].ComparisonData.CurrentMoveReferenceID
		if frame.CharData[focusCharID].ComparisonData.AStateAttack == true {
			descriptor = append(descriptor, &DebugDescriptor{Primary: "Attack", Secondary: stateMem.moveID, Tertiary: int64(frame.CharData[focusCharID].IkemenMoveID)})

		} else {
			descriptor = append(descriptor, &DebugDescriptor{Primary: "StateChange", Secondary: stateMem.moveID, Tertiary: int64(frame.CharData[focusCharID].IkemenMoveID)})

		}
	}

	if frame.CharData[focusCharID].ComparisonData.SelfHit == true {
		descriptor = append(descriptor, &DebugDescriptor{Primary: "SelfHit"})
	}
	if frame.CharData[focusCharID].ComparisonData.SelfGuard == true {
		descriptor = append(descriptor, &DebugDescriptor{Primary: "SelfGuard"})
	}
	if frame.CharData[focusCharID].ComparisonData.MoveHit == true {
		descriptor = append(descriptor, &DebugDescriptor{Primary: "MoveHit"})
	}
	if frame.CharData[focusCharID].ComparisonData.MoveGuarded == true {
		descriptor = append(descriptor, &DebugDescriptor{Primary: "MoveGuarded"})
	}

	hitEnemy := false
	for i, char := range frame.CharData {
		if i != focusCharID && char.ComparisonData.Blockstun > 0 || char.ComparisonData.HitStun > 0 {
			hitEnemy = true
			break
		}
	}
	if stateMem.attackState == true && frame.CharData[focusCharID].ComparisonData.AStateAttack == false &&
		frame.CharData[focusCharID].ComparisonData.MoveHit == false && frame.CharData[focusCharID].ComparisonData.MoveGuarded == false &&
		!hitEnemy {
		descriptor = append(descriptor, &DebugDescriptor{Primary: "MoveWhiff"})
	}
	stateMem.attackState = frame.CharData[focusCharID].ComparisonData.AStateAttack

	for _, helper := range frame.CharData[focusCharID].HelperData {
		if helper.MoveHit == true {
			descriptor = append(descriptor, &DebugDescriptor{Primary: "HelperMoveHit", Secondary: int64(helper.CompData.HelperID)})
		}
		if helper.MoveGuarded == true {
			descriptor = append(descriptor, &DebugDescriptor{Primary: "HelperMoveGuarded", Secondary: int64(helper.CompData.HelperID)})
		}
	}
	x.CharData[focusCharID].DebugDescriptors = descriptor
}

func (x *CBRData_Case) addCaseDebugInfo(replay *CBRRawFrames_ReplayFile, focusCharID int, framestart int, frameend int) {
	if !debugUtil.debugActive {
		return
	}
	x.DebugDescriptors = nil
	for i := framestart; i <= frameend; i++ {
		x.DebugDescriptors = append(x.DebugDescriptors, replay.Frame[i].CharData[focusCharID].DebugDescriptors...)
	}
}
