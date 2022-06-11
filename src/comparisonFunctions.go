package main

import (
	"fmt"
	"math"
)

/*---FILE DESCRIPTION---
Contains comparison functions that are used to determine which case is the most similar to the current
gamestate. Similarity is measured from identical to completly dissimilar (0 to 1), meaning its a dissimilarity metric.
Currently designed to only work in 1v1 scenarios but can be expanded upon.
---FILE DESCRIPTION---*/

//used to keep track of which character is the CBR AI controlled and which is the opponent
type CharReferenceIDs struct {
	curFocusCharNr  int
	curEnemyCharNr  int
	caseFocusCharNr int
	caseEnemyCharNr int
	focusID         string
	enemyID         string
	sameEnemy       bool
}

type HelperMapping struct {
	curHelperIndex  int
	caseHelperIndex int
	comparisonValue float32
	debugMap        NavMapString_Float32
}

//main comparison functions containing and executing all sub comparison functions.
func comparisonFunction(curGamestate *CBRRawFrames, caseData *CBRData, replayIndex int, caseIndex int, framedata *Framedata, refIds *CharReferenceIDs) (comparisonValue float32) {

	debugUtil.resetWorkingCase()
	debugUtil.addCompOutputs()
	debugUtil.addCaseData(caseData.ReplayFile[replayIndex].Case[caseIndex])
	debugUtil.addReplayIndex(replayIndex)
	debugUtil.addCaseIndex(caseIndex)

	workingFrame := curGamestate.ReplayFile[0].Frame[len(curGamestate.ReplayFile[0].Frame)-1]
	workingCase := caseData.ReplayFile[replayIndex].Case[caseIndex]

	//runs all comparison values and adds them together. The result of all comparison values added up should be a value ranging from 0-1.
	//The cbrParameters.cps are used to determine how much a single sub-comparison function affects the final comparisonValue
	var valueBuffer float32 = 0
	comparisonValue = 0

	valueBuffer = relativeXPositionComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.XRelativePosition
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("xRelativePosition", valueBuffer)

	valueBuffer = relativeYPositionComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.YRelativePosition
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("yRelativePosition", valueBuffer)

	valueBuffer = yVelocityComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.yVelocityComparison
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("yVelocity", valueBuffer)

	valueBuffer = xVelocityComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.xVelocityComparison
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("xVelocity", valueBuffer)

	valueBuffer = inputBufferButtonComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.inputBufferButton
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("inputBufferButton", valueBuffer)

	valueBuffer = inputBufferDirectionComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.inputBufferDirection
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("inputBufferDirection", valueBuffer)

	valueBuffer = airborneStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.airborneState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("airborneState", valueBuffer)

	valueBuffer = lyingDownStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.lyingDownState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("lyingDownState", valueBuffer)

	valueBuffer = hitStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.hitState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("hitState", valueBuffer)

	valueBuffer = blockStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.blockState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("blockState", valueBuffer)

	valueBuffer = attackStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.attackState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("attackState", valueBuffer)

	valueBuffer = frameAdvComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.frameAdv
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("frameAdv", valueBuffer)

	valueBuffer = frameAdvInitatorComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.frameAdvInitiator
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("frameAdvInitiator", valueBuffer)

	valueBuffer = comboSimilarityComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.comboSimilarity
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("comboSimilarity", valueBuffer)

	valueBuffer = nearWallComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.nearWall
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("nearWall", valueBuffer)

	valueBuffer = moveIDComparison(workingFrame, workingCase, refIds, framedata) * cbrParameters.cps.moveID
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("moveID", valueBuffer)

	valueBuffer = pressureMoveIDComparison(workingFrame, workingCase, refIds, framedata) * cbrParameters.cps.pressureMoveID
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("pressureMoveID", valueBuffer)

	valueBuffer = getHitComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.getHit
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("getHit", valueBuffer)

	valueBuffer = didHitComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.didHit
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("didHit", valueBuffer)

	valueBuffer = ikemenVarComparison(workingFrame, workingCase, refIds, ikemenVarImport)
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("ikemenVar", valueBuffer)

	valueBuffer = caseReuseCost(cbrProcess, caseIndex, replayIndex) * cbrParameters.cps.caseReuse
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("caseReuse", valueBuffer)

	valueBuffer = roundStateCost(workingFrame, workingCase) * cbrParameters.cps.roundState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("roundState", valueBuffer)

	valueBuffer = enemyYVelocityComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyYVelocityComparison
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyYVelocity", valueBuffer)

	valueBuffer = enemyXVelocityComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyXVelocityComparison
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyXVelocity", valueBuffer)

	valueBuffer = enemyAirborneStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyAirborneState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyAirborneState", valueBuffer)

	valueBuffer = enemyLyingDownStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyLyingDownState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyLyingDownState", valueBuffer)

	valueBuffer = enemyHitStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyHitState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyHitState", valueBuffer)

	valueBuffer = enemyBlockStateComparison(workingFrame, workingCase, refIds) * cbrParameters.cps.enemyBlockState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyBlockState", valueBuffer)

	valueBuffer = enemyAttackStateComparison(workingFrame, workingCase, refIds, framedata) * cbrParameters.cps.enemyAttackState
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyAttackState", valueBuffer)

	valueBuffer = enemyMoveIDComparison(workingFrame, workingCase, refIds, framedata) * cbrParameters.cps.enemyMoveID
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyMoveID", valueBuffer)

	valueBuffer = enemyPressureMoveIDComparison(workingFrame, workingCase, refIds, framedata) * cbrParameters.cps.enemyPressureMoveID
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("enemyPressureMoveID", valueBuffer)

	//calculate the comparison value for every mapping possebility and then choosing the best mapping
	enemyHelperMap := enemyHelperMapping(workingFrame, workingCase, refIds, framedata)
	for i := range enemyHelperMap {
		enemyHelperMap[i].debugMap = NavMapString_Float32{m: map[string]float32{}, keys: []string{}}
		valueBuffer = enemyHelperRelativePositionXComparison(workingFrame, workingCase, refIds, enemyHelperMap[i]) * cbrParameters.cps.enemyHelperRelativePositionX
		enemyHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("eHelper_%v_RelativePositionX", enemyHelperMap[i].caseHelperIndex), valueBuffer)
		enemyHelperMap[i].comparisonValue += valueBuffer

		valueBuffer = enemyHelperRelativePositionYComparison(workingFrame, workingCase, refIds, enemyHelperMap[i]) * cbrParameters.cps.enemyHelperRelativePositionY
		enemyHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("eHelper_%v_RelativePositionY", enemyHelperMap[i].caseHelperIndex), valueBuffer)
		enemyHelperMap[i].comparisonValue += valueBuffer

		valueBuffer = enemyHelperVelocityXComparison(workingFrame, workingCase, refIds, enemyHelperMap[i]) * cbrParameters.cps.enemyHelperXVelocityComparison
		enemyHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("eHelper_%v_VelocityX", enemyHelperMap[i].caseHelperIndex), valueBuffer)
		enemyHelperMap[i].comparisonValue += valueBuffer

		valueBuffer = enemyHelperVelocityYComparison(workingFrame, workingCase, refIds, enemyHelperMap[i]) * cbrParameters.cps.enemyHelperYVelocityComparison
		enemyHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("eHelper_%v_VelocityY", enemyHelperMap[i].caseHelperIndex), valueBuffer)
		enemyHelperMap[i].comparisonValue += valueBuffer

	}
	var usedUpCaseIndex = map[int]bool{}
	var usedUpCurIndex = map[int]bool{}
	bufferVal, trimmedEnemyHelperMap := findBestHelperMapping(enemyHelperMap, usedUpCaseIndex, usedUpCurIndex)
	comparisonValue += bufferVal
	debugUtil.addHelperMapping(trimmedEnemyHelperMap)

	//calculate the comparison value for every mapping possebility and then choosing the best mapping
	focusHelperMap := focusCharHelperMapping(workingFrame, workingCase, refIds, framedata)
	for i := range focusHelperMap {
		focusHelperMap[i].debugMap = NavMapString_Float32{m: map[string]float32{}, keys: []string{}}
		valueBuffer = helperRelativePositionXComparison(workingFrame, workingCase, refIds, focusHelperMap[i]) * cbrParameters.cps.helperRelativePositionX
		focusHelperMap[i].comparisonValue += valueBuffer
		focusHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("Helper_%v_RelativePositionX", focusHelperMap[i].caseHelperIndex), valueBuffer)

		valueBuffer = helperRelativePositionYComparison(workingFrame, workingCase, refIds, focusHelperMap[i]) * cbrParameters.cps.helperRelativePositionY
		focusHelperMap[i].comparisonValue += valueBuffer
		focusHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("Helper_%v_RelativePositionY", focusHelperMap[i].caseHelperIndex), valueBuffer)

		valueBuffer = helperVelocityXComparison(workingFrame, workingCase, refIds, focusHelperMap[i]) * cbrParameters.cps.helperXVelocityComparison
		focusHelperMap[i].comparisonValue += valueBuffer
		focusHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("Helper_%v_VelocityX", focusHelperMap[i].caseHelperIndex), valueBuffer)

		valueBuffer = helperVelocityYComparison(workingFrame, workingCase, refIds, focusHelperMap[i]) * cbrParameters.cps.helperYVelocityComparison
		focusHelperMap[i].comparisonValue += valueBuffer
		focusHelperMap[i].debugMap.debugMapAdd(fmt.Sprintf("Helper_%v_VelocityY", focusHelperMap[i].caseHelperIndex), valueBuffer)

	}
	usedUpCaseIndex = map[int]bool{}
	usedUpCurIndex = map[int]bool{}
	bufferVal, trimmedFocusHelperMap := findBestHelperMapping(focusHelperMap, usedUpCaseIndex, usedUpCurIndex)
	comparisonValue += bufferVal
	debugUtil.addHelperMapping(trimmedFocusHelperMap)

	valueBuffer = objectOrderComparison(trimmedFocusHelperMap, trimmedEnemyHelperMap, workingFrame, workingCase, refIds) * cbrParameters.cps.objectOrder
	comparisonValue += valueBuffer
	debugUtil.addCompOutputValue("objectOrder", valueBuffer)

	debugUtil.addCompValue(comparisonValue)
	debugUtil.addComparison()

	return comparisonValue
}

