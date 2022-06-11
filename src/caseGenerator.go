package main

/*---FILE DESCRIPTION---
Contains the process for processing a replay of annotated frames (rawFrames) and turn them into a replay of cases.
Starts a new case either when a state change occurs, including the frame that contains the input which causes the change.
Or starts a new case when the character is controllable and has been in the same state for long enough.
---FILE DESCRIPTION---*/

type Generator struct {
	maxCaseLength           int32
	minCaseLength           int32
	targetCaseLength        int32
	currentState            int64 //used to store the state ID the character is currently in, to know when states switch
	currentExecutedCommands []string
	constraintBuffer        []byte
}

var generator = Generator{
	maxCaseLength:           20, //should split cases if they grow too long
	minCaseLength:           8,  //If a case is shorter than this, it should be merged with a another case to be long enough.
	targetCaseLength:        12, //determines how long cases are when a player can control the character during the case.
	currentExecutedCommands: []string{},
	constraintBuffer:        []byte{},
}

// used to mark a weather a command was executed, started or both in a frame
const (
	cmdEx_Start            = 0
	cmdEx_Executed         = 1
	cmdEx_ExecutedAndStart = 2
)

//main frames to case conversion function
//Takes a replay of rawFrames and converts it into a replay containing cases.
func (x *CBRRawFrames_ReplayFile) rawFramesToCBRReplay(framedata *Framedata) *CBRData_ReplayFile {
	generator.currentState = -1
	generator.currentExecutedCommands = []string{}
	generator.constraintBuffer = []byte{}

	var inputs []int32
	for _, val := range x.Frame {
		inputs = append(inputs, val.CharData[0].Input)
	}

	cbrReplay := CBRData_ReplayFile{}
	if x != nil {
		//set up before looping through all frames
		cbrReplay.CbrFocusCharNr = x.CbrFocusCharNr //saves which character in the replay is the one the CBR ai is supposed to mimic
		cbrReplay.CharTeam = x.CharTeam
		cbrReplay.CharName = x.CharName
		var workingCase *CBRData_Case
		workingCase = nil
		for i := len(x.Frame) - 1; i >= 0; i-- {

			//prepend to the input of the frames in the CBRData. Stores only input relevant data.
			frameBuff := &CBRData_Frame{Input: x.Frame[i].CharData[x.CbrFocusCharNr].Input, Facing: x.Frame[i].CharData[x.CbrFocusCharNr].Facing}
			cbrReplay.Frame = append([]*CBRData_Frame{frameBuff}, cbrReplay.Frame...)

			//Creates a new case on this frame if no case generation in progress
			if workingCase == nil {
				workingCase = startCaseGeneration(x.Frame[i], x.CbrFocusCharNr, int32(i))
				generator.currentState = x.Frame[i].CharData[x.CbrFocusCharNr].ComparisonData.CurrentMoveReferenceID
			}

			workingCase.checkControllable(x.Frame[i], x.CbrFocusCharNr, int32(i))
			x.Frame[i].addCaseDebugInfo(x.Frame[i], int(x.CbrFocusCharNr))
			//if any command was executed we ignore every other check and try to find the start of the command input
			if len(generator.currentExecutedCommands) > 0 {
				//logic to check for the beginning of a executed command
				workingCase, cbrReplay = x.commandBeginningCheck(workingCase, cbrReplay, i)
			} else {
				//logic that happens when a state transition happens
				workingCase, cbrReplay = x.caseSplitCheck(workingCase, cbrReplay, i)
			}
		}
		//finishing the last case
		if workingCase != nil {
			workingCase = finishCaseGeneration(x.Frame[0], workingCase, x.CbrFocusCharNr, int32(0), x)
			cbrReplay.Case = append([]*CBRData_Case{workingCase}, cbrReplay.Case...)
			workingCase = nil
		}

		//goes through the case array again to split cases that are too big
		for i := 0; i < len(cbrReplay.Case); i++ {
			if cbrReplay.Case[i].Controllable && cbrReplay.Case[i].checkMaxCaseSize() {
				endIndex := cbrReplay.Case[i].FrameEndId
				splitIndex := cbrReplay.Case[i].FrameStartId + (cbrReplay.Case[i].FrameEndId-cbrReplay.Case[i].FrameStartId)/2
				cbrReplay.Case = cbrReplay.splitCase(i, x, x.CbrFocusCharNr, endIndex, splitIndex)
				i--
			} else if cbrReplay.Case[i].checkMaxLingeringCaseSize() {
				endIndex := cbrReplay.Case[i].FrameEndId
				splitIndex := cbrReplay.Case[i].ControllableLastFrame + (cbrReplay.Case[i].FrameEndId-cbrReplay.Case[i].ControllableLastFrame)/2
				cbrReplay.Case = cbrReplay.splitCase(i, x, x.CbrFocusCharNr, endIndex, splitIndex)
				i--
			}
		}
	}
	return &cbrReplay
}

