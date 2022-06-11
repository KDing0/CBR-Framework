package main

import "math/rand"

/*---FILE DESCRIPTION---
Contains thelogic which determines when a new case needs to be selected for the AI.
The determination of which case is best is done in "comparisonFunctions.go".
Checks for a new case if the current case is done, or if the gamestate changes drastically enough.
---FILE DESCRIPTION---*/

//CBRProcess is used to during case selection to store data that needs to be continously updated, like the currently running case.
type CBRProcess struct {
	curScope            int
	randSeed            int64
	lastCase            *CBRCaseReference
	caseUsageReplayFile map[int]map[int]CaseUsageData
	framesSinceStart    int64
}

var cbrProcess = CBRProcess{
	curScope:         0,
	randSeed:         0,
	framesSinceStart: 0,
	lastCase: &CBRCaseReference{
		cbrCase:       &CBRData_Case{},
		comparisonVal: -1,
		caseIndex:     -1,
		replayIndex:   -1,
		frameIndex:    -1},
}

type CBRCaseReference struct {
	cbrCase       *CBRData_Case
	comparisonVal float32
	caseIndex     int
	replayIndex   int
	frameIndex    int32
}

type CaseUsageData struct {
	lastFrameUsed int64
	timesUsed     int64
}

//main AI loop, every frame checks if another case needs to be found, and if so calls comparison functions.
//Should only be called once per frame as it ticks up a frame counter to retrive the input for the current frame and facing
func cbrLoop() (input int32, facing bool) {
	//if no cases exist in the case base yet, dont replay anything

	aiData.rawFrames.midFightLearn(aiData.framedata)
	if aiData.cbrData == nil || aiData.cbrData.ReplayFile == nil || len(aiData.cbrData.ReplayFile) < 1 ||
		aiData.cbrData.ReplayFile[0].Case == nil || len(aiData.cbrData.ReplayFile[0].Case) < 1 {
		if midFightLearning.active == false {
			print("\n NO CASES FOUND")
		}
		return 0, false
	}

	rand.Seed(cbrProcess.randSeed)
	cbrProcess.lastCase.frameIndex++ //increments frame counter since case started
	var bufferCase *CBRData_Case
	var bufferCompVar float32
	var bufferReplayIndex int
	var bufferCaseIndex int
	//if no case in use or big statechange happens, find a new case.
	if cbrProcess.lastCase.cbrCase == nil || significantGameScopeChangeCheck(aiData.curGamestate, cbrProcess.lastCase) {
		debugUtil.addFrame(cbrProcess.framesSinceStart)
		debugUtil.setSelectionReason("significantScopeChange")

		cbrProcess.lastCase.frameIndex = 0 //resets frame counter as new case is selected
		cbrProcess.lastCase.cbrCase, cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.comparisonVal = findBestCase(aiData.curGamestate, aiData.cbrData, aiData.framedata)

		debugUtil.addCurGamestateFrame(aiData.curGamestate)
		debugUtil.addChosenCase(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.comparisonVal)
		cbrProcess.caseUsageIncrease(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.framesSinceStart)
		if debugUtil.debugActive == true {
			debugUtil.prepareDebugData(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1],
				aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame[len(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame)-1],
				aiData.cbrData, len(debugUtil.comparisons)-1)
		}
	} else {
		//if case is over evaluate how good the next case is, unless another case is significantly better, take the next case.
		if caseOver(cbrProcess.lastCase.cbrCase, cbrProcess.lastCase.frameIndex) {
			debugUtil.addFrame(cbrProcess.framesSinceStart)

			cbrProcess.lastCase.frameIndex = 0 //resets frame counter as new case is selected
			nextCase, noNextCase := getNextCase(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex)
			//go through all constraints and if any constraints are not fulfilled, skip the next case
			if noNextCase == false {
				for k := range nextCase.ExecutionConditions {
					if aiData.curGamestate.constraintCheck(CBRinter.commandExecutionConditions[nextCase.ExecutionConditions[k]], nextCase.ExecutionConditions[k]) == false {
						noNextCase = true
					}
				}
			}
			//if the next case can be used
			if noNextCase == false {
				cbrProcess.lastCase.cbrCase = nextCase
				cbrProcess.lastCase.caseIndex++
				refIds := refIdsMapping(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1], aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame[len(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame)-1], aiData.cbrData, cbrProcess.lastCase.replayIndex)
				cbrProcess.lastCase.comparisonVal = comparisonFunction(aiData.curGamestate, aiData.cbrData, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.caseIndex, aiData.framedata, &refIds)
				debugUtil.resetFrame(cbrProcess.framesSinceStart)

				bufferCase, bufferCaseIndex, bufferReplayIndex, bufferCompVar = findBestCase(aiData.curGamestate, aiData.cbrData, aiData.framedata)
				debugUtil.setSelectionReason("caseOver - NextCaseChosen")
				//if the found case is significatly better than the next case
				debugUtil.addNextCase(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex)
				if significantlyBetterCaseCheck(cbrProcess.lastCase.comparisonVal, bufferCompVar) {
					cbrProcess.lastCase.cbrCase = bufferCase
					cbrProcess.lastCase.caseIndex = bufferCaseIndex
					cbrProcess.lastCase.comparisonVal = bufferCompVar
					cbrProcess.lastCase.replayIndex = bufferReplayIndex
					debugUtil.setSelectionReason("caseOver - BetterFound")
				}
				debugUtil.addCurGamestateFrame(aiData.curGamestate)
				debugUtil.addChosenCase(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.comparisonVal)
				if debugUtil.debugActive == true {
					debugUtil.prepareDebugData(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1],
						aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame[len(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame)-1],
						aiData.cbrData, len(debugUtil.comparisons)-1)
				}
			} else { //if replay ended find a new case
				debugUtil.setSelectionReason("caseOver - NoNextCase")
				cbrProcess.lastCase.cbrCase, cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.comparisonVal = findBestCase(aiData.curGamestate, aiData.cbrData, aiData.framedata)
				debugUtil.addCurGamestateFrame(aiData.curGamestate)
				debugUtil.addChosenCase(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.lastCase.comparisonVal)
				if debugUtil.debugActive == true {
					debugUtil.prepareDebugData(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1],
						aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame[len(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame)-1],
						aiData.cbrData, len(debugUtil.comparisons)-1)
				}
			}
			cbrProcess.caseUsageIncrease(cbrProcess.lastCase.caseIndex, cbrProcess.lastCase.replayIndex, cbrProcess.framesSinceStart)
		}
	}

	//return the input and facing of the current frame from the current case
	curFrameIndex := cbrProcess.lastCase.cbrCase.FrameStartId + cbrProcess.lastCase.frameIndex
	curInput := aiData.cbrData.ReplayFile[cbrProcess.lastCase.replayIndex].Frame[curFrameIndex].Input
	curFacing := aiData.cbrData.ReplayFile[cbrProcess.lastCase.replayIndex].Frame[curFrameIndex].Facing

	cbrProcess.framesSinceStart++
	return curInput, curFacing
}

