package main

import (
	"math"
	"strings"
)

/*---FILE DESCRIPTION---
The interface file contains functions to let the game interface with the CBR system.
All data has to be requested or sent from the game, the CBR system never calls information from the game.
The interface functions need to be adjusted depending on the game the CBR system is used for.
The interface functions format the input data to be compatible with the CBR system and
format the output data of the CBR system to be compatible with the game system again.
---FILE DESCRIPTION---*/

type InterfaceVars struct {
	commandBufferSave          []map[string]int32
	commandExecutionConditions map[int32][]byte
	blockStun                  []int32
	hitStun                    []int32
	recording                  bool
	replaying                  bool
	midFight                   bool

	recordingCharIndex int
	replayingCharIndex int

	replayFrames int32

	ctrlState   []bool
	allCtrlTime int32

	playerNames map[int]string
}

var CBRinter = InterfaceVars{
	blockStun: []int32{},
	hitStun:   []int32{},
	recording: false,
	replaying: false,
	midFight:  false,

	recordingCharIndex: 0,
	replayingCharIndex: 1,

	replayFrames: 0,

	ctrlState:   []bool{},
	allCtrlTime: 0,

	playerNames: map[int]string{},
}

type InvulType int32

const (
	IT_ST_S InvulType = 1 << iota
	IT_ST_C
	IT_ST_A
	IT_ST_L
	IT_ST_N
	IT_ST_U

	IT_AT_NA
	IT_AT_NT
	IT_AT_NP
	IT_AT_SA
	IT_AT_ST
	IT_AT_SP
	IT_AT_HA
	IT_AT_HT
	IT_AT_HP

	IT_MT_I
	IT_MT_H
	IT_MT_A
	IT_MT_U

	IT_MT_MNS = IT_MT_I
	IT_MT_PLS = IT_MT_H

	IT_ST_MASK = 63
	IT_ST_D    = IT_ST_L
	IT_ST_F    = IT_ST_N
	IT_ST_P    = IT_ST_U
	IT_ST_SCA  = IT_ST_S | IT_ST_C | IT_ST_A

	IT_AT_AA  = IT_AT_NA | IT_AT_SA | IT_AT_HA
	IT_AT_AT  = IT_AT_NT | IT_AT_ST | IT_AT_HT
	IT_AT_AP  = IT_AT_NP | IT_AT_SP | IT_AT_HP
	IT_AT_ALL = IT_AT_AA | IT_AT_AT | IT_AT_AP
	IT_AT_AN  = IT_AT_NA | IT_AT_NT | IT_AT_NP
	IT_AT_AS  = IT_AT_SA | IT_AT_ST | IT_AT_SP
	IT_AT_AH  = IT_AT_HA | IT_AT_HT | IT_AT_HP
)

type MovementState int32

const (
	MT_Standing = 1 << iota
	MT_Crouching
	MT_Airborne
	MT_LyingDown
	MT_Special
	MT_Undefined
)

type AttackState int32

const (
	AS_Idle = 1 << (iota + 15)
	AS_Hit
	AS_Attack
	AS_Undefined
)

type CommandBuffer struct {
	Bb, Db, Fb, Ub                         int32
	ab, bb, cb, xb, yb, zb, sb, db, wb, mb int32
	B, D, F, U                             int8
	a, b, c, x, y, z, s, d, w, m           int8
}

func DeleteData() bool {
	return deleteAllData("save/cbrData/")
}
func DeleteCharData(cbrFocusCharNr int32, charName []string) bool {

	return deleteCharData("save/cbrData/", charName[cbrFocusCharNr]+"_"+getPlayerName(int(cbrFocusCharNr)))
}

func SetPlayerName(index int, name string) {
	CBRinter.playerNames[index] = name
}
func getPlayerName(index int) string {
	name := ""
	n, ok := CBRinter.playerNames[index]
	if ok {
		name = n
	}

	return name

}
func SetRecordingCharIndex(index int) {
	CBRinter.recordingCharIndex = index
}
func SetReplayingCharIndex(index int) {
	CBRinter.replayingCharIndex = index
}
func GetRecordingCharIndex() int {
	return CBRinter.recordingCharIndex
}
func GetReplayingCharIndex() int {
	return CBRinter.replayingCharIndex
}

