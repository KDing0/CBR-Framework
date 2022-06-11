package cbr

/*---FILE DESCRIPTION---
Contains functions that control the managing of AI relevant data and to start/end recording or turn the AI on/off
---FILE DESCRIPTION---*/

//calls the cbrAILoop to figure out what input to replay. If a new case is selected also sends the input buffer
func readCbrFrameInput(facing int32) (input int32) {
	input, storedFacing := cbrLoop()
	bFacing := facingToBool(facing)
	if bFacing != storedFacing {
		input = swapBitsAtPos(input, 2, 3)
	}
	return input
}

func setAIControlledCharacter(charIndex int) {
	aiData.aiControlledIndex = charIndex
}

func getAIControlledCharacter() (charIndex int) {
	return aiData.aiControlledIndex
}

func (x *CBRRawFrames) clearRecording() bool {
	x.ReplayFile = nil
	return true
}

//Starts recording in a new replay File
func (x *AIData) InitializeRecording() bool {
	if x.rawFrames.ReplayFile == nil {
		x.rawFrames = &CBRRawFrames{ReplayFile: []*CBRRawFrames_ReplayFile{}}
	}
	x.recording = true
	return true
}
func (x *CBRRawFrames) initalize(replayIndex int) bool {
	if x.ReplayFile == nil {
		x.ReplayFile = append(x.ReplayFile, &CBRRawFrames_ReplayFile{})
	} else {
		if len(x.ReplayFile) > replayIndex {
			x.ReplayFile[replayIndex] = &CBRRawFrames_ReplayFile{}
		}
	}

	return true
}

func (x *AIData) EndRecording() bool {
	x.recording = false
	return true
}

func (x *AIData) endCBRReplaying() bool {
	x.replaying = false
	return true
}
func (x *AIData) startCBRReplaying(cbrFocusCharNr int32, charName []string, charTeam []int32, framedata *Framedata) bool {
	aiData.curGamestate.initalize(len(aiData.curGamestate.ReplayFile) - 1)
	resetCBRProcess()
	aiData.replaying = true
	for i := range charName {
		aiData.curGamestate.setCharData(cbrFocusCharNr, charName[i], charTeam[i])
	}
	aiData.framedata = framedata
	return true
}
func (x *AIData) startCBRRecording(cbrFocusCharNr int32, charName []string, charTeam []int32) {
	aiData.InitializeRecording()
	aiData.rawFrames.AddReplay()
	for i := range charName {
		aiData.rawFrames.setCharData(cbrFocusCharNr, charName[i], charTeam[i])
	}
}
func (x *AIData) endCBRRecording(framedata *Framedata, directory string) {
	//AiData.rawFrames.clearRecording()

	if aiData.bufferCbrData != nil {
		aiData.saveCBRBuffer(directory)
		aiData.bufferCbrData = nil
	}

	//replace the last replay file if we were recording with midFightLearning active
	if midFightLearning.active == true {
		aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1] = nil
		replay := aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1].rawFramesToCBRReplay(framedata)
		aiData.rawFrames.ReplayFile = []*CBRRawFrames_ReplayFile{}
		aiData.cbrData = loadCBRData(directory, replay.CharName[replay.CbrFocusCharNr]+"_"+getPlayerName(int(replay.CbrFocusCharNr)))
		aiData.EndRecording()
		aiData.cbrData.insertReplaytoCaseData(replay)
		setCbrMetaData(aiData.cbrData, replay)
		saveCBRData(aiData.cbrData, directory, aiData.cbrData.CharName+"_"+aiData.cbrData.PlayerName) //aiData.cbrData.ReplayFile[0].CharName[aiData.cbrData.ReplayFile[0].CbrFocusCharNr]
	} else {
		replay := aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1].rawFramesToCBRReplay(framedata)
		aiData.rawFrames.ReplayFile = []*CBRRawFrames_ReplayFile{}
		cbrData := loadCBRData(directory, replay.CharName[replay.CbrFocusCharNr]+"_"+getPlayerName(int(replay.CbrFocusCharNr)))
		aiData.EndRecording()
		cbrData.insertReplaytoCaseData(replay)
		setCbrMetaData(cbrData, replay)
		saveCBRData(cbrData, directory, cbrData.CharName+"_"+cbrData.PlayerName) //cbrData.ReplayFile[0].CharName[cbrData.ReplayFile[0].CbrFocusCharNr]
	}
}