//logic that happens to check weather a case should be split/finished creating
func (x *CBRRawFrames_ReplayFile) caseSplitCheck(workingCase *CBRData_Case, cbrReplay CBRData_ReplayFile, i int) (*CBRData_Case, CBRData_ReplayFile) {
	//if a state transition is happening and the min size was reached, or the next state is an attack, make a new case
	if checkStateTransition(x.Frame[i], x.CbrFocusCharNr) || checkGotHit(x.Frame[i], x.CbrFocusCharNr) || checkLandedHit(x.Frame[i], x.CbrFocusCharNr) {
		//end generating on the frame prior to the statechange and generate a case on the frame causing the statechange
		commandExec := workingCase.checkCommandExecution(x.Frame[i], x.CbrFocusCharNr, int32(i))
		if !commandExec {
			if checkCaseLengthReverse(int32(i), *workingCase, generator.minCaseLength) {
				workingCase = finishCaseGeneration(x.Frame[i], workingCase, x.CbrFocusCharNr, int32(i), x)
				cbrReplay.Case = append([]*CBRData_Case{workingCase}, cbrReplay.Case...)
				workingCase = nil
			} else {
				generator.currentState = x.Frame[i].CharData[x.CbrFocusCharNr].ComparisonData.CurrentMoveReferenceID
			}
		} else {
			//if a command was executed that triggered a statechange, save the constraints that need to be fulfilled for that trigger to happen
			//and start searching for the beginning of the command execution
			constraintIndex := x.Frame[i].CharData[x.CbrFocusCharNr].ExecutionConditions
			constraints := CBRinter.commandExecutionConditions[x.Frame[i].CharData[x.CbrFocusCharNr].ExecutionConditions]
			constBuff := constraintCheck(constraints, *x.Frame[i].CharData[x.CbrFocusCharNr], constraintIndex)

			//if a command was executed and started on the same frame immediately add the cases constraint
			//else add the constraint to a buffer to use for later
			if generator.currentExecutedCommands == nil || len(generator.currentExecutedCommands) < 1 {
				if constBuff == true {
					addCaseExecutionConditions(constraintIndex, workingCase)
				}
				workingCase = finishCaseGeneration(x.Frame[i], workingCase, x.CbrFocusCharNr, int32(i), x)
				cbrReplay.Case = append([]*CBRData_Case{workingCase}, cbrReplay.Case...)
				workingCase = nil
			} else {
				if constBuff == true {
					addCaseExecutionConditions(constraintIndex, workingCase)
				}
			}
			generator.currentState = x.Frame[i].CharData[x.CbrFocusCharNr].ComparisonData.CurrentMoveReferenceID
		}
	}
	return workingCase, cbrReplay
}