func caseOver(curCase *CBRData_Case, frameIndex int32) bool {
	return curCase.FrameStartId+frameIndex > curCase.FrameEndId
}

//triggers if some major change happens while a case is playing out
func significantGameScopeChangeCheck(curGamestate *CBRRawFrames, lastCase *CBRCaseReference) bool {
	sigChange := false
	sigChange = wasHit(curGamestate)
	return sigChange
}
func wasHit(curGamestate *CBRRawFrames) bool {
	curData := curGamestate.ReplayFile[len(curGamestate.ReplayFile)-1].Frame[len(curGamestate.ReplayFile[len(curGamestate.ReplayFile)-1].Frame)-1].CharData[curGamestate.ReplayFile[len(curGamestate.ReplayFile)-1].CbrFocusCharNr]
	if curData.ComparisonData.SelfGuard || curData.ComparisonData.SelfHit {
		return true
	}
	return false
}

//triggers if a case ends and the next case is significantly worse than another case in the caseBase
func significantlyBetterCaseCheck(nextInReplayComparisonVal float32, bestFoundInCaseBaseCompVar float32) bool {
	return nextInReplayComparisonVal > bestFoundInCaseBaseCompVar+cbrParameters.betterCaseThreshold
}
func getNextCase(curCaseIndex int, curCaseReplayIndex int) (nextCase *CBRData_Case, noNextCase bool) {
	if curCaseIndex+1 >= len(aiData.cbrData.ReplayFile[curCaseReplayIndex].Case) {
		return nil, true
	} else {
		return aiData.cbrData.ReplayFile[curCaseReplayIndex].Case[curCaseIndex+1], false
	}
}