func SetRecording(recording bool) {
	print("\nRecording ")
	print(recording)
	CBRinter.recording = recording
}
func GetRecording() bool {
	return CBRinter.recording
}
func SetReplaying(replaying bool) {
	print("\nReplaying ")
	print(replaying)
	CBRinter.replaying = replaying
}
func GetReplaying() bool {
	return CBRinter.replaying
}

func SetMidFightLearning(midFightLearning bool) {
	CBRinter.midFight = midFightLearning
}
func GetMidfight() bool {
	return CBRinter.midFight
}
func CheckAiActivityChange() bool {
	if CBRinter.recording != aiData.recording || CBRinter.replaying != aiData.replaying || (CBRinter.replaying == true &&
		CBRinter.recording == true && CBRinter.midFight != midFightLearning.active) || (midFightLearning.active == true && CBRinter.midFight == false) {
		return true
	}
	return false
}

func UpdateAiActivity(recordingFocusCharNr int32, charName []string, charTeam []int32, framedata *Framedata, replayingFocusCharNr int32, discard bool) bool {
	if CBRinter.recording == true && CBRinter.replaying == true && CBRinter.midFight == true && (aiData.recording == false || aiData.replaying == false || midFightLearning.active == false) {
		//start midFightrecording
		ToggleMidFightRecording(recordingFocusCharNr, charName, charTeam, framedata, replayingFocusCharNr, 1, discard)

	} else if (CBRinter.recording == false || CBRinter.replaying == false || CBRinter.midFight == false) && aiData.recording == true && aiData.replaying == true && midFightLearning.active == true {
		//end midfightrecording
		ToggleMidFightRecording(recordingFocusCharNr, charName, charTeam, framedata, replayingFocusCharNr, -1, discard)
	}

	if CBRinter.recording == true && aiData.recording != true {
		//start recording
		ToggleRecording(recordingFocusCharNr, charName, charTeam, framedata, 1, discard)
	}
	if CBRinter.recording == false && aiData.recording != false {
		//end recording
		ToggleRecording(recordingFocusCharNr, charName, charTeam, framedata, -1, discard)
	}
	if CBRinter.replaying == true && aiData.replaying != true {
		//start replaying
		ToggleCBRReplaying(replayingFocusCharNr, charName, charTeam, framedata, 1)
	}
	if CBRinter.replaying == false && aiData.replaying != false {
		//end replaying
		ToggleCBRReplaying(replayingFocusCharNr, charName, charTeam, framedata, -1)
	}

	return false
}
func EndAiActivity(framedata *Framedata, bufferSave bool) {
	if aiData.replaying {
		print("\nCBRReplay End")
		aiData.endCBRReplaying()
		resetDebug("save/cbrData/")
	}
	if aiData.recording {
		print("\nCBRRecord End")
		if bufferSave == true {
			aiData.bufferStoreCBRRecording(framedata, "save/cbrData/")
		} else {
			aiData.endCBRRecording(framedata, "save/cbrData/")
		}

	}
}
func SaveCBRBuffer() {
	aiData.saveCBRBuffer("save/cbrData/")
}

///Interface functions for the outside program-----------
func ToggleRecording(cbrFocusCharNr int32, charName []string, charTeam []int32, framedata *Framedata, startOrEnd int32, discard bool) bool {
	for i := range charName {
		charName[i] = strings.ToValidUTF8(charName[i], "")
	}

	if aiData.recording == true && startOrEnd <= 0 {
		print("\nCBRRecord End")
		if discard == true {
			aiData.discardCBRRecording()
		} else {
			aiData.endCBRRecording(framedata, "save/cbrData/")
		}
	} else if startOrEnd >= 0 {
		CBRinter.replayFrames = 0
		print("\nCBRRecord Start")
		aiData.startCBRRecording(cbrFocusCharNr, charName, charTeam)
	}
	return true
}