//when a command was executed keeps running through prior frames until the start of the command input is found
func (x *CBRRawFrames_ReplayFile) commandBeginningCheck(workingCase *CBRData_Case, cbrReplay CBRData_ReplayFile, i int) (*CBRData_Case, CBRData_ReplayFile) {

	//if another command was executed during our search, add this to the commands we are looking for
	if workingCase.checkCommandExecution(x.Frame[i], x.CbrFocusCharNr, int32(i)) {
		//check if the constraints for the command are fulfilled and if they are ...
		constraintIndex := x.Frame[i].CharData[x.CbrFocusCharNr].ExecutionConditions
		constraints := CBRinter.commandExecutionConditions[x.Frame[i].CharData[x.CbrFocusCharNr].ExecutionConditions]
		constBuff := constraintCheck(constraints, *x.Frame[i].CharData[x.CbrFocusCharNr], constraintIndex)
		if constBuff == true && constraints != nil && len(constraints) > 0 {
			//... add the new constraints
			addCaseExecutionConditions(constraintIndex, workingCase)
		}
	}
	//if the beginning of a executed command was found, remove it from the list against which we check
	for j := len(generator.currentExecutedCommands) - 1; j >= 0; j-- {
		k := generator.currentExecutedCommands[j]
		_, ok := x.Frame[i].CharData[x.CbrFocusCharNr].CommandExecution[k]
		if ok {
			generator.currentExecutedCommands = unorderedRemove(generator.currentExecutedCommands, j)
		}
	}
	//if all executed commands are resolved i.e the start of their execution was found
	if len(generator.currentExecutedCommands) == 0 {

		//check if the constraints were already resolved when the command execution was started
		//if its not already resolved dont save the constraint, as we are expecting the command to resolve the constraint.
		for k := len(workingCase.ExecutionConditions) - 1; k >= 0; k-- {
			constraintIndex := workingCase.ExecutionConditions[k]
			constraints := CBRinter.commandExecutionConditions[constraintIndex]
			constBuff := constraintCheck(constraints, *x.Frame[i].CharData[x.CbrFocusCharNr], constraintIndex)
			if constBuff == false {
				removeCaseExecutionConditions(k, workingCase)
			}
		}

		//finish generating this case, and add the case to the front of the replay
		workingCase = finishCaseGeneration(x.Frame[i], workingCase, x.CbrFocusCharNr, int32(i), x)
		cbrReplay.Case = append([]*CBRData_Case{workingCase}, cbrReplay.Case...)
		workingCase = nil
	}
	return workingCase, cbrReplay
}

func unorderedRemove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

//These functions need to be edited when adding more comparison parameters as they determine what information cases contain -------------------
func startCaseGeneration(frame *CBRRawFrames_Frame, CbrFocusCharNr int32, frameIndex int32) *CBRData_Case {

	cbrCase := CBRData_Case{}
	cbrCase.FrameEndId = frameIndex
	cbrCase.Controllable = true
	cbrCase.CommandExecuteFrame = -1
	cbrCase.ControllableLastFrame = -1
	return &cbrCase
}

func finishCaseGeneration(frame *CBRRawFrames_Frame, cbrCase *CBRData_Case, CbrFocusCharNr int32, frameIndex int32, rawFrameArr *CBRRawFrames_ReplayFile) *CBRData_Case {
	cbrCase.FrameStartId = frameIndex
	cbrCase.WorldCBRComparisonData = frame.WorldCBRComparisonData

	for i := range frame.CharData {
		if len(cbrCase.CharCBRComparisonData) <= i {
			cbrCase.CharCBRComparisonData = append(cbrCase.CharCBRComparisonData, &CharCBRComparisonData{})
		}
		cbrCase.CharCBRComparisonData[i].ComparisonData = frame.CharData[i].ComparisonData

		//Adding generic Vars that are Ikemen specific
		for _, val := range ikemenVarImport[i] {
			if val.HelperID >= 0 {
				continue
			}

			if val.Float == true {
				if cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars == nil {
					cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars = map[int32]float32{}
				}
				cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars[val.VarNr] = frame.CharData[i].GenericFloatVars[val.VarNr]
			} else {
				if cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars == nil {
					cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars = map[int32]int32{}
				}
				cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars[val.VarNr] = frame.CharData[i].GenericIntVars[val.VarNr]
			}
		}

		for j := range frame.CharData[i].HelperData {
			cbrCase.CharCBRComparisonData[i].HelperData = append(cbrCase.CharCBRComparisonData[i].HelperData, &HelperComparisonData{})
			cbrCase.CharCBRComparisonData[i].HelperData[j] = frame.CharData[i].HelperData[j].CompData

			//Adding generic Vars that are Ikemen specific from helpers
			for _, val := range ikemenVarImport[i] {
				if val.HelperID < 0 || val.HelperID != frame.CharData[i].HelperData[j].CompData.HelperID || val.VarNr < 0 {
					continue
				}
				if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars == nil {
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars = &GenericVars{}
				}
				if val.Float == true {
					if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars == nil {
						cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars = map[int32]float32{}
					}
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars[val.VarNr] = frame.CharData[i].HelperData[j].GenericFloatVars[val.VarNr]
				} else {
					if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars == nil {
						cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars = map[int32]int32{}
					}
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars[val.VarNr] = frame.CharData[i].HelperData[j].GenericIntVars[val.VarNr]
				}
			}

		}
	}

	cbrCase.addCaseDebugInfo(rawFrameArr, int(CbrFocusCharNr), int(cbrCase.FrameStartId), int(cbrCase.FrameEndId))

	return cbrCase
}