//used to restart the cbr process
func resetCBRProcess() bool {
	cbrProcess.curScope = -1
	cbrProcess.lastCase.cbrCase = nil
	cbrProcess.lastCase.caseIndex = -1
	cbrProcess.lastCase.replayIndex = -1
	cbrProcess.lastCase.comparisonVal = -1
	cbrProcess.lastCase.frameIndex = -1
	cbrProcess.caseUsageReplayFile = map[int]map[int]CaseUsageData{}
	cbrProcess.framesSinceStart = 0
	return true
}

//goes through the caseBase and find the best cases, then randomly choose among the best cases.
//cbrParameters.topSelectionThreshold determines how much worse than the best case other cases are allowed to be to be eligable for selection
func findBestCase(curGamestate *CBRRawFrames, data *CBRData, framedata *Framedata) (bestCase *CBRData_Case, caseIndex int, replayIndex int, comparisonValue float32) {
	var bestCaseRefs []CBRCaseReference
	var highestCompVal = float32(0)
	var lowestCompVal = float32(0)
	var bufferCompVal float32
	var bestCtrlVal float32 = -1
	var bestCtrlCaseIndex = -1
	var bestCtrlCaseReplayIndex = -1

	refIds := CharReferenceIDs{}
	for i := range data.ReplayFile {
		refIds = refIdsMapping(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1], aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame[len(aiData.curGamestate.ReplayFile[len(aiData.curGamestate.ReplayFile)-1].Frame)-1], data, i)
		for j, cbrCase := range data.ReplayFile[i].Case {

			bufferCompVal = comparisonFunction(aiData.curGamestate, data, i, j, framedata, &refIds)

			//We try to filter out multiple controllable cases in a row, since we split these cases into smaller pieces and
			//want to avoid them being over representated in the random selection.

			//if the case is controllable...
			if cbrCase.Controllable == true {
				//... and its the first controllable case in sequence...
				if bestCtrlCaseIndex == -1 {
					//...set the case values
					bestCtrlCaseIndex = j
					bestCtrlCaseReplayIndex = i
					bestCtrlVal = bufferCompVal
				} else if bufferCompVal <= bestCtrlVal {
					//... else if the current comparison value is better replace the prior case
					bestCtrlCaseIndex = j
					bestCtrlCaseReplayIndex = i
					bestCtrlVal = bufferCompVal
				}
			} else { // if we had controllable cases and are now entering a non controllable case, add the best case among the ctrl cases to bestCaseRefs
				if bestCtrlCaseIndex != -1 {
					highestCompVal, lowestCompVal, bestCaseRefs = bestCaseRefEvaluation(curGamestate, data, bestCtrlVal, highestCompVal, lowestCompVal, bestCaseRefs, bestCtrlCaseReplayIndex, bestCtrlCaseIndex)
					bestCtrlCaseIndex = -1
					bestCtrlCaseReplayIndex = -1
					bestCtrlVal = -1
				}
				//add the non controllable case to best Case Refs
				highestCompVal, lowestCompVal, bestCaseRefs = bestCaseRefEvaluation(curGamestate, data, bufferCompVal, highestCompVal, lowestCompVal, bestCaseRefs, i, j)
			}

		}
	}
	if bestCtrlCaseIndex != -1 {
		highestCompVal, lowestCompVal, bestCaseRefs = bestCaseRefEvaluation(curGamestate, data, bestCtrlVal, highestCompVal, lowestCompVal, bestCaseRefs, bestCtrlCaseReplayIndex, bestCtrlCaseIndex)
		bestCtrlCaseIndex = -1
		bestCtrlCaseReplayIndex = -1
		bestCtrlVal = -1
	}

	randIndex := 0
	if len(bestCaseRefs) > 1 {
		randIndex = rand.Intn(len(bestCaseRefs) - 1)
	}

	debugUtil.debugGetFrameInfo(curGamestate, refIds)
	return bestCaseRefs[randIndex].cbrCase, bestCaseRefs[randIndex].caseIndex, bestCaseRefs[randIndex].replayIndex, bestCaseRefs[randIndex].comparisonVal
}