//relative distance to the opponent compared to case
func relativeXPositionComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos
	curEnemyPosX := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.CharPos.XPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos
	caseEnemyPosX := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CharPos.XPos
	return float32(math.Min(math.Abs((math.Abs(float64(curPlayerPosX-curEnemyPosX))-math.Abs(float64(casePlayerPosX-caseEnemyPosX)))/float64(cbrParameters.maxXPositionComparison)), 1))
}

func relativeYPositionComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerPosY := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.YPos
	curEnemyPosY := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.CharPos.YPos
	casePlayerPosY := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.YPos
	caseEnemyPosY := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CharPos.YPos
	return float32(math.Min(math.Abs((math.Abs(float64(curPlayerPosY-curEnemyPosY))-math.Abs(float64(casePlayerPosY-caseEnemyPosY)))/float64(cbrParameters.maxYPositionComparison)), 1))

}

func yVelocityComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerVelY := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.Velocity.YVel
	casePlayerVelY := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Velocity.YVel

	diff := float32(math.Abs(float64(curPlayerVelY - casePlayerVelY)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret

}

func xVelocityComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerVelX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.Velocity.XVel
	casePlayerVelX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Velocity.XVel

	diff := float32(math.Abs(float64(curPlayerVelX - casePlayerVelX)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret

}

func inputBufferDirectionComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	var compValue int32
	compValue = 0
	gamestateBuffer := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.InputBuffer.InputDirection
	for i, val := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.InputBuffer.InputDirection {
		//if the direction input beeing held differs from the case, add max dissimilarity
		if (val <= 0 && gamestateBuffer[i] > 0) || (val > 0 && gamestateBuffer[i] <= 0) {
			compValue += 1
		} else if val > 0 {
			compValue += minInt32(Abs(val-gamestateBuffer[i]), cbrParameters.maxInputBufferDifference) / cbrParameters.maxInputBufferDifference
		}
	}

	return float32(compValue) / float32(len(gamestateBuffer))
}
func inputBufferButtonComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	var compValue int32
	compValue = 0
	gamestateBuffer := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.InputBuffer.InputButton
	for i, val := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.InputBuffer.InputButton {
		//if the direction input beeing held differs from the case, add max dissimilarity
		if (val <= 0 && gamestateBuffer[i] > 0) || (val > 0 && gamestateBuffer[i] <= 0) {
			compValue += 1
		} else if val > 0 {
			compValue += maxInt32(Abs(val-gamestateBuffer[i]), cbrParameters.maxInputBufferDifference) / cbrParameters.maxInputBufferDifference
		}
	}
	return float32(compValue) / float32(len(gamestateBuffer))
}

func airborneStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curMoveState := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.MStateAir
	caseMoveState := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MStateAir
	if curMoveState != caseMoveState {
		return 1
	}
	return 0
}
func lyingDownStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curMoveState := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.MStateLying
	caseMoveState := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MStateLying
	if curMoveState != caseMoveState {
		return 1
	}
	return 0
}
func getHitComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	hitComp := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.SelfHit == true && singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.SelfHit != true
	blockComp := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.SelfGuard == true && singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.SelfGuard != true
	if hitComp == true || blockComp == true {
		return 1
	}
	return 0
}
func didHitComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	hitComp := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.MoveHit == true && singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MoveHit != true
	blockComp := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.MoveGuarded == true && singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.MoveGuarded != true
	if hitComp == true || blockComp == true {
		return 1
	}
	return 0
}
func blockStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.AStateHit
	caseState := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateHit
	if curState != caseState {
		return 1
	}

	curBlockStun := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.Blockstun
	caseBlockStun := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Blockstun

	if (curBlockStun > 0) != (caseBlockStun > 0) {
		return 1
	}

	stunDiff := Abs(curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.Blockstun - singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Blockstun)
	diffNorm := float32(minInt32(cbrParameters.maxBlockstunDifference, stunDiff)) / float32(cbrParameters.maxBlockstunDifference)
	return diffNorm
}
func hitStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.AStateHit
	caseState := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateHit
	if curState != caseState {
		return 1
	}

	curBlockStun := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.HitStun
	caseBlockStun := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.HitStun

	if (curBlockStun > 0) != (caseBlockStun > 0) {
		return 1
	}

	stunDiff := Abs(curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.HitStun - singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.HitStun)
	diffNorm := float32(minInt32(cbrParameters.maxHitstunDifference, stunDiff)) / float32(cbrParameters.maxHitstunDifference)
	return diffNorm
}

func attackStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.AStateHit
	caseState := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.AStateHit
	if curState != caseState {
		return 1
	}

	return 0
}

func frameAdvComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curAdv := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.FrameAdv
	if curAdv != 0 {
		curAdv = curAdv / Abs(curAdv)
	}

	caseAdv := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.FrameAdv
	if caseAdv != 0 {
		caseAdv = caseAdv / Abs(caseAdv)
	}
	return float32(Abs(curAdv-caseAdv)) / 2
}

func frameAdvInitatorComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curInitiator := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.FrameAdv
	caseInitiator := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.FrameAdv

	if curInitiator != caseInitiator {
		return 1
	}
	return 0
}

func comboSimilarityComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curMovesUsed := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.ComboMovesUsed
	curPressure := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.Pressure
	caseMovesUsed := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.ComboMovesUsed
	casePressure := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.Pressure

	if (curMovesUsed > 0) != (caseMovesUsed > 0) {
		return 1
	}
	if curPressure != casePressure {
		return 1
	}

	diff := float32(Abs(curMovesUsed - caseMovesUsed))
	ret := minFloat32(1, diff/cbrParameters.comboLength)

	return ret
}

func nearWallComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos
	curStageSize := calcStageSize(*curGamestate.WorldCBRComparisonData.StageData)
	caseStageSize := calcStageSize(*singleCase.WorldCBRComparisonData.StageData)

	var curFacingWall float64 = 0
	var curBackWall float64 = 0
	var caseFacingWall float64 = 0
	var caseBackWall float64 = 0
	if curGamestate.CharData[refIds.curFocusCharNr].Facing {
		curFacingWall = math.Abs(float64(curGamestate.WorldCBRComparisonData.StageData.RightWallPos - curPlayerPosX))
		curBackWall = math.Abs(float64(curGamestate.WorldCBRComparisonData.StageData.LeftWallPos - curPlayerPosX))
		caseFacingWall = math.Abs(float64(singleCase.WorldCBRComparisonData.StageData.RightWallPos - casePlayerPosX))
		caseBackWall = math.Abs(float64(singleCase.WorldCBRComparisonData.StageData.LeftWallPos - casePlayerPosX))
	} else {
		curBackWall = math.Abs(float64(curGamestate.WorldCBRComparisonData.StageData.RightWallPos - curPlayerPosX))
		curFacingWall = math.Abs(float64(curGamestate.WorldCBRComparisonData.StageData.LeftWallPos - curPlayerPosX))
		caseBackWall = math.Abs(float64(singleCase.WorldCBRComparisonData.StageData.RightWallPos - casePlayerPosX))
		caseFacingWall = math.Abs(float64(singleCase.WorldCBRComparisonData.StageData.LeftWallPos - casePlayerPosX))
	}

	return float32(math.Abs(curBackWall-caseBackWall))/curStageSize + float32(math.Abs(curFacingWall-caseFacingWall))/caseStageSize
}

/*
func oldnearWallComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].cPositionX
	casePlayerPosX := singleCase.CharPos[refIds.caseFocusCharNr].XPos
	curStageSize := calcStageSize(*curGamestate.StageData)
	caseStageSize := calcStageSize(*singleCase.StageData)
	curNearLeftWall := math.Abs(float64(curGamestate.StageData.LeftWallPos - curPlayerPosX)) <= float64(cbrParameters.nearWallDist*curStageSize)
	curNearRightWall := math.Abs(float64(curGamestate.StageData.RightWallPos - curPlayerPosX)) <= float64(cbrParameters.nearWallDist*curStageSize)
	caseNearLeftWall := math.Abs(float64(singleCase.StageData.LeftWallPos - casePlayerPosX)) <= float64(cbrParameters.nearWallDist*caseStageSize)
	caseNearRightWall := math.Abs(float64(singleCase.StageData.RightWallPos - casePlayerPosX)) <= float64(cbrParameters.nearWallDist*caseStageSize)

	if (curNearLeftWall || curNearRightWall) != (caseNearLeftWall || caseNearRightWall){
		return 1
	}

	if curNearLeftWall == false  && curNearRightWall == false && caseNearLeftWall == false && caseNearRightWall == false{
		return 0
	}

	//check if the AI is with their back to the wall or facing the wall
	curAgressor := curGamestate.CharData[refIds.curFocusCharNr].Facing == true && curNearRightWall || curGamestate.CharData[refIds.curFocusCharNr].Facing == false || curNearLeftWall
	caseAgressor := singleCase.Facing == true && caseNearRightWall || singleCase.Facing == false || caseNearLeftWall
	if curAgressor != caseAgressor{
		return 1
	}

	return 0
}

*/

func moveIDComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (comparisonScore float32) {
	curMoveRefId := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CurrentMoveReferenceID
	if curMoveRefId != singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CurrentMoveReferenceID {
		return 1
	}
	/* well add a more complicated logic for pressure/whiff/block checkinglater baed on the code below
	charFrameData, ok := framedata.CharData[refIds.enemyID]
	if !ok {
		return 0
	}

	moveFrameData, ok := charFrameData.Movedata[curMoveRefId]
	if !ok{
		return 0
	}
	enemyMoveWhiffed := moveFrameData.Startup + moveFrameData.Active >= curGamestate.CharData[refIds.curEnemyCharNr].CurrentMoveFrame
	*/
	return 0
}

func ikemenVarComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, ikemenVars map[int][]CBRVariableImportance) (comparisonScore float32) {
	var dissimilarityCounter float32 = 0
	for _, val := range ikemenVars[refIds.curFocusCharNr] {
		if val.VarNr == -1 && val.HelperID > 0 {
			curHelperFound := false
			caseHelperFound := false
			for _, helper := range curGamestate.CharData[refIds.curFocusCharNr].HelperData {
				if helper.CompData.HelperID == val.HelperID {
					curHelperFound = true
				}
			}
			for _, helper := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {
				if helper.HelperID == val.HelperID {
					caseHelperFound = true
				}
			}
			if curHelperFound != caseHelperFound {
				dissimilarityCounter += val.MaxDissimilarityCost
			}
			continue
		}

		if val.HelperID < 0 {
			if val.Float == true {
				curFloat := curGamestate.CharData[refIds.curFocusCharNr].GenericFloatVars[val.VarNr]
				caseFloat := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.GenericVars.GenericFloatVars[val.VarNr]
				diffFloat := math.Abs(float64(curFloat - caseFloat))
				resultFloat := math.Floor(diffFloat/val.VariableIncrements) / val.DivisionIntervals
				finalFloat := minFloat32(float32(resultFloat)*val.MaxDissimilarityCost, val.MaxDissimilarityCost)
				dissimilarityCounter += finalFloat
			} else {
				curInt := curGamestate.CharData[refIds.curFocusCharNr].GenericIntVars[val.VarNr]
				caseInt := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.GenericVars.GenericIntVars[val.VarNr]
				diffFloat := math.Abs(float64(curInt - caseInt))
				resultFloat := math.Floor(diffFloat/val.VariableIncrements) / val.DivisionIntervals
				finalFloat := minFloat32(float32(resultFloat)*val.MaxDissimilarityCost, val.MaxDissimilarityCost)
				dissimilarityCounter += finalFloat
			}
		} else {
			if val.Float == true {
				foundCur, foundCase := false, false
				var curFloat float32 = 0
				var caseFloat float32 = 0
				for _, helper := range curGamestate.CharData[refIds.curFocusCharNr].HelperData {
					if helper.CompData.HelperID == val.HelperID {
						curFloat = helper.GenericFloatVars[val.VarNr]
						foundCur = true
					}
				}
				for _, helper := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {
					if helper.HelperID == val.HelperID {
						caseFloat = helper.GenericVars.GenericFloatVars[val.VarNr]
						foundCase = true
					}
				}
				if foundCase && foundCur {
					diffFloat := math.Abs(float64(curFloat - caseFloat))
					resultFloat := math.Floor(diffFloat/val.VariableIncrements) / val.DivisionIntervals
					finalFloat := minFloat32(float32(resultFloat)*val.MaxDissimilarityCost, val.MaxDissimilarityCost)
					dissimilarityCounter += finalFloat
				} else {
					if foundCase {
						dissimilarityCounter += val.MaxDissimilarityCost
					}
				}

			} else {
				foundCur, foundCase := false, false
				var curInt int32 = 0
				var caseInt int32 = 0
				for _, helper := range curGamestate.CharData[refIds.curFocusCharNr].HelperData {
					if helper.CompData.HelperID == val.HelperID {
						curInt = helper.GenericIntVars[val.VarNr]
						foundCur = true
					}
				}
				for _, helper := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {
					if helper.HelperID == val.HelperID {
						caseInt = helper.GenericVars.GenericIntVars[val.VarNr]
						foundCase = true
					}
				}
				if foundCase && foundCur {
					diffFloat := math.Abs(float64(curInt - caseInt))
					resultFloat := math.Floor(diffFloat/val.VariableIncrements) / val.DivisionIntervals
					finalFloat := minFloat32(float32(resultFloat)*val.MaxDissimilarityCost, val.MaxDissimilarityCost)
					dissimilarityCounter += finalFloat
				} else {
					if foundCase {
						dissimilarityCounter += val.MaxDissimilarityCost
					}
				}
			}
		}
	}

	return dissimilarityCounter
}

func pressureMoveIDComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (comparisonScore float32) {
	curMoveRefId := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CurrentMoveReferenceID
	curPressure := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.AStateAttack && curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.AStateHit
	if curPressure == false {
		return 0
	}
	if curMoveRefId != singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CurrentMoveReferenceID {
		return 1
	}
	/* well add a more complicated logic for pressure/whiff/block checkinglater baed on the code below
	charFrameData, ok := framedata.CharData[refIds.enemyID]
	if !ok {
		return 0
	}

	moveFrameData, ok := charFrameData.Movedata[curMoveRefId]
	if !ok{
		return 0
	}
	enemyMoveWhiffed := moveFrameData.Startup + moveFrameData.Active >= curGamestate.CharData[refIds.curEnemyCharNr].CurrentMoveFrame
	*/
	return 0
}

func enemyYVelocityComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerVelY := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.Velocity.YVel
	casePlayerVelY := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Velocity.YVel

	diff := float32(math.Abs(float64(curPlayerVelY - casePlayerVelY)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret

}

func enemyXVelocityComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curPlayerVelX := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.Velocity.XVel
	casePlayerVelX := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Velocity.XVel

	diff := float32(math.Abs(float64(curPlayerVelX - casePlayerVelX)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret

}

func enemyAirborneStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.MStateAir
	caseState := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.MStateAir
	if curState != caseState {
		return 1
	}
	return 0
}
func enemyLyingDownStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.MStateLying
	caseState := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.MStateLying
	if curState != caseState {
		return 1
	}
	return 0
}
func enemyBlockStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.AStateHit
	caseState := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateHit
	if curState != caseState {
		return 1
	}
	//return 0

	curBlockStun := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.Blockstun
	caseBlockStun := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Blockstun

	if (curBlockStun > 0) != (caseBlockStun > 0) {
		return 1
	}

	stunDiff := Abs(curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.Blockstun - singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.Blockstun)
	diffNorm := float32(minInt32(cbrParameters.maxBlockstunDifference, stunDiff)) / float32(cbrParameters.maxBlockstunDifference)
	return diffNorm
}
func enemyHitStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.AStateHit
	caseState := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateHit
	if curState != caseState {
		return 1
	}
	//return 0

	curBlockStun := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.HitStun
	caseBlockStun := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.HitStun

	if (curBlockStun > 0) != (caseBlockStun > 0) {
		return 1
	}

	stunDiff := Abs(curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.HitStun - singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.HitStun)
	diffNorm := float32(minInt32(cbrParameters.maxHitstunDifference, stunDiff)) / float32(cbrParameters.maxHitstunDifference)
	return diffNorm
}