///Interface functions for the outside program-----------
func ToggleMidFightRecording(recordingFocusCharNr int32, charName []string, charTeam []int32, framedata *Framedata, replayingFocusCharNr int32, startOrEnd int32, discard bool) bool {
	for i := range charName {
		charName[i] = strings.ToValidUTF8(charName[i], "")
		if charName[0] != charName[i] {
			print("\nMidFightLearning only works when using the same characters.")
			return false
		}
	}
	resetMidFightLearning()
	if midFightLearning.active == false && startOrEnd >= 0 {
		if aiData.recording == true {
			print("\nCBRRecord Restart")
			if discard == true {
				aiData.discardCBRRecording()
			} else {
				aiData.endCBRRecording(framedata, "save/cbrData/")
			}
			resetDebug("save/cbrData/")
		}
		print("\nCBRRecord Start")
		CBRinter.replayFrames = 0
		aiData.startCBRRecording(recordingFocusCharNr, charName, charTeam)

		if aiData.replaying == true {
			print("\nCBRReplay Restart")
			aiData.endCBRReplaying()
			resetDebug("save/cbrData/")
		}
		print("\nCBRReplay Start")
		aiData.cbrData = loadCBRData("save/cbrData/", charName[replayingFocusCharNr]+"_"+getPlayerName(int(replayingFocusCharNr)))
		aiData.startCBRReplaying(replayingFocusCharNr, charName, charTeam, framedata)
		setAIControlledCharacter(int(replayingFocusCharNr))
		midFightLearning.active = true
	} else if startOrEnd <= 0 {
		if aiData.recording == true {
			print("\nCBRRecord End")
			if discard == true {
				aiData.discardCBRRecording()
			} else {
				aiData.endCBRRecording(framedata, "save/cbrData/")
			}
		}
		if aiData.replaying == true {
			print("\nCBRReplay End")
			aiData.endCBRReplaying()
			resetDebug("save/cbrData/")
		}
		midFightLearning.active = false
	}

	return true
}

func EndRecording(framedata *Framedata) bool {
	if aiData.recording == true {
		print("\nCBRRecord End")
		aiData.endCBRRecording(framedata, "save/cbrData/")
	}
	return true
}

func ToggleCBRReplaying(cbrFocusCharNr int32, charName []string, charTeam []int32, framedata *Framedata, startOrEnd int32) bool {
	for i := range charName {
		charName[i] = strings.ToValidUTF8(charName[i], "")
	}

	if aiData.replaying == true && startOrEnd <= 0 {
		print("\nCBRReplay End")
		aiData.endCBRReplaying()
		resetDebug("save/cbrData/")
	} else if startOrEnd >= 0 {
		print("\nCBRReplay Start")
		aiData.cbrData = loadCBRData("save/cbrData/", charName[cbrFocusCharNr]+"_"+getPlayerName(int(cbrFocusCharNr)))
		aiData.startCBRReplaying(cbrFocusCharNr, charName, charTeam, framedata)
		setAIControlledCharacter(int(cbrFocusCharNr))
	}
	return aiData.replaying
}
func EndCBRReplaying() bool {
	if aiData.replaying == true {
		print("\nCBRReplay End")
		aiData.endCBRReplaying()
		resetDebug("save/cbrData/")
	}
	return aiData.replaying
}

func CheckCBRReplaying() bool {
	return aiData.replaying
}
func CheckCBRRecording() bool {
	return aiData.recording
}