func bestCaseRefEvaluation(curGamestate *CBRRawFrames, data *CBRData, bufferCompVal float32, highestCompVal float32, lowestCompVal float32, bestCaseRefs []CBRCaseReference, i int, j int) (float32, float32, []CBRCaseReference) {
	// if the current case is far more similar than all prior cases empty bestCaseRefs and add it
	//the constraints of the case need to be fulfilled to be valid
	executionFail := false
	if bufferCompVal < lowestCompVal-cbrParameters.topSelectionThreshold || len(bestCaseRefs) < 1 {
		for k := range data.ReplayFile[i].Case[j].ExecutionConditions {
			if curGamestate.constraintCheck(CBRinter.commandExecutionConditions[data.ReplayFile[i].Case[j].ExecutionConditions[k]], data.ReplayFile[i].Case[j].ExecutionConditions[k]) == false {
				executionFail = true
			}
		}
		if executionFail == false {
			bestCaseRefs = nil
			bestCaseRefs = append(bestCaseRefs, CBRCaseReference{cbrCase: data.ReplayFile[i].Case[j], replayIndex: i, caseIndex: j, comparisonVal: bufferCompVal})
			highestCompVal = bufferCompVal
			lowestCompVal = bufferCompVal
		}

	} else {
		//if the current case is within range of the best case prior found and its constraints are fulfilled add it to bestCaseRefs
		if bufferCompVal-cbrParameters.topSelectionThreshold <= lowestCompVal {
			for k := range data.ReplayFile[i].Case[j].ExecutionConditions {
				if curGamestate.constraintCheck(CBRinter.commandExecutionConditions[data.ReplayFile[i].Case[j].ExecutionConditions[k]], data.ReplayFile[i].Case[j].ExecutionConditions[k]) == false {
					executionFail = true
				}
			}
			if executionFail == false {
				bestCaseRefs = append(bestCaseRefs, CBRCaseReference{cbrCase: data.ReplayFile[i].Case[j], replayIndex: i, caseIndex: j, comparisonVal: bufferCompVal})
				if bufferCompVal > highestCompVal {
					highestCompVal = bufferCompVal
				}
				if bufferCompVal < lowestCompVal {
					lowestCompVal = bufferCompVal
				}
			}

		}
	}
	return highestCompVal, lowestCompVal, bestCaseRefs
}

//checks weather the current gamestate fulfills the constraints given to the function
func (x *CBRRawFrames) constraintCheck(constraints []byte, constraintNr int32) bool {
	curFrame := x.ReplayFile[len(x.ReplayFile)-1].Frame[len(x.ReplayFile[len(x.ReplayFile)-1].Frame)-1].CharData[x.ReplayFile[len(x.ReplayFile)-1].CbrFocusCharNr]
	return constraintCheck(constraints, *curFrame, constraintNr)

}

func refIdsMapping(curGamestate *CBRRawFrames_ReplayFile, curFrame *CBRRawFrames_Frame, caseData *CBRData, replayIndex int) CharReferenceIDs {
	//finds which characters in the current gamestae maps to which character in the case for comparison.
	var refIds = CharReferenceIDs{curEnemyCharNr: -1, curFocusCharNr: -1, caseEnemyCharNr: -1, caseFocusCharNr: -1, focusID: "", enemyID: ""}
	refIds.curFocusCharNr = int(curGamestate.CbrFocusCharNr)
	refIds.curEnemyCharNr = -1
	for i := range curFrame.CharData {
		if i != refIds.curFocusCharNr {
			refIds.curEnemyCharNr = i
			break
		}
	}
	refIds.caseFocusCharNr = int(caseData.ReplayFile[replayIndex].CbrFocusCharNr)
	refIds.caseEnemyCharNr = -1
	for i := range caseData.ReplayFile[replayIndex].CharTeam {
		if i != refIds.caseFocusCharNr {
			refIds.caseEnemyCharNr = i
			break
		}
	}
	refIds.focusID = caseData.ReplayFile[replayIndex].CharName[refIds.caseFocusCharNr]
	refIds.enemyID = caseData.ReplayFile[replayIndex].CharName[refIds.caseEnemyCharNr]

	refIds.sameEnemy = refIds.enemyID == curGamestate.CharName[refIds.curEnemyCharNr]

	return refIds
}

func (x *CBRProcess) caseUsageIncrease(caseIndex int, replayIndex int, frame int64) {
	val, ok := x.caseUsageReplayFile[replayIndex]
	if !ok {
		x.caseUsageReplayFile[replayIndex] = map[int]CaseUsageData{}
		val, _ = x.caseUsageReplayFile[replayIndex]
	}
	val2, ok := val[caseIndex]
	if !ok {
		x.caseUsageReplayFile[replayIndex][caseIndex] = CaseUsageData{}
		val2 = x.caseUsageReplayFile[replayIndex][caseIndex]
	}
	val2.timesUsed++
	val2.lastFrameUsed = frame
	x.caseUsageReplayFile[replayIndex][caseIndex] = val2
}
