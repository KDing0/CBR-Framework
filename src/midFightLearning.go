package main

/*---FILE DESCRIPTION---
The interface file contains functions to let the game interface with the CBR system.
All data has to be requested or sent from the game, the CBR system never calls information from the game.
The interface functions need to be adjusted depending on the game the CBR system is used for.
The interface functions format the input data to be compatible with the CBR system and
format the output data of the CBR system to be compatible with the game system again.
---FILE DESCRIPTION---*/
type MidFightLearning struct {
	ctrlFrames           int
	active               bool
	framesSinceCaseGen   int
	lastReplay           CBRData_ReplayFile
	caseUnchangedCounter []int
	stableCaseIndex      int
	firstCaseAddition    bool
}

var midFightLearning = MidFightLearning{
	ctrlFrames:           0,
	active:               false,
	framesSinceCaseGen:   0,
	lastReplay:           CBRData_ReplayFile{},
	caseUnchangedCounter: []int{},
	stableCaseIndex:      -1,
	firstCaseAddition:    true,
}

func resetMidFightLearning() bool {
	midFightLearning.ctrlFrames = 0
	midFightLearning.framesSinceCaseGen = 0
	midFightLearning.lastReplay = CBRData_ReplayFile{}
	midFightLearning.caseUnchangedCounter = []int{}
	midFightLearning.stableCaseIndex = -1
	midFightLearning.firstCaseAddition = true

	return true
}

func (x CBRRawFrames) midFightLearn(framedata *Framedata) {
	if midFightLearning.active != true {
		midFightLearning.ctrlFrames = 0
		midFightLearning.framesSinceCaseGen = 0
		midFightLearning.lastReplay = CBRData_ReplayFile{}
		midFightLearning.caseUnchangedCounter = []int{}
		midFightLearning.stableCaseIndex = -1
		midFightLearning.firstCaseAddition = true

		if midFightLearning.active != true {
			return
		}
	}

	if x.ReplayFile != nil && len(x.ReplayFile) > 0 && x.ReplayFile[len(x.ReplayFile)-1].Frame != nil {
		curReplay := x.ReplayFile[len(x.ReplayFile)-1]
		curFrame := x.ReplayFile[len(x.ReplayFile)-1].Frame[len(x.ReplayFile[len(x.ReplayFile)-1].Frame)-1]
		if curFrame.CharData[curReplay.CbrFocusCharNr].ComparisonData.Controllable {
			midFightLearning.ctrlFrames++
		} else {
			midFightLearning.ctrlFrames = 0
		}
		midFightLearning.framesSinceCaseGen++

		if midFightLearning.framesSinceCaseGen > 12 && midFightLearning.ctrlFrames > 12 {
			var oldCbrData = []CBRData_Case{}
			if aiData.cbrData.ReplayFile != nil {
				for i := range aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case {
					oldCbrData = append(oldCbrData, *aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case[i])
				}
			}

			frameStartingIndex := 0
			if aiData.cbrData != nil && aiData.cbrData.ReplayFile != nil && len(aiData.cbrData.ReplayFile) > 0 && aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case != nil && len(aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case) > 0 && midFightLearning.stableCaseIndex != -1 {
				frameStartingIndex = int(aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case[midFightLearning.stableCaseIndex].FrameEndId + 1)
			}

			bufferRawFrames := *aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1]
			bufferRawFrames.Frame = bufferRawFrames.Frame[frameStartingIndex:]

			newReplay := bufferRawFrames.rawFramesToCBRReplay(framedata)
			if midFightLearning.firstCaseAddition == false {
				aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case = aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case[:midFightLearning.stableCaseIndex+1]
				aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Frame = aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Frame[:frameStartingIndex]
				aiData.cbrData.appendCBRReplay(newReplay)
			} else {
				aiData.cbrData.ReplayFile = append(aiData.cbrData.ReplayFile, newReplay)
				midFightLearning.firstCaseAddition = false
			}

			maxCounter := 5
			instableCasesNr := len(midFightLearning.caseUnchangedCounter) - midFightLearning.stableCaseIndex + 1
			maxCounter = maxCounter - int(maxInt32(int32(instableCasesNr/10), 0))
			for i := midFightLearning.stableCaseIndex + 1; i < len(aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case); i++ {
				if len(midFightLearning.caseUnchangedCounter) <= i {
					midFightLearning.caseUnchangedCounter = append(midFightLearning.caseUnchangedCounter, 0)
				}

				if oldCbrData != nil && i == midFightLearning.stableCaseIndex+1 && (aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1].Case[i].isSameCase(oldCbrData, i) || maxCounter <= 0) {

					if i == midFightLearning.stableCaseIndex+1 {
						midFightLearning.caseUnchangedCounter[i]++
						if midFightLearning.caseUnchangedCounter[i] >= maxCounter {
							midFightLearning.stableCaseIndex = i
							instableCasesNr = len(midFightLearning.caseUnchangedCounter) - midFightLearning.stableCaseIndex + 1
							maxCounter = maxCounter - int(maxInt32(int32(instableCasesNr/10), 0))
						}
					} else {
						midFightLearning.caseUnchangedCounter[i] = 0

					}
				} else {
					midFightLearning.caseUnchangedCounter[i] = 0
				}
			}

			midFightLearning.framesSinceCaseGen = 0

		}
	}
}

func (x *CBRData) appendCBRReplay(newReplay *CBRData_ReplayFile) *CBRData {
	var lastFrameIndex int32 = 0

	if x.ReplayFile != nil && len(x.ReplayFile) > 0 && x.ReplayFile[len(x.ReplayFile)-1].Case != nil && len(x.ReplayFile[len(x.ReplayFile)-1].Case) > 0 {
		lastFrameIndex = x.ReplayFile[len(x.ReplayFile)-1].Case[len(x.ReplayFile[len(x.ReplayFile)-1].Case)-1].FrameEndId + 1
	}
	for i := range newReplay.Case {
		newReplay.Case[i].FrameEndId += lastFrameIndex
		newReplay.Case[i].FrameStartId += lastFrameIndex
	}
	x.ReplayFile[len(x.ReplayFile)-1].Case = append(x.ReplayFile[len(x.ReplayFile)-1].Case, newReplay.Case...)
	x.ReplayFile[len(x.ReplayFile)-1].Frame = append(x.ReplayFile[len(x.ReplayFile)-1].Frame, newReplay.Frame...)

	return x
}

func (x *CBRData_Case) isSameCase(caseArr []CBRData_Case, index int) bool {

	if len(caseArr) > index && x.FrameStartId == caseArr[index].FrameStartId && x.FrameEndId == caseArr[index].FrameEndId {
		return true
	}
	return false
}