/*
func CheckRawFrameReplaying() bool{
	if CheckRawFrame() == true {
		index := aiData.replayIndex
		ret := len(aiData.replayFrames) > index && aiData.replaying
		if ret == false {aiData.replaying = false}
		return ret
	}else{
		return false
	}
}
func ToggleRawFramesReplaying(playerNr int) bool {
	if aiData.replaying == true {
		aiData.replaying = false
		aiData.rawFrameReplay = false
	}else{
		replayFrames := RawFramesToReplay(*aiData.rawFrames, playerNr, 0)
		*aiData.rawFrames = CBRRawFrames{}
		aiData.InitializeReplaying(replayFrames)
		aiData.replaying = true
		aiData.rawFrameReplay = true
	}
	return aiData.replaying
}
func ReadRawFrameInput(facing int32) int32 {
	index := aiData.replayIndex
	input := aiData.replayFrames[index].Input
	storedFacing := aiData.replayFrames[index].Facing
	bFacing := facingToBool(facing)
	if bFacing != storedFacing{
		input = swapBitsAtPos(input, 2, 3)
	}
	return input
}
func CheckRawFrame() bool{
	return aiData.rawFrameReplay
}*/

func ReadCbrFrameInput(facing int32) (input int32) {
	input = readCbrFrameInput(facing)
	return
}

func GetAIControlledCharacter() (charIndex int) {
	return getAIControlledCharacter()
}

/*
func ReadCbrFrameInput(Raw) (int32, []InputBuffer)  {
	index := aiData.replayIndex
	input := aiData.replayFrames[index].Input
	storedFacing := aiData.replayFrames[index].Facing
	bFacing := facingToBool(facing)
	if bFacing != storedFacing{
		input = swapBitsAtPos(input, 2, 3)
	}
	return input
}

func IncrementReplayIndex() bool {
	ret := false
	if aiData.replaying {
		aiData.replayIndex++
		ret = true
	}
	return ret
}*/

func ResetCommandExecConditions() bool {
	CBRinter.commandExecutionConditions = map[int32][]byte{}
	return true
}
func AddCommandExecConditions(index int32, b byte) bool {
	CBRinter.commandExecutionConditions[index] = append(CBRinter.commandExecutionConditions[index], b)
	return true
}

func AddFrame() bool {
	if aiData.recording == true {
		CBRinter.replayFrames++
		aiData.rawFrames.AddFrame()
	}
	if CheckCBRReplaying() && len(aiData.curGamestate.ReplayFile) > 0 {
		aiData.curGamestate.queueFrame(cbrParameters.curGamestateQueLength)
	}
	return true
}

func AddCharData() bool {
	if aiData.recording == true {
		aiData.rawFrames.AddCharData(1)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.AddCharData(1)
	}
	return true
}

func CheckFrameInsertable() bool {
	return aiData.recording == true && len(aiData.rawFrames.ReplayFile) > 0
}

func AddHelperData(charNr int) bool {
	if aiData.recording == true {
		aiData.rawFrames.addHelperData(charNr)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.addHelperData(charNr)
	}
	return true
}

//---Adding char data for the replay
func ReplayRecordCharData(cbrFocusCharNr int32, charName string, charTeam int32) bool {
	charName = strings.ToValidUTF8(charName, "")

	if aiData.recording == true {
		aiData.rawFrames.setCharData(cbrFocusCharNr, charName, charTeam)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setCharData(cbrFocusCharNr, charName, charTeam)
	}
	return true
}

func ReplayRecordStageData(leftWallPos float32, rightWallPos float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setStageData(leftWallPos, rightWallPos)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setStageData(leftWallPos, rightWallPos)
	}
	return true
}
func ReplayRecordRoundState(roundState int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setRoundState(roundState)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setRoundState(roundState)
	}
	return true
}

//---Adding Player relevant data into a frame in a replay
func ReplayRecordInputs(playerNr int, inputs int32, facing float32) bool {
	bFacing := floatFacingToBool(facing)
	if aiData.recording == true {
		aiData.rawFrames.setPlayerInput(playerNr, inputs, bFacing)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setPlayerInput(playerNr, inputs, bFacing)
	}
	return true
}

