package cbr

/*---FILE DESCRIPTION---
Contains functions to insert data from the game into the AIData
---FILE DESCRIPTION---*/

//Naming for inputs, individual bits represent pressed buttons 0001 = 1 = up, 0011 = 3 = up + down
const (
	IB_PU int32 = 1 << iota
	IB_PD
	IB_PL
	IB_PR
	IB_A
	IB_B
	IB_C
	IB_X
	IB_Y
	IB_Z
	IB_S
	IB_D
	IB_W
	IB_M
	IB_anybutton = IB_A | IB_B | IB_C | IB_X | IB_Y | IB_Z | IB_S | IB_D | IB_W | IB_M
)

///Interface functions for the outside program-----------

// switch left/right input for mirrored replay
func swapBitsAtPos(byte int32, posA int, posB int) int32 {
	x := byte
	p1 := posA
	p2 := posB
	n := 1
	set1 := (x >> p1) & ((1 << n) - 1)
	set2 := (x >> p2) & ((1 << n) - 1)
	xor1 := set1 ^ set2
	xor := (xor1 << p1) | (xor1 << p2)
	result := x ^ xor
	return result
}

//Adds empty replay
func (x *CBRRawFrames) AddReplay() bool {
	//appending a new replay file to the CBR replay data everytime a new recording is started.
	x.ReplayFile = append(x.ReplayFile, &CBRRawFrames_ReplayFile{})
	return true
}

//Adds empty Frame
func (x *CBRRawFrames) AddFrame() bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		//x.ReplayFile[replayId].Frame = append(x.ReplayFile[replayId].Frame, &CBRRawFrames_Frame{})
		x.ReplayFile[replayId].addFrame()
		//cbrData.ReplayFile[replayId].Frame[frameId].CharData = append(cbrData.ReplayFile[replayId].Frame[frameId].CharData, &CBRData_CharData{Input: input})
	}
	return true
}

//Adds empty Frame
func (x *CBRRawFrames) queueFrame(queLength int) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		if len(x.ReplayFile[replayId].Frame) > queLength {
			x.queueRemoveFrame(0)
		} else {
			x.ReplayFile[replayId].addFrame()
		}

	}
	return true
}
func (x *CBRRawFrames_ReplayFile) addFrame() bool {
	if x != nil {
		x.Frame = append(x.Frame, &CBRRawFrames_Frame{})
		x.Frame[len(x.Frame)-1].WorldCBRComparisonData = &WorldCBRComparisonData{}
		//cbrData.ReplayFile[replayId].Frame[frameId].CharData = append(cbrData.ReplayFile[replayId].Frame[frameId].CharData, &CBRData_CharData{Input: input})
	}
	return true
}

//Removes a frame from an array, which is beeing treated as a que, makes sure the que stays the same length
func (x *CBRRawFrames) queueRemoveFrame(index int) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		for i := index; i < len(x.ReplayFile[replayId].Frame); i++ {
			if i == len(x.ReplayFile[replayId].Frame)-1 {
				x.ReplayFile[replayId].Frame[i] = &CBRRawFrames_Frame{}
				x.ReplayFile[replayId].Frame[i].WorldCBRComparisonData = &WorldCBRComparisonData{}
			} else {
				x.ReplayFile[replayId].Frame[i] = x.ReplayFile[replayId].Frame[i+1]
			}
		}
	}
	return true
}

//Adds empty CharData to Frame, every frame has data from multiple players stored
func (x *CBRRawFrames) AddCharData(players int) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		x.ReplayFile[replayId].Frame[frameId].AddCharData(players)
		/*
			for i := 0; i < players; i++ {
				x.ReplayFile[replayId].Frame[frameId].CharData = append(x.ReplayFile[replayId].Frame[frameId].CharData, &CBRRawFrames_CharData{})
			}*/
	}
	return true
}
func (x *CBRRawFrames_Frame) AddCharData(players int) bool {
	//cbrData.ReplayFile[replayId].Frame[frameId].CharData[playerNr].Input = input
	for i := 0; i < players; i++ {
		x.CharData = append(x.CharData, &CBRRawFrames_CharData{})
		x.CharData[len(x.CharData)-1].ComparisonData = &ComparisonData{}
	}
	return true
}
func (x *CBRRawFrames) addHelperData(charNr int) bool {
	lastReplay := len(x.ReplayFile) - 1
	frameId := len(x.ReplayFile[lastReplay].Frame) - 1

	x.ReplayFile[lastReplay].Frame[frameId].CharData[charNr].HelperData = append(x.ReplayFile[lastReplay].Frame[frameId].CharData[charNr].HelperData, &CBRRawFrames_HelperData{})
	x.ReplayFile[lastReplay].Frame[frameId].CharData[charNr].HelperData[len(x.ReplayFile[lastReplay].Frame[frameId].CharData[charNr].HelperData)-1].CompData = &HelperComparisonData{}
	return true
}