func (x *AIData) bufferStoreCBRRecording(framedata *Framedata, directory string) {
	//AiData.rawFrames.clearRecording()

	//replace the last replay file if we were recording with midFightLearning active
	if midFightLearning.active == true {
		aiData.cbrData.ReplayFile[len(aiData.cbrData.ReplayFile)-1] = nil
		replay := aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1].rawFramesToCBRReplay(framedata)
		aiData.rawFrames.ReplayFile = []*CBRRawFrames_ReplayFile{}
		if aiData.bufferCbrData == nil {
			aiData.bufferCbrData = loadCBRData(directory, replay.CharName[replay.CbrFocusCharNr]+"_"+getPlayerName(int(replay.CbrFocusCharNr)))
		}
		aiData.EndRecording()
		aiData.bufferCbrData.insertReplaytoCaseData(replay)
		setCbrMetaData(aiData.bufferCbrData, replay)
		aiData.bufferCbrData = aiData.cbrData
		//saveCBRData(aiData.cbrData, directory, aiData.cbrData.CharName+"_"+aiData.cbrData.PlayerName)//aiData.cbrData.ReplayFile[0].CharName[aiData.cbrData.ReplayFile[0].CbrFocusCharNr]
	} else {
		replay := aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1].rawFramesToCBRReplay(framedata)
		aiData.rawFrames.ReplayFile = []*CBRRawFrames_ReplayFile{}
		if aiData.bufferCbrData == nil {
			aiData.bufferCbrData = loadCBRData(directory, replay.CharName[replay.CbrFocusCharNr]+"_"+getPlayerName(int(replay.CbrFocusCharNr)))
		}
		setCbrMetaData(aiData.bufferCbrData, replay)
		aiData.EndRecording()
		aiData.bufferCbrData.insertReplaytoCaseData(replay)

		//setCbrMetaData(cbrData, replay)
		//saveCBRData(cbrData, directory, cbrData.CharName+"_"+cbrData.PlayerName)//cbrData.ReplayFile[0].CharName[cbrData.ReplayFile[0].CbrFocusCharNr]
	}
}

func (x *AIData) saveCBRBuffer(directory string) {
	//AiData.rawFrames.clearRecording()
	if aiData.bufferCbrData != nil {
		saveCBRData(aiData.bufferCbrData, directory, aiData.bufferCbrData.CharName+"_"+aiData.bufferCbrData.PlayerName)
	}
}

func (x *AIData) discardCBRRecording() {
	//AiData.rawFrames.clearRecording()
	aiData.bufferCbrData = nil
	aiData.rawFrames.ReplayFile[len(aiData.rawFrames.ReplayFile)-1] = nil
	aiData.rawFrames.ReplayFile = []*CBRRawFrames_ReplayFile{}
	aiData.EndRecording()
}

func setCbrMetaData(cbr *CBRData, lastReplay *CBRData_ReplayFile) {
	playerName := ""

	val, ok := CBRinter.playerNames[int(lastReplay.CbrFocusCharNr)]

	if ok {
		playerName = val
	}
	cbr.PlayerName = playerName
	cbr.CharName = lastReplay.CharName[lastReplay.CbrFocusCharNr]
}

/*
func (x *CBRData) Initalize() bool {
	if x == nil{
		print("CBRRcord start")
		data, err := ioutil.ReadFile("CBRReplays.data")
		if err != nil {
			x = &CBRData{
				ReplayFile: []*CBRData_ReplayFile{},
			}
		}else{
			x = &CBRData{
				ReplayFile: []*CBRData_ReplayFile{},
			}
			err := proto.Unmarshal(data, x)
			if err != nil {
				log.Fatal("unmarshaling error: ", err)
			}
		}
	}else{
		x.CBRClose()
	}
	return true
}

func (x *CBRData) CBRClose() bool {
	print("CBRRcord end")
	data, err := proto.Marshal(x)
	err = proto.Unmarshal(data, x)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	ioutil.WriteFile("CBRReplays.data", data, 0644)
	x = nil
	return true
}
*/
/*
//Converts unprocessed frames to an array ready for replaying. Only playerNr data is kept.
func RawFramesToReplay(frames CBRRawFrames, playerNr int, replayNr int) []*CBRData_Frame {
	var replayFrames []*CBRData_Frame
	for i := 0; i < len(frames.ReplayFile[replayNr].Frame); i++ {
		in := frames.ReplayFile[replayNr].Frame[i].CharData[playerNr].Input
		facing := frames.ReplayFile[replayNr].Frame[i].CharData[playerNr].Facing
		fr := CBRData_Frame{Input: in, Facing: facing}
		replayFrames = append(replayFrames, &fr)
	}
	return replayFrames
}


//Start replaying with given given array of frames
func (x *AIData) InitializeReplaying(frames []*CBRData_Frame) bool {
	ret := false
	if frames != nil && len(frames) > 0 {
		x.replayFrames = frames
		x.replaying = true
		x.replayIndex = 0
		ret = true
	}
	return ret
}
*/