func ReplayRecordCharState(playerNr int, movementState int32, attackState int32, controllable bool) bool {
	if aiData.recording == true {
		aiData.rawFrames.setCharState(playerNr, (movementState&MT_Standing) > 0, (movementState&MT_Crouching) > 0, (movementState&MT_Airborne) > 0, (movementState&MT_LyingDown) > 0, (attackState&AS_Idle) > 0, (attackState&AS_Hit) > 0, (attackState&AS_Attack) > 0, controllable)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setCharState(playerNr, (movementState&MT_Standing) > 0, (movementState&MT_Crouching) > 0, (movementState&MT_Airborne) > 0, (movementState&MT_LyingDown) > 0, (attackState&AS_Idle) > 0, (attackState&AS_Hit) > 0, (attackState&AS_Attack) > 0, controllable)
	}
	return true
}

func ReplayRecordFramedata(playerNr int, currentMoveFrame int32, currentMoveReferenceID int64) bool {
	if aiData.recording == true {
		aiData.rawFrames.setFramedata(playerNr, currentMoveFrame, currentMoveReferenceID)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setFramedata(playerNr, currentMoveFrame, currentMoveReferenceID)
	}
	return true
}
func ReplayRecordMeters(playerNr int, lifePercentage float32, meterPercentage float32, meterMax float32, dizzyPercentage float32, guardPointsPercentage float32, recoverableHpPercentage float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setMeters(playerNr, lifePercentage, meterPercentage, meterMax, dizzyPercentage, guardPointsPercentage, recoverableHpPercentage)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setMeters(playerNr, lifePercentage, meterPercentage, meterMax, dizzyPercentage, guardPointsPercentage, recoverableHpPercentage)
	}
	return true
}

func ReplayRecordVelocity(playerNr int, horizontalVelocity float32, verticalVelocity float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setVelocity(playerNr, horizontalVelocity, verticalVelocity)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setVelocity(playerNr, horizontalVelocity, verticalVelocity)
	}
	return true
}

func ReplayRecordStun(playerNr int, blockstun int32, hitStun int32) bool {
	selfHit := false
	selfGuard := false
	for len(CBRinter.blockStun) <= playerNr {
		CBRinter.blockStun = append(CBRinter.blockStun, 0)
	}
	for len(CBRinter.hitStun) <= playerNr {
		CBRinter.hitStun = append(CBRinter.hitStun, 0)
	}

	if hitStun > 0 && Abs(CBRinter.hitStun[playerNr]-hitStun) > 1 {
		selfHit = true
	}
	if blockstun > 0 && Abs(CBRinter.blockStun[playerNr]-blockstun) > 1 {
		selfGuard = true
	}
	CBRinter.hitStun[playerNr] = hitStun
	CBRinter.blockStun[playerNr] = blockstun

	if aiData.recording == true {
		aiData.rawFrames.setStun(playerNr, blockstun, hitStun)
		aiData.rawFrames.setSelfHit(playerNr, selfGuard, selfHit)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setStun(playerNr, blockstun, hitStun)
		aiData.curGamestate.setSelfHit(playerNr, selfGuard, selfHit)
	}
	return true
}

func ReplayRecordAttackHit(playerNr int, moveHit bool, moveGuarded bool) bool {

	if aiData.recording == true {
		aiData.rawFrames.setAttackHit(playerNr, moveHit, moveGuarded)

	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setAttackHit(playerNr, moveHit, moveGuarded)
	}
	return true
}

func ReplayRecordPosition(playerNr int, positionX float32, positionY float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setPosition(playerNr, positionX, positionY)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setPosition(playerNr, positionX, positionY)
	}
	return true
}

func ReplayRecordInputBuffer(playerNr int, directionBuffer []int32, buttonBuffer []int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setInputBuffer(playerNr, directionBuffer, buttonBuffer)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setInputBuffer(playerNr, directionBuffer, buttonBuffer)
	}
	return true
}