// splits a case that contains only movement into more cases
func (x *CBRData_ReplayFile) splitCase(caseIndex int, rawFrameArr *CBRRawFrames_ReplayFile, CbrFocusCharNr int32, endIndex int32, splitIndex int32) []*CBRData_Case {
	x.Case[caseIndex].FrameEndId = splitIndex - 1
	if x.Case[caseIndex].ControllableLastFrame >= splitIndex {
		x.Case[caseIndex].ControllableLastFrame = -1
	}

	frame := rawFrameArr.Frame[splitIndex] //the frame at the beginning of the split case

	cbrCase := CBRData_Case{}
	cbrCase.FrameStartId = splitIndex
	cbrCase.FrameEndId = endIndex
	cbrCase.WorldCBRComparisonData = frame.WorldCBRComparisonData

	cbrCase.ControllableLastFrame = splitIndex
	cbrCase.Controllable = true
	for i := range frame.CharData {

		if len(cbrCase.CharCBRComparisonData) <= i {
			cbrCase.CharCBRComparisonData = append(cbrCase.CharCBRComparisonData, &CharCBRComparisonData{})
		}
		cbrCase.CharCBRComparisonData[i].ComparisonData = frame.CharData[i].ComparisonData

		//Adding generic Vars that are Ikemen specific
		for _, val := range ikemenVarImport[i] {
			if val.HelperID >= 0 {
				continue
			}

			if val.Float == true {
				if cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars == nil {
					cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars = map[int32]float32{}
				}
				cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericFloatVars[val.VarNr] = frame.CharData[i].GenericFloatVars[val.VarNr]
			} else {
				if cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars == nil {
					cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars = map[int32]int32{}
				}
				cbrCase.CharCBRComparisonData[i].ComparisonData.GenericVars.GenericIntVars[val.VarNr] = frame.CharData[i].GenericIntVars[val.VarNr]
			}
		}

		for j := range frame.CharData[i].HelperData {
			cbrCase.CharCBRComparisonData[i].HelperData = append(cbrCase.CharCBRComparisonData[i].HelperData, &HelperComparisonData{})
			cbrCase.CharCBRComparisonData[i].HelperData[j] = frame.CharData[i].HelperData[j].CompData

			//Adding generic Vars that are Ikemen specific from helpers
			for _, val := range ikemenVarImport[i] {
				if val.HelperID < 0 || val.HelperID != frame.CharData[i].HelperData[j].CompData.HelperID || val.VarNr < 0 {
					continue
				}
				if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars == nil {
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars = &GenericVars{}
				}
				if val.Float == true {
					if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars == nil {
						cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars = map[int32]float32{}
					}
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericFloatVars[val.VarNr] = frame.CharData[i].HelperData[j].GenericFloatVars[val.VarNr]
				} else {
					if cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars == nil {
						cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars = map[int32]int32{}
					}
					cbrCase.CharCBRComparisonData[i].HelperData[j].GenericVars.GenericIntVars[val.VarNr] = frame.CharData[i].HelperData[j].GenericIntVars[val.VarNr]
				}
			}
		}

	}

	cbrCase.addCaseDebugInfo(rawFrameArr, int(CbrFocusCharNr), int(cbrCase.FrameStartId), int(cbrCase.FrameEndId))

	return insert(x.Case, &cbrCase, caseIndex+1)
}

