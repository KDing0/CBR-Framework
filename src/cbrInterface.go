package main

import (
	"github.com/ikemen-engine/Ikemen-GO/src/cbr"
	"strconv"
)

/*---FILE DESCRIPTION---
Functions required to send information about the game-state to the CBR AI.
Functions for getting information from the CBR AI can be found in \cbr\interface.go
---FILE DESCRIPTION---*/

//Interface to send game-state data to the CBR AI every Frame
func (s *System) cbrAddFrame() bool {
	//saves player inputs for CBR AI
	if cbr.CheckFrameInsertable() || cbr.CheckCBRReplaying() {
		var iBit InputBits
		cbr.AddFrame()
		cbr.ReplayRecordStageData(s.cam.XMin, s.cam.XMax)
		cbr.ReplayRecordRoundState(s.chars[0][0].roundState())

		for i := range s.chars {
			if s.chars[i] != nil && s.chars[i][0] != nil {

				cbr.AddCharData()
				//inserts playable character information into the current frame of a replay
				if sys.netInput != nil {
					iBit.SetInput(0)
					//fmt.Printf("%v - %v%v \n", iBit, s.chars[i][0].name, i)
					//for o := range sys.inputRemap{
					//	fmt.Printf("Remap%v: %v \n",o, sys.inputRemap[o])
					//}
				} else {
					iBit.SetInput(sys.inputRemap[i])
					//fmt.Printf("%v - %v%v \n", iBit, s.chars[i][0].name, i)
					//for o := range sys.inputRemap{
					//	fmt.Printf("Remap%v: %v \n",o, sys.inputRemap[o])
					//}
				}

				//getCurrentMoveFrame(i) getCurrentMoveReference(i, s.chars[i][0])
				cbr.ReplayRecordCharState(i, int32(s.chars[i][0].ss.stateType), int32(s.chars[i][0].ss.moveType), s.chars[i][0].ctrl())
				cbr.ReplayRecordInputs(i, int32(iBit), s.chars[i][0].facing)
				cbr.ReplayRecordFramedata(i, s.chars[i][0].ss.time, int64(s.chars[i][0].ss.no))
				cbr.ReplayRecordMeters(i, float32(s.chars[i][0].life)/float32(s.chars[i][0].lifeMax), float32(s.chars[i][0].power)/float32(s.chars[i][0].powerMax), float32(s.chars[i][0].powerMax), float32(s.chars[i][0].dizzyPoints)/float32(s.chars[i][0].dizzyPointsMax), float32(s.chars[i][0].guardPoints)/float32(s.chars[i][0].guardPointsMax), float32(s.chars[i][0].redLife)/float32(s.chars[i][0].lifeMax))
				cbr.ReplayRecordVelocity(i, s.chars[i][0].vel[0], s.chars[i][0].vel[1])
				cbr.ReplayRecordAttackHit(i, s.chars[i][0].moveGuarded() == 1, s.chars[i][0].moveHit() == 1)
				cbr.ReplayRecordPosition(i, s.chars[i][0].pos[0], s.chars[i][0].pos[1])
				cbr.ReplayRecordIkemenSpecific(i, s.chars[i][0].ss.no, s.chars[i][0].ss.time)
				cbr.ReplayRecordFrameAdv(i, int32(checkFrameAdvantageState(i)))
				_, movesUsed, pressure := checkComboState(i)
				cbr.ReplayRecordComboInfo(i, pressure, movesUsed)

				if s.chars[i][0].ghv.guarded {
					cbr.ReplayRecordStun(i, s.chars[i][0].ghv.hittime, 0)
				} else {
					cbr.ReplayRecordStun(i, 0, s.chars[i][0].ghv.hittime)
				}

				//unspecified variables that players can use when making their characters for all kinds of purposes
				//since the variables can be used to stop certain attacks from coming out, we save them to check for conditions.
				genericInt := s.chars[i][0].ivar[:]
				genericFloat := s.chars[i][0].fvar[:]

				//adds the input buffer of the current character
				if s.chars[i][0].cmd[0].Buffer != nil {
					buf := s.chars[i][0].cmd[0].Buffer
					var dir []int32
					var button []int32

					dir = append(dir, buf.Ub, buf.Db, buf.Fb, buf.Bb)
					button = append(button, buf.ab, buf.bb, buf.cb, buf.db, buf.mb, buf.sb, buf.wb, buf.xb, buf.yb, buf.zb)
					cbr.ReplayRecordInputBuffer(i, dir, button)
				}
				cbr.ReplayRecordGenericVars(i, genericInt, genericFloat)

				//commandExecution added to frame recordings
				//records when a command is started to be input and when the command is executed
				//needed to determine how cases are generated
				if len(cbrCommandBufferState) <= i {
					cbrCommandBufferState = append(cbrCommandBufferState, Char{})
				}
				//runs complicated function to check if a command was executed, uses cbrCommandBufferState to save a prior character state
				//because this function is called after the input is already executed, which would make checks fail without the prior character state
				commandExecuted, commandIds, execIndex := commandCheck(i, cbrCommandBufferState[i], *s.chars[i][0])
				if commandExecuted == true {
					cbr.ReplayRecordCommands(i, commandIds, 1, execIndex)
				}
				cbrCommandBufferState[i] = *s.chars[i][0] //updates the prior character state

				//checks if a command buffer is in use.
				for j := range s.chars[i][0].cmd {
					for k := range s.chars[i][0].cmd[j].Commands {
						for l := range s.chars[i][0].cmd[j].Commands[k] {
							var stringIds []string
							stringIds = append(stringIds, s.chars[i][0].cmd[j].Commands[k][l].name+"_"+strconv.Itoa(l))
							if s.chars[i][0].cmd[j].Commands[k][l].cur > 0 {
								cbr.ReplayRecordCommands(i, stringIds, 0, -1)
							} else {
								cbr.ReplayRecordCommands(i, stringIds, -1, -1)
							}
						}
					}
				}
				//inserts the data of all helpers for the character into the frame of a file
				for j := range s.chars[i][1:] {
					if s.chars[i][j+1] != nil && s.chars[i][j+1].helperIndex >= 0 {

						helperNr := j + 1

						cbr.AddHelperData(i)
						cbr.HelperReplayRecordPosition(i, s.chars[i][helperNr].pos[0], s.chars[i][helperNr].pos[1], s.chars[i][helperNr].facing)
						cbr.HelperReplayRecordState(i, int32(s.chars[i][helperNr].ss.stateType), int32(s.chars[i][helperNr].ss.moveType), s.chars[i][helperNr].ctrl(), s.chars[i][helperNr].ghv.hittime)
						cbr.HelperReplayRecordFramedata(i, getHelperMoveFrame(i, j), getHelperMoveReference(i, j))
						cbr.HelperReplayRecordAttackHit(i, s.chars[i][helperNr].moveHit() == 1, s.chars[i][helperNr].moveGuarded() == 1)
						cbr.HelperReplayRecordMeters(i, float32(s.chars[i][helperNr].life)/float32(s.chars[i][helperNr].lifeMax))
						cbr.HelperReplayRecordVelocity(i, s.chars[i][helperNr].vel[0], s.chars[i][helperNr].vel[1])
						cbr.HelperReplayRecordGenericVars(i, s.chars[i][helperNr].helperId, s.chars[i][helperNr].ivar[:], s.chars[i][helperNr].fvar[:])

						if s.chars[i][helperNr].curFrame.Ex == nil {
							cbr.HelperReplayRecordCollisionBoxes(i, false, false)
						} else {
							hurtboxBool := s.chars[i][helperNr].curFrame.Ex[0] != nil
							hitboxBool := s.chars[i][helperNr].curFrame.Ex[1] != nil
							cbr.HelperReplayRecordCollisionBoxes(i, hurtboxBool, hitboxBool)
						}

					}
				}
			}
		}
	}
	return true
}