//Saves in the replay when a player begins inputting a command "0" and a command gets executed "1"
//commandBuffer is used to determine when a a command started to be input, if a -1 is saved in there the buffer is currently empty
func ReplayRecordCommands(playerNr int, commandIds []string, commandState int32, execId int32) bool {

	for _, commandId := range commandIds {

		commandId = strings.ToValidUTF8(commandId, "")

		if aiData.recording == true {
			if len(CBRinter.commandBufferSave) <= playerNr {
				mapBuffer := make(map[string]int32)
				CBRinter.commandBufferSave = append(CBRinter.commandBufferSave, mapBuffer)
			}
			_, ok := CBRinter.commandBufferSave[playerNr][commandId]
			if !ok {
				CBRinter.commandBufferSave[playerNr][commandId] = -1
			}
			if commandState == -1 {
				CBRinter.commandBufferSave[playerNr][commandId] = -1
			}
			if commandState == 0 && CBRinter.commandBufferSave[playerNr][commandId] == -1 {
				CBRinter.commandBufferSave[playerNr][commandId] = 0
				aiData.rawFrames.setCharCommands(playerNr, commandId, commandState)
			}
			if commandState == 1 {
				if CBRinter.commandBufferSave[playerNr][commandId] == -1 {
					aiData.rawFrames.setCharCommandsPrevFrame(playerNr, commandId, 0, -1)
				}
				aiData.rawFrames.setCharCommandsPrevFrame(playerNr, commandId, commandState, execId)
				CBRinter.commandBufferSave[playerNr][commandId] = -1
			}

		}
	}

	return true
}

func ReplayRecordGenericVars(playerNr int, genericInt []int32, genericFloat []float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setGenericVars(playerNr, genericInt, genericFloat)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setGenericVars(playerNr, genericInt, genericFloat)
	}
	return true
}
func ReplayRecordIkemenSpecific(playerNr int, MoveID int32, MoveFrame int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setIkemenSpecific(playerNr, MoveID, MoveFrame)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setIkemenSpecific(playerNr, MoveID, MoveFrame)
	}
	return true
}

func ReplayRecordFrameAdv(playerNr int, frameAdv int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setFrameAdv(playerNr, frameAdv)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setFrameAdv(playerNr, frameAdv)
	}
	return true
}

func ReplayRecordComboInfo(playerNr int, movesUsed int32, pressure bool) bool {
	if aiData.recording == true {
		aiData.rawFrames.setComboInfo(playerNr, movesUsed, pressure)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.setComboInfo(playerNr, movesUsed, pressure)
	}
	return true
}

func ReplayRecordCommandBuffer(playerNr int, CommandID string, commandIndex int32, execute bool) bool {
	return true
}

/*
func ReplayRecordInputBuffer(playerNr int, InputId1 int32,  InputId2 int32, InputId3 int32, CommandIndex int32, ResetTimer int32, BufferTimer int32, TameIndex int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.setInputBuffer(playerNr, InputId1,  InputId2, InputId3, CommandIndex, ResetTimer, BufferTimer, TameIndex)
	}
	if CheckCBRReplaying(){
		aiData.curGamestate.setInputBuffer(playerNr, InputId1,  InputId2, InputId3, CommandIndex, ResetTimer, BufferTimer, TameIndex)
	}
	return true
}*/