func (x *CBRRawFrames) setCharData(cbrFocusCharNr int32, charName string, charTeam int32) bool {
	if x.ReplayFile != nil {
		x.ReplayFile[len(x.ReplayFile)-1].CharName = append(x.ReplayFile[len(x.ReplayFile)-1].CharName, charName)
		x.ReplayFile[len(x.ReplayFile)-1].CharTeam = append(x.ReplayFile[len(x.ReplayFile)-1].CharTeam, charTeam)
		x.ReplayFile[len(x.ReplayFile)-1].CbrFocusCharNr = cbrFocusCharNr
	}
	return true
}

func (x *CBRRawFrames) setStageData(leftWallPos float32, rightWallPos float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		x.ReplayFile[replayId].Frame[frameId].WorldCBRComparisonData.StageData = &StageData{LeftWallPos: leftWallPos, RightWallPos: rightWallPos}
	}

	return true
}
func (x *CBRRawFrames) setRoundState(roundState int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		x.ReplayFile[replayId].Frame[frameId].WorldCBRComparisonData.RoundState = roundState
	}

	return true
}

//------------------Sets the data like player input and character facing in the the CharData of a frame
func (x *CBRRawFrames) setPlayerInput(playerNr int, input int32, facing bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].Input = input
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].Facing = facing
		}
	}
	return true
}
func (x *CBRRawFrames) setCharState(playerNr int, MStateStanding bool, MStateCrouching bool, MStateAir bool, MStateLying bool, AStateIdle bool, AStateHit bool, AStateAttack bool, controllable bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MStateStanding = MStateStanding
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MStateCrouching = MStateCrouching
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MStateAir = MStateAir
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MStateLying = MStateLying
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.AStateIdle = AStateIdle
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.AStateHit = AStateHit
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.AStateAttack = AStateAttack
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Controllable = controllable
		}
	}
	return true
}
func (x *CBRRawFrames) setFramedata(playerNr int, currentMoveFrame int32, currentMoveReferenceID int64) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.CurrentMoveFrame = currentMoveFrame
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.CurrentMoveReferenceID = currentMoveReferenceID

		}
	}
	return true
}

func (x *CBRRawFrames) setMeters(playerNr int, lifePercentage float32, meterPercentage float32, meterMax float32, dizzyPercentage float32, guardPointsPercentage float32, recoverableHpPercentage float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].LifePercentage = lifePercentage
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].MeterPercentage = meterPercentage
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].MeterMax = meterMax
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].DizzyPercentage = dizzyPercentage
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].GuardPointsPercentage = guardPointsPercentage
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].RecoverableHpPercentage = recoverableHpPercentage
		}
	}
	return true
}

func (x *CBRRawFrames) setVelocity(playerNr int, horizontalVelocity float32, verticalVelocity float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Velocity = &Velocity{}
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Velocity.XVel = horizontalVelocity
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Velocity.YVel = verticalVelocity
		}
	}
	return true
}

func (x *CBRRawFrames) setStun(playerNr int, blockstun int32, hitStun int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Blockstun = blockstun
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.HitStun = hitStun
		}
	}
	return true
}
func (x *CBRRawFrames) setSelfHit(playerNr int, selfGuard bool, selfHit bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.SelfGuard = selfGuard
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.SelfHit = selfHit
		}
	}
	return true
}

func (x *CBRRawFrames) setAttackHit(playerNr int, moveHit bool, moveGuarded bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MoveHit = moveHit
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.MoveGuarded = moveGuarded
		}
	}
	return true
}
func (x *CBRRawFrames) setPosition(playerNr int, positionX float32, positionY float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.CharPos = &Position{}
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.CharPos.XPos = positionX
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.CharPos.YPos = positionY
		}
	}
	return true
}
func (x *CBRRawFrames) setGenericVars(playerNr int, genericInt []int32, genericFloat []float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {

			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].GenericFloatVars = append(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].GenericFloatVars, genericFloat...)
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].GenericIntVars = append(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].GenericIntVars, genericInt...)
		}
	}
	return true
}

func (x *CBRRawFrames) setIkemenSpecific(playerNr int, MoveID int32, MoveFrame int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].IkemenMoveID = MoveID
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].IkemenMoveFrame = MoveFrame
		}
	}
	return true
}

func (x *CBRRawFrames) setFrameAdv(playerNr int, frameAdv int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.FrameAdv = frameAdv
		}
	}
	return true
}

func (x *CBRRawFrames) setComboInfo(playerNr int, movesUsed int32, pressure bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.ComboMovesUsed = movesUsed
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.Pressure = pressure
		}
	}
	return true
}