//called while fighting every frame to update the CBR ai recording/replaying state.
//used so that setting changes dont take effect every time a setting is changed, but instead collectively after settings were changed.
var aiUpdateDelay = 10

func updateCbrAiState(discard bool, skipDelay bool) {
	focusCharNr := cbr.GetReplayingCharIndex()
	recordingCharNr := cbr.GetRecordingCharIndex()
	if cbr.CheckAiActivityChange() && (sys.chars[0][0].alive() || sys.chars[1][0].alive()) {
		if !skipDelay && aiUpdateDelay > 0 {
			aiUpdateDelay--
			return
		}
		aiUpdateDelay = 10
		var charName []string
		var charTeam []int32
		for i := range sys.chars {
			if sys.chars[i] != nil && sys.chars[i][0] != nil {
				charTeam = append(charTeam, int32(sys.chars[i][0].teamside))
				charName = append(charName, sys.chars[i][0].name+"_"+sys.cgi[i].author)
			}
		}

		conditions := getTrimmedByteCode(focusCharNr, *sys.chars[focusCharNr][0])
		cbr.ResetCommandExecConditions()
		for i, arr := range conditions {
			for _, val := range arr {
				cbr.AddCommandExecConditions(i, byte(val))
			}
		}

		cbr.UpdateAiActivity(int32(recordingCharNr), charName, charTeam, getFramedata(), int32(focusCharNr), discard)
	}
}

//ends all cbrAI activity
func endCbrActivity(discard bool) {
	cbr.SetRecording(false)
	cbr.SetMidFightLearning(false)
	cbr.SetReplaying(false)
	updateCbrAiState(discard, true)
}

//--------more complex data read functions-------------
type FrameAdvantageState struct {
	initalStunState bool
	initalEnemyNr   int
	Frameadvantage  int
	selfInitator    bool

	returningFramedataTime int
}

var stunStates = map[int]FrameAdvantageState{}