//----------------------------------------------------------------------------------------------

func addCaseExecutionConditions(conditionIndex int32, cbrCase *CBRData_Case) *CBRData_Case {
	cbrCase.ExecutionConditions = append(cbrCase.ExecutionConditions, conditionIndex)
	return cbrCase
}
func removeCaseExecutionConditions(conditionIndex int, cbrCase *CBRData_Case) *CBRData_Case {
	cbrCase.ExecutionConditions[conditionIndex] = cbrCase.ExecutionConditions[len(cbrCase.ExecutionConditions)-1] // Copy last element to index i.
	cbrCase.ExecutionConditions[len(cbrCase.ExecutionConditions)-1] = -1                                          // Erase last element (write zero value).
	cbrCase.ExecutionConditions = cbrCase.ExecutionConditions[:len(cbrCase.ExecutionConditions)-1]                // Truncate slice.
	return cbrCase
}

func (x CBRData_Case) checkCommandExecution(frame *CBRRawFrames_Frame, CbrFocusCharNr int32, frameIndex int32) bool {
	for k, v := range frame.CharData[CbrFocusCharNr].CommandExecution {
		if v == cmdEx_Executed {
			generator.currentExecutedCommands = append(generator.currentExecutedCommands, k)
			if x.CommandExecuteFrame < frameIndex {
				x.CommandExecuteFrame = frameIndex
			}
			return true
		}
		if v == cmdEx_ExecutedAndStart {
			return true
		}

	}
	return false
}

func checkCaseLengthReverse(frameIndex int32, workingCase CBRData_Case, length int32) bool {
	return workingCase.FrameEndId-frameIndex > length
}

//check if the character is changing states during this frame
func checkStateTransition(frame *CBRRawFrames_Frame, CbrFocusCharNr int32) bool {

	return generator.currentState != frame.CharData[CbrFocusCharNr].ComparisonData.CurrentMoveReferenceID
}

//check if the character got hit or blocked a attack
func checkGotHit(frame *CBRRawFrames_Frame, CbrFocusCharNr int32) bool {

	return frame.CharData[CbrFocusCharNr].ComparisonData.SelfHit || frame.CharData[CbrFocusCharNr].ComparisonData.SelfGuard
}

//check if the character got hit or blocked a attack
func checkLandedHit(frame *CBRRawFrames_Frame, CbrFocusCharNr int32) bool {

	return frame.CharData[CbrFocusCharNr].ComparisonData.MoveGuarded || frame.CharData[CbrFocusCharNr].ComparisonData.MoveHit
}

//check if the case reached a certain size
func (x CBRData_Case) checkMaxCaseSize() bool {
	return x.FrameEndId-x.FrameStartId > generator.maxCaseLength
}
func (x CBRData_Case) checkMinCaseSize() bool {
	return x.FrameEndId+1-x.FrameStartId < generator.minCaseLength
}

//check if the case reached a certain size
func (x CBRData_Case) checkMaxLingeringCaseSize() bool {
	if x.ControllableLastFrame != -1 {
		return x.FrameEndId-x.ControllableLastFrame > generator.maxCaseLength
	}
	return false
}

//check if the character was controllable during the frame and if it was not the case is marked as a non controllable case.
func (x *CBRData_Case) checkControllable(frame *CBRRawFrames_Frame, CbrFocusCharNr int32, frameNr int32) bool {

	if frame.CharData[CbrFocusCharNr].ComparisonData.Controllable == false {
		x.Controllable = false
	}
	if x.Controllable == true && (x.ControllableLastFrame > frameNr || x.ControllableLastFrame == -1) {
		x.ControllableLastFrame = frameNr
	}
	return x.Controllable
}

func insert(s []*CBRData_Case, element *CBRData_Case, i int) []*CBRData_Case {
	s = append(s, nil /* use the zero value of the element type */)
	copy(s[i+1:], s[i:])
	s[i] = element
	return s
}