func (x *CBRRawFrames) setInputBuffer(playerNr int, directionBuffer []int32, buttonBuffer []int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			buffer := InputBuffer{InputDirection: directionBuffer, InputButton: buttonBuffer}
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ComparisonData.InputBuffer = &buffer
		}
	}
	return true
}

func (x *CBRRawFrames) setCharCommands(playerNr int, CommandId string, commandState int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			if x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution == nil {
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution = make(map[string]int32)
			}

			val, ok := x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId]
			if commandState != 1 {
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId] = commandState
			} else if !ok || val != 0 {
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId] = commandState
			}
		}
	}
	return true
}

//Makeshift fix for my commandBuffer always detecting a executed command 1 frame too late
func (x *CBRRawFrames) setCharCommandsPrevFrame(playerNr int, CommandId string, commandState int32, execId int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 2
		if frameId >= 0 && len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr {
			if x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution == nil {
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution = make(map[string]int32)
			}

			val, ok := x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId]
			//if the command"Executed" is inserted, check if for that command the "Start" or "Execute" is already in the frame and weather it is the opposite command as the one we are inserting...
			if !ok || val == commandState {
				//... if not add the "executed" command normally
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId] = commandState
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ExecutionConditions = execId
			} else {
				//... if so add a "ExecutedAndStart" command for the caseGenerator to discern later
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].CommandExecution[CommandId] = cmdEx_ExecutedAndStart // cmdEx_ExecutedAndStart means the command was both executed and started
				x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].ExecutionConditions = execId
			}
		}
	}

	return true
}

/*
func (x *CBRRawFrames) setInputBuffer(playerNr int, InputId1 int32,  InputId2 int32, InputId3 int32, CommandIndex int32, ResetTimer int32, BufferTimer int32, TameIndex int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		inputBuffer := InputBuffer{InputId1: InputId1, InputId2: InputId2, InputId3: InputId3, CommandIndex: CommandIndex, ResetTimer: ResetTimer, BufferTimer: BufferTimer, TameIndex: TameIndex}
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr{
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].InputBuffer = append(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].InputBuffer, &inputBuffer)
		}
	}
	return true
}
*/

//------------------Sets the data of a helper object for a frame
func (x *CBRRawFrames) helperSetPosition(playerNr int, positionX float32, positionY float32, facing bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.PositionX = positionX
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.PositionY = positionY
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].Facing = facing
		}
	}
	return true
}

func (x *CBRRawFrames) helperSetState(playerNr int, MStateStanding bool, MStateCrouching bool, MStateAir bool, MStateLying bool, AStateIdle bool, AStateHit bool, AStateAttack bool, controllable bool, stun int32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MStateStanding = MStateStanding
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MStateCrouching = MStateCrouching
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MStateAir = MStateAir
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MStateLying = MStateLying
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.AStateIdle = AStateIdle
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.AStateHit = AStateHit
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.AStateAttack = AStateAttack
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].Controllable = controllable
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].HitStun = stun
		}
	}
	return true
}
func (x *CBRRawFrames) helperSetFramedata(playerNr int, currentMoveFrame int32, currentMoveReferenceID int64) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.CurrentMoveFrame = currentMoveFrame
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.CurrentMoveReferenceID = currentMoveReferenceID
		}
	}
	return true
}
func (x *CBRRawFrames) helperSetAttackHit(playerNr int, moveHit bool, moveGuarded bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MoveHit = moveHit
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].MoveGuarded = moveGuarded
		}
	}
	return true
}
func (x *CBRRawFrames) helperSetMeters(playerNr int, lifePercentage float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].LifePercentage = lifePercentage
		}
	}
	return true
}
func (x *CBRRawFrames) helperSetVelocity(playerNr int, horizontalVelocity float32, verticalVelocity float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.Velocity = &Velocity{}
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.Velocity.XVel = horizontalVelocity
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.Velocity.YVel = verticalVelocity
		}
	}
	return true
}

func (x *CBRRawFrames) helperSetGenericVars(playerNr int, helperID int32, genericInt []int32, genericFloat []float32) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.HelperID = helperID
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].GenericFloatVars = append(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].GenericFloatVars, genericFloat...)
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].GenericIntVars = append(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].GenericIntVars, genericInt...)
		}
	}
	return true
}

func (x *CBRRawFrames) helperSetCollisionBoxes(playerNr int, hitbox bool, hurtbox bool) bool {
	if x.ReplayFile != nil {
		var replayId = len(x.ReplayFile) - 1
		var frameId = len(x.ReplayFile[replayId].Frame) - 1
		var helperNr = len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) - 1
		if len(x.ReplayFile[replayId].Frame[frameId].CharData) > playerNr && len(x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData) > helperNr {
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.HitboxOut = hitbox
			x.ReplayFile[replayId].Frame[frameId].CharData[playerNr].HelperData[helperNr].CompData.HurtboxOut = hurtbox
		}
	}
	return true
}