func checkFrameAdvantageState(charIndex int) int {
	if sys.chars[charIndex] == nil || sys.chars[charIndex][0] == nil {
		return -1
	}

	val, ok := stunStates[charIndex]
	if !ok {
		val = FrameAdvantageState{initalStunState: false, initalEnemyNr: -1, returningFramedataTime: 0, Frameadvantage: 0}
		stunStates[charIndex] = val
	}

	hitStun := sys.chars[charIndex][0].ghv.hittime > 0
	ctrl := sys.chars[charIndex][0].ctrl()
	attack := sys.chars[charIndex][0].ss.moveType&MT_A > 0

	enemyHitStun := -1
	enemyCtrl := false
	enemyAttack := -1
	for i := range sys.chars {
		if i != charIndex && sys.chars[i] != nil {
			if sys.chars[i][0].ss.moveType&MT_A > 0 {
				enemyAttack = i
			}
			if sys.chars[i][0].ghv.hittime > 0 {
				enemyHitStun = i
			}
			if val.initalEnemyNr == i && sys.chars[i][0].ctrl() {
				enemyCtrl = true
			}
		}
	}

	if val.initalStunState == false {
		if hitStun == true && enemyAttack >= 0 {
			val.initalStunState = true
			val.initalEnemyNr = enemyAttack
			val.Frameadvantage = 0
			val.selfInitator = false

		}
		if attack == true && enemyHitStun >= 0 {
			val.initalStunState = true
			val.initalEnemyNr = enemyHitStun
			val.Frameadvantage = 0
			val.selfInitator = true
		}

	} else {
		if ctrl && !enemyCtrl {
			val.Frameadvantage = 1
			val.returningFramedataTime = 10
		}
		if !ctrl && enemyCtrl {
			val.Frameadvantage = -1
			val.returningFramedataTime = 10
		}

		if ctrl && enemyCtrl {
			val.initalStunState = false
		}
	}

	returnAdv := 0
	if val.returningFramedataTime > 0 {
		val.returningFramedataTime--
		if val.selfInitator {
			returnAdv = val.Frameadvantage * 2
		}

	}

	stunStates[charIndex] = val
	return returnAdv
}

type ComboState struct {
	inCombo       bool
	initalEnemyNr int

	pressure  bool
	curMoveID int64
	movesUsed int32

	didGatling bool
}

var comboStates = map[int]ComboState{}

func checkComboState(charIndex int) (combo bool, pressure bool, movesUsed int32) {
	if sys.chars[charIndex] == nil || sys.chars[charIndex][0] == nil {
		return false, false, 0
	}

	val, ok := comboStates[charIndex]
	if !ok {
		val = ComboState{inCombo: false, initalEnemyNr: -1, curMoveID: -999, movesUsed: 0}
		comboStates[charIndex] = val
	}
	hitStun := sys.chars[charIndex][0].ghv.hittime > 0
	ctrl := sys.chars[charIndex][0].ctrl()
	attack := sys.chars[charIndex][0].ss.moveType&MT_A > 0

	enemyHitStun := -1
	enemyCtrl := false

	for i := range sys.chars {
		if i != charIndex && sys.chars[i] != nil {
			if sys.chars[i][0].ghv.hittime > 0 {
				enemyHitStun = i
			}
			if val.initalEnemyNr == i && sys.chars[i][0].ctrl() {
				enemyCtrl = true
			}
		}
	}

	if val.inCombo == false {

		if attack == true && enemyHitStun >= 0 {
			val.inCombo = true
			val.initalEnemyNr = enemyHitStun
			if sys.chars[val.initalEnemyNr][0].ghv.guarded {
				val.pressure = true
			}
		}

	} else {
		if ctrl && enemyCtrl || hitStun {
			val.inCombo = false
			val.pressure = false
			val.movesUsed = 0
			val.curMoveID = -999
			val.didGatling = false
		}

		if val.inCombo && !ctrl {
			curMove := getCurrentMoveReference(charIndex, sys.chars[charIndex][0])

			if sys.chars[val.initalEnemyNr][0].ghv.guarded != val.pressure {
				val.curMoveID = curMove
				val.movesUsed = 1
				val.didGatling = false
				val.pressure = sys.chars[val.initalEnemyNr][0].ghv.guarded
			}

			if val.curMoveID != curMove || val.didGatling == true {

				val.curMoveID = curMove
				val.movesUsed++
				val.didGatling = false
				//fmt.Printf("%v\n", val.curMoveID)
				//getCurrentMoveReference(charIndex, sys.chars[charIndex][0])

			}
		}

	}

	comboStates[charIndex] = val
	return val.inCombo, val.pressure, val.movesUsed
}
func getComboState(charIndex int) (combo bool, pressure bool, movesUsed int32) {
	if sys.chars[charIndex] == nil || sys.chars[charIndex][0] == nil {
		return false, false, 0
	}

	val, ok := comboStates[charIndex]
	if !ok {
		val = ComboState{inCombo: false, initalEnemyNr: -1, curMoveID: -999, movesUsed: 0}
		comboStates[charIndex] = val
	}

	return val.inCombo, val.pressure, val.movesUsed
}