//---Adding Helper relevant data into a frame in a replay
func HelperReplayRecordPosition(playerNr int, positionX float32, positionY float32, facing float32) bool {
	bFacing := floatFacingToBool(facing)
	if aiData.recording == true {
		aiData.rawFrames.helperSetPosition(playerNr, positionX, positionY, bFacing)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetPosition(playerNr, positionX, positionY, bFacing)
	}
	return true
}
func HelperReplayRecordState(playerNr int, movementState int32, attackState int32, controllable bool, stun int32) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetState(playerNr, (movementState&MT_Standing) > 0, (movementState&MT_Crouching) > 0, (movementState&MT_Airborne) > 0, (movementState&MT_LyingDown) > 0, (attackState&AS_Idle) > 0, (attackState&AS_Hit) > 0, (attackState&AS_Attack) > 0, controllable, stun)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetState(playerNr, (movementState&MT_Standing) > 0, (movementState&MT_Crouching) > 0, (movementState&MT_Airborne) > 0, (movementState&MT_LyingDown) > 0, (attackState&AS_Idle) > 0, (attackState&AS_Hit) > 0, (attackState&AS_Attack) > 0, controllable, stun)
	}
	return true
}
func HelperReplayRecordFramedata(playerNr int, currentMoveFrame int32, currentMoveReferenceID int64) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetFramedata(playerNr, currentMoveFrame, currentMoveReferenceID)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetFramedata(playerNr, currentMoveFrame, currentMoveReferenceID)
	}
	return true
}
func HelperReplayRecordAttackHit(playerNr int, moveHit bool, moveGuarded bool) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetAttackHit(playerNr, moveHit, moveGuarded)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetAttackHit(playerNr, moveHit, moveGuarded)
	}
	return true
}
func HelperReplayRecordMeters(playerNr int, lifePercentage float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetMeters(playerNr, lifePercentage)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetMeters(playerNr, lifePercentage)
	}
	return true
}
func HelperReplayRecordVelocity(playerNr int, horizontalVelocity float32, verticalVelocity float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetVelocity(playerNr, horizontalVelocity, verticalVelocity)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetVelocity(playerNr, horizontalVelocity, verticalVelocity)
	}
	return true
}

func HelperReplayRecordGenericVars(playerNr int, helperID int32, genericInt []int32, genericFloat []float32) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetGenericVars(playerNr, helperID, genericInt, genericFloat)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetGenericVars(playerNr, helperID, genericInt, genericFloat)
	}
	return true

}

func HelperReplayRecordCollisionBoxes(playerNr int, hitbox bool, hurtbox bool) bool {
	if aiData.recording == true {
		aiData.rawFrames.helperSetCollisionBoxes(playerNr, hitbox, hurtbox)
	}
	if CheckCBRReplaying() {
		aiData.curGamestate.helperSetCollisionBoxes(playerNr, hitbox, hurtbox)
	}
	return true

}

//saveFile := cbr.ReplayFile[0].CharName[cbr.ReplayFile[0].CbrFocusCharNr]

//utility functions
//facing "TRUE" == facing to the right
func facingToBool(i int32) bool {
	var bFacing bool
	if i == -1 {
		bFacing = false
	} else if i == 1 {
		bFacing = true
	} else {
		print("There are more than 2 facings you dumbfuck go fix the function facingToBool")
	}
	return bFacing
}

func floatFacingToBool(f float32) bool {
	i := int(f)
	var bFacing bool
	if i == -1 {
		bFacing = false
	} else if i == 1 {
		bFacing = true
	} else {
		print("There are more than 2 facings you dumbfuck go fix the function facingToBool")
	}
	return bFacing
}

func floatBoolToFacing(b bool) float32 {
	var iFacing float32
	if b == false {
		iFacing = -1
	} else if b == true {
		iFacing = 1
	} else {
		print("There are more than 2 facings you dumbfuck go fix the function boolToFacing")
	}
	return iFacing
}

func boolToFacing(b bool) int32 {
	var iFacing int32
	if b == false {
		iFacing = -1
	} else if b == true {
		iFacing = 1
	} else {
		print("There are more than 2 facings you dumbfuck go fix the function boolToFacing")
	}
	return iFacing
}

///Interface functions for the outside program-----------

type CBRVariableImportance struct {
	Float                bool
	VarNr                int32
	HelperID             int32
	DifferenceRange      float64
	VariableIncrements   float64
	MaxDissimilarityCost float32
	DivisionIntervals    float64
}

var ikemenVarImport map[int][]CBRVariableImportance

func AddIkemenVars(val CBRVariableImportance, playerNr int) {
	val.DivisionIntervals = math.Floor(val.DifferenceRange / val.VariableIncrements)
	if ikemenVarImport == nil {
		ikemenVarImport = map[int][]CBRVariableImportance{}
	}
	if ikemenVarImport[playerNr] == nil {
		ikemenVarImport[playerNr] = []CBRVariableImportance{}
	}

	ikemenVarImport[playerNr] = append(ikemenVarImport[playerNr], val)
}