func caseReuseCost(process CBRProcess, caseIndex int, replayIndex int) (comparisonScore float32) {
	val, ok := process.caseUsageReplayFile[replayIndex]
	if !ok {
		return 0
	}
	val2, ok := val[caseIndex]
	if !ok {
		return 0
	}

	repetition := float64(val2.timesUsed * cbrParameters.repetitionFrames)
	if repetition == 0 {
		return 0
	}
	diff := math.Max(0, repetition-float64(process.framesSinceStart-val2.lastFrameUsed))
	return float32(diff / repetition)
}
func roundStateCost(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case) (comparisonScore float32) {
	if curGamestate.WorldCBRComparisonData.RoundState == singleCase.WorldCBRComparisonData.RoundState {
		return 0
	}
	return 1
}
func enemyAttackStateComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (comparisonScore float32) {
	curState := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.AStateAttack
	caseState := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.AStateAttack
	if curState != caseState {
		return 1
	}
	return 0

}

func enemyMoveIDComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (comparisonScore float32) {
	if refIds.sameEnemy == false {
		return 1
	}
	curMoveRefId := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.CurrentMoveReferenceID
	if curMoveRefId != singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CurrentMoveReferenceID {
		return 1
	}
	/* well add a more complicated logic for pressure/whiff/block checkinglater baed on the code below
	charFrameData, ok := framedata.CharData[refIds.enemyID]
	if !ok {
		return 0
	}

	moveFrameData, ok := charFrameData.Movedata[curMoveRefId]
	if !ok{
		return 0
	}
	enemyMoveWhiffed := moveFrameData.Startup + moveFrameData.Active >= curGamestate.CharData[refIds.curEnemyCharNr].CurrentMoveFrame
	*/
	return 0
}
func enemyPressureMoveIDComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (comparisonScore float32) {
	curPressure := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.AStateAttack && curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.AStateHit
	if curPressure == false {
		return 0
	}
	if refIds.sameEnemy == false {
		return 1
	}
	curMoveRefId := curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.CurrentMoveReferenceID
	if curMoveRefId != singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CurrentMoveReferenceID {
		return 1
	}
	/* well add a more complicated logic for pressure/whiff/block checkinglater baed on the code below
	charFrameData, ok := framedata.CharData[refIds.enemyID]
	if !ok {
		return 0
	}

	moveFrameData, ok := charFrameData.Movedata[curMoveRefId]
	if !ok{
		return 0
	}
	enemyMoveWhiffed := moveFrameData.Startup + moveFrameData.Active >= curGamestate.CharData[refIds.curEnemyCharNr].CurrentMoveFrame
	*/
	return 0
}

//relative distance to the opponent compared to case
func enemyHelperRelativePositionXComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curHelperPosX := curGamestate.CharData[refIds.curEnemyCharNr].HelperData[hMap.curHelperIndex].CompData.PositionX
	caseHelperPosX := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData[hMap.caseHelperIndex].PositionX

	return float32(math.Abs((math.Abs(float64(curPlayerPosX-curHelperPosX)) - math.Abs(float64(casePlayerPosX-caseHelperPosX))) / float64(cbrParameters.maxXPositionComparison)))
}

//relative distance to the opponent compared to case
func enemyHelperRelativePositionYComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.YPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.YPos

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curHelperPosY := curGamestate.CharData[refIds.curEnemyCharNr].HelperData[hMap.curHelperIndex].CompData.PositionY
	caseHelperPosY := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData[hMap.caseHelperIndex].PositionY

	return float32(math.Abs((math.Abs(float64(curPlayerPosX-curHelperPosY)) - math.Abs(float64(casePlayerPosX-caseHelperPosY))) / float64(cbrParameters.maxXPositionComparison)))
}

//relative distance to the opponent compared to case
func helperRelativePositionXComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos
	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curHelperPosX := curGamestate.CharData[refIds.curFocusCharNr].HelperData[hMap.curHelperIndex].CompData.PositionX
	caseHelperPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData[hMap.caseHelperIndex].PositionX

	return float32(math.Abs((math.Abs(float64(curPlayerPosX-curHelperPosX)) - math.Abs(float64(casePlayerPosX-caseHelperPosX))) / float64(cbrParameters.maxXPositionComparison)))
}

//relative distance to the opponent compared to case
func helperRelativePositionYComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {
	curPlayerPosX := curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.YPos
	casePlayerPosX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.YPos

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curHelperPosY := curGamestate.CharData[refIds.curFocusCharNr].HelperData[hMap.curHelperIndex].CompData.PositionY
	caseHelperPosY := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData[hMap.caseHelperIndex].PositionY

	return float32(math.Abs((math.Abs(float64(curPlayerPosX-curHelperPosY)) - math.Abs(float64(casePlayerPosX-caseHelperPosY))) / float64(cbrParameters.maxXPositionComparison)))
}

func helperVelocityXComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curVelX := curGamestate.CharData[refIds.curFocusCharNr].HelperData[hMap.curHelperIndex].CompData.Velocity.XVel
	caseVelX := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData[hMap.caseHelperIndex].Velocity.XVel

	diff := float32(math.Abs(float64(curVelX - caseVelX)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret
}

func helperVelocityYComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curVelY := curGamestate.CharData[refIds.curFocusCharNr].HelperData[hMap.curHelperIndex].CompData.Velocity.YVel
	caseVelY := singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData[hMap.caseHelperIndex].Velocity.YVel

	diff := float32(math.Abs(float64(curVelY - caseVelY)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret
}
func enemyHelperVelocityXComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curVelX := curGamestate.CharData[refIds.curEnemyCharNr].HelperData[hMap.curHelperIndex].CompData.Velocity.XVel
	caseVelX := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData[hMap.caseHelperIndex].Velocity.XVel

	diff := float32(math.Abs(float64(curVelX - caseVelX)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret
}
func enemyHelperVelocityYComparison(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, hMap HelperMapping) (comparisonScore float32) {

	if hMap.caseHelperIndex == -1 || hMap.curHelperIndex == -1 {
		return 1
	}
	curVelY := curGamestate.CharData[refIds.curEnemyCharNr].HelperData[hMap.curHelperIndex].CompData.Velocity.YVel
	caseVelY := singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData[hMap.caseHelperIndex].Velocity.YVel

	diff := float32(math.Abs(float64(curVelY - caseVelY)))
	div := diff / cbrParameters.maxVelocityComparison
	ret := minFloat32(1, div)

	return ret
}

//Since we somehow have to know which helper to compare to which helper, and some helpers might exist multiple times we need to create a mapping
func enemyHelperMapping(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (hMap []HelperMapping) {
	hMap = []HelperMapping{}
	var curHelperIndex = map[int]int32{}
	var missedIndex = map[int]bool{}

	for i, helper := range curGamestate.CharData[refIds.curEnemyCharNr].HelperData {
		_, ok := framedata.CharData[refIds.enemyID]
		if ok {
			_, ok = framedata.CharData[refIds.enemyID].ProjectileData[helper.CompData.CurrentMoveReferenceID]
		}
		if helper.CompData.AStateAttack || ok || helper.CompData.HitboxOut || helper.CompData.HurtboxOut {
			curHelperIndex[i] = helper.CompData.HelperID
			missedIndex[i] = true
		}
	}

	for i, caseHelper := range singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData {
		_, ok := framedata.CharData[refIds.enemyID]
		if ok {
			_, ok = framedata.CharData[refIds.enemyID].ProjectileData[caseHelper.CurrentMoveReferenceID]
		}
		if caseHelper.AStateAttack || ok || caseHelper.HitboxOut || caseHelper.HurtboxOut {
			foundMatch := false
			for j, val := range curHelperIndex {
				if val == caseHelper.HelperID {
					foundMatch = true
					hMap = append(hMap, HelperMapping{curHelperIndex: j, caseHelperIndex: i, comparisonValue: 0})
					missedIndex[j] = false
				}
			}
			if foundMatch == false {
				hMap = append(hMap, HelperMapping{curHelperIndex: -1, caseHelperIndex: i, comparisonValue: 0})
			}

		}
	}

	for i, bool := range missedIndex {
		if bool == true {
			hMap = append(hMap, HelperMapping{curHelperIndex: i, caseHelperIndex: -1, comparisonValue: 0})
		}
	}

	return hMap
}

//Since we somehow have to know which helper to compare to which helper, and some helpers might exist multiple times we need to create a mapping
func focusCharHelperMapping(curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs, framedata *Framedata) (hMap []HelperMapping) {
	hMap = []HelperMapping{}
	var curHelperIndex = map[int]int32{}
	var missedIndex = map[int]bool{}

	for i, helper := range curGamestate.CharData[refIds.curFocusCharNr].HelperData {
		_, ok := framedata.CharData[refIds.focusID]
		if ok {
			_, ok = framedata.CharData[refIds.focusID].ProjectileData[helper.CompData.CurrentMoveReferenceID]
		}
		if helper.CompData.AStateAttack || ok || helper.CompData.HitboxOut || helper.CompData.HurtboxOut {
			curHelperIndex[i] = helper.CompData.HelperID
			missedIndex[i] = true
		}
	}

	for i, caseHelper := range singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData {
		_, ok := framedata.CharData[refIds.focusID]
		if ok {
			_, ok = framedata.CharData[refIds.focusID].ProjectileData[caseHelper.CurrentMoveReferenceID]
		}
		if caseHelper.AStateAttack || ok || caseHelper.HitboxOut || caseHelper.HurtboxOut {
			foundMatch := false
			for j, val := range curHelperIndex {
				if val == caseHelper.HelperID {
					foundMatch = true
					hMap = append(hMap, HelperMapping{curHelperIndex: j, caseHelperIndex: i, comparisonValue: 0})
					missedIndex[j] = false
				}
			}
			if foundMatch == false {
				hMap = append(hMap, HelperMapping{curHelperIndex: -1, caseHelperIndex: i, comparisonValue: 0})
			}

		}
	}

	for i, bool := range missedIndex {
		if bool == true {
			hMap = append(hMap, HelperMapping{curHelperIndex: i, caseHelperIndex: -1, comparisonValue: 0})
		}
	}

	return hMap
}

//calculate the comparison value for every mapping possebility and then choosing the best mapping
//Basically traveling salesman problem, for now its just a brute force solution
func findBestHelperMapping(hMap []HelperMapping, usedUpCaseIndex map[int]bool, usedUpCurIndex map[int]bool) (bestCompVal float32, filteredMap []HelperMapping) {
	startingIndex := -2
	bestCompVal = 0
	filteredMap = []HelperMapping{}

	for _, val := range hMap {
		_, ok := usedUpCaseIndex[val.caseHelperIndex]
		_, ok2 := usedUpCurIndex[val.curHelperIndex]

		if !ok && !ok2 && (val.caseHelperIndex == startingIndex || startingIndex == -2) {
			startingIndex = val.caseHelperIndex
			caseBuffer := usedUpCaseIndex
			curBuffer := usedUpCurIndex
			caseBuffer[val.caseHelperIndex] = true
			curBuffer[val.curHelperIndex] = true

			compVal, bufferMap := findBestHelperMapping(hMap, caseBuffer, curBuffer)
			compVal += +val.comparisonValue
			if compVal < bestCompVal || bestCompVal == 0 {
				bestCompVal = compVal
				filteredMap = append(filteredMap, val)
				filteredMap = append(filteredMap, bufferMap...)
			}
		}
	}

	return bestCompVal, filteredMap

}

type ObjectOrder struct {
	charID   int
	helperID int
	XPos     float32
}

func objectOrderComparison(focusHMap []HelperMapping, enemyHMap []HelperMapping, curGamestate *CBRRawFrames_Frame, singleCase *CBRData_Case, refIds *CharReferenceIDs) (compVal float32) {
	var curOrder []ObjectOrder
	var curOrderReversed []ObjectOrder
	var caseOrder []ObjectOrder
	insertIndex := 0

	curOrder = append(curOrder, ObjectOrder{refIds.curFocusCharNr, -1, curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos})
	curOrderReversed = append(curOrderReversed, ObjectOrder{refIds.curFocusCharNr, -1, curGamestate.CharData[refIds.curFocusCharNr].ComparisonData.CharPos.XPos})

	obj := ObjectOrder{refIds.curEnemyCharNr, -1, curGamestate.CharData[refIds.curEnemyCharNr].ComparisonData.CharPos.XPos}
	insertIndex, curOrder = objectOrderInsert(obj, curOrder)
	insertIndex = (len(curOrder) / 2) - (insertIndex - (len(curOrder) / 2)) - 1
	curOrderReversed = objectOrderInsertAt(obj, curOrderReversed, insertIndex)

	for _, val := range focusHMap {
		if val.curHelperIndex != -1 && val.caseHelperIndex != -1 {
			obj := ObjectOrder{refIds.curFocusCharNr, val.curHelperIndex, curGamestate.CharData[refIds.curFocusCharNr].HelperData[val.curHelperIndex].CompData.PositionX}
			insertIndex, curOrder = objectOrderInsert(obj, curOrder)
			insertIndex = (len(curOrder) / 2) - (insertIndex - (len(curOrder) / 2)) - 1
			curOrderReversed = objectOrderInsertAt(obj, curOrderReversed, insertIndex)
		}
	}
	for _, val := range enemyHMap {
		if val.curHelperIndex != -1 && val.caseHelperIndex != -1 {
			obj := ObjectOrder{refIds.curEnemyCharNr, val.curHelperIndex, curGamestate.CharData[refIds.curEnemyCharNr].HelperData[val.curHelperIndex].CompData.PositionX}
			insertIndex, curOrder = objectOrderInsert(obj, curOrder)
			insertIndex = (len(curOrder) / 2) - (insertIndex - (len(curOrder) / 2)) - 1
			curOrderReversed = objectOrderInsertAt(obj, curOrderReversed, insertIndex)
		}
	}

	//insert case Order
	caseOrder = append(caseOrder, ObjectOrder{refIds.caseFocusCharNr, -1, singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].ComparisonData.CharPos.XPos})

	obj = ObjectOrder{refIds.caseEnemyCharNr, -1, singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].ComparisonData.CharPos.XPos}
	_, caseOrder = objectOrderInsert(obj, caseOrder)

	for _, val := range focusHMap {
		if val.curHelperIndex != -1 && val.caseHelperIndex != -1 {
			obj := ObjectOrder{refIds.caseFocusCharNr, val.caseHelperIndex, singleCase.CharCBRComparisonData[refIds.caseFocusCharNr].HelperData[val.caseHelperIndex].PositionX}
			_, caseOrder = objectOrderInsert(obj, caseOrder)
		}
	}

	for _, val := range enemyHMap {
		if val.curHelperIndex != -1 && val.caseHelperIndex != -1 {
			obj := ObjectOrder{refIds.caseEnemyCharNr, val.caseHelperIndex, singleCase.CharCBRComparisonData[refIds.caseEnemyCharNr].HelperData[val.caseHelperIndex].PositionX}
			_, caseOrder = objectOrderInsert(obj, caseOrder)
		}
	}

	var normalErrors int32 = 0
	var mirroredErrors int32 = 0

	for i := range caseOrder {
		if curOrder[i].helperID != caseOrder[i].helperID || curOrder[i].charID != caseOrder[i].charID {
			normalErrors++
		}
		if curOrderReversed[i].helperID != caseOrder[i].helperID || curOrderReversed[i].charID != caseOrder[i].charID {
			mirroredErrors++
		}
	}

	return float32(Min(normalErrors, mirroredErrors)) / float32(len(caseOrder))

}

func objectOrderInsert(obj ObjectOrder, arr []ObjectOrder) (index int, arrOut []ObjectOrder) {
	index = 0
	for i, val := range arr {
		if obj.XPos > val.XPos {
			index = i + 1
			break
		}
	}
	temp := append([]ObjectOrder{}, arr[index:]...)
	arr = append(arr[0:index], obj)
	arr = append(arr, temp...)

	return index, arr
}
func objectOrderInsertAt(obj ObjectOrder, arr []ObjectOrder, index int) (arrOut []ObjectOrder) {
	index = 0
	temp := append([]ObjectOrder{}, arr[index:]...)
	arr = append(arr[0:index], obj)
	arr = append(arr, temp...)

	return arr
}

func calcStageSize(data StageData) float32 {
	return float32(math.Abs(float64(data.LeftWallPos - data.RightWallPos)))
}

func calcStageMiddle(data *StageData) float32 {
	return (data.LeftWallPos + data.RightWallPos) / 2
}

func mirrorPosition(data *StageData, xPos float32) float32 {
	middle := calcStageMiddle(data)
	return middle + (middle - xPos)
}

func maxInt32(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
func minInt32(a int32, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func minFloat32(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
