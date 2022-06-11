package main

/*---FILE DESCRIPTION---
The parameter file contains structs with parameters that changes how the CBR AI interprets or works with the gamestate and its cases.
This file alo contains the aiData struct which is where general CBRAI data is stored.
---FILE DESCRIPTION---*/

var cbrParameters = CBRParameters{
	gameStateMaxChange:       5,    //every frame during gameplay if we check the current gamestate against the currently running case and it got this much worse, we check for a new case
	betterCaseThreshold:      0.30, //when the next case in a sequence would be played, if another case exists that is this much better, switch to that case
	topSelectionThreshold:    0.00, //when a best case is found to avoid always selecting the same case in similar situation choose a random worse. This parameter determines how much worse at most.
	maxXPositionComparison:   300,  //If X position comparisons are more than this amount shifted, similarity for them is 1
	maxYPositionComparison:   200,  //If Y position comparisons are more than this amount shifted, similarity for them is 1
	maxVelocityComparison:    10,   //The max difference in velocity before the comparison function hits max dissimilarity
	curGamestateQueLength:    12,   // how many frames in a row are stored to compare against in the comparison function. Stored in aiData.curGamestate
	maxInputBufferDifference: 30,   //How many frames of input buffering can be off before reaching max dissimilarity
	maxHitstunDifference:     10,   //how many frames of difference of beeing in hitstun is allowed till max dissimilarity
	maxBlockstunDifference:   10,   //how many frames of difference of beeing in blockstun is allowed till max dissimilarity
	maxAttackStateDiff:       20,   //how many frames of beeing in an attack state are allowed till max dissimilarity
	nearWallDist:             0.13, //percent of how close compared to current stage size a character has to be, to be considered near the wall
	repetitionFrames:         60,   //amount of frames after which a case was used where it will be taxed for beeing used again. Multiplied the more a case is used.
	comboLength:              20,

	cps: comparisonParameters{
		//parameters that determine how strongly different comparison functions are evaluated
		XRelativePosition:    1.0,
		YRelativePosition:    1.0,
		xVelocityComparison:  0.25,
		yVelocityComparison:  0.25,
		inputBufferDirection: 1.0,
		inputBufferButton:    1.0,
		airborneState:        1.0,
		lyingDownState:       1.0,
		hitState:             1.0,
		blockState:           1.0,
		attackState:          1.0,
		nearWall:             0.3,
		moveID:               0.5,
		pressureMoveID:       0.8,
		getHit:               1.0,
		didHit:               1.0,
		frameAdv:             0.3,
		frameAdvInitiator:    0.1,
		comboSimilarity:      1.0,

		objectOrder: 0.3,
		caseReuse:   0.5,
		roundState:  100.0,

		helperRelativePositionX:   0.5,
		helperRelativePositionY:   0.5,
		helperXVelocityComparison: 0.25,
		helperYVelocityComparison: 0.25,

		enemyXVelocityComparison: 0.25,
		enemyYVelocityComparison: 0.25,
		enemyAirborneState:       1.0,
		enemyLyingDownState:      1.0,
		enemyHitState:            1.0,
		enemyBlockState:          1.0,
		enemyAttackState:         1.0,
		enemyMoveID:              0.5,
		enemyPressureMoveID:      0.8,

		enemyHelperRelativePositionX:   1.0,
		enemyHelperRelativePositionY:   1.0,
		enemyHelperXVelocityComparison: 0.25,
		enemyHelperYVelocityComparison: 0.25,
	},
}

//All data relevant for the AI to operate
//See CBRData.proto and CBRRawFrames.proto for the structure of the data
var aiData = AIData{
	cbrData:       &CBRData{},
	bufferCbrData: nil,
	rawFrames:     &CBRRawFrames{},
	recording:     false,
	replaying:     false,
	//replayIndex: 0,
	//rawFrameReplay: false,
	curGamestate:      &CBRRawFrames{},
	framedata:         &Framedata{},
	aiControlledIndex: -1,
}

type AIData struct {
	cbrData       *CBRData //Data storage after the AI processed a replay
	bufferCbrData *CBRData
	rawFrames     *CBRRawFrames //Data storage of a replay before processing
	recording     bool
	replaying     bool
	//replayFrames []*CBRData_Frame	//Frames that the AI is ready to send over to the game for replaying
	//replayIndex int
	//rawFrameReplay bool
	curGamestate      *CBRRawFrames
	framedata         *Framedata
	aiControlledIndex int // index denoting the character that the CBRAI controlls.

}

//parameters that are relevant for the comparison function or need to be adjusted when the comparison function is adjusted
type CBRParameters struct {
	gameStateMaxChange       float32
	betterCaseThreshold      float32
	topSelectionThreshold    float32
	maxXPositionComparison   float32
	maxYPositionComparison   float32
	maxVelocityComparison    float32
	cps                      comparisonParameters
	curGamestateQueLength    int
	maxInputBufferDifference int32
	maxHitstunDifference     int32
	maxBlockstunDifference   int32
	maxAttackStateDiff       int32
	nearWallDist             float32
	repetitionFrames         int64
	comboLength              float32
}

//parameters that determine how important the corresponding sub comparison functions are
type comparisonParameters struct {
	XRelativePosition    float32
	YRelativePosition    float32
	xVelocityComparison  float32
	yVelocityComparison  float32
	inputBufferDirection float32
	inputBufferButton    float32
	airborneState        float32
	lyingDownState       float32
	hitState             float32
	blockState           float32
	attackState          float32
	nearWall             float32
	unitOrder            float32
	moveID               float32
	pressureMoveID       float32
	getHit               float32
	didHit               float32
	frameAdv             float32
	frameAdvInitiator    float32
	comboSimilarity      float32

	objectOrder float32
	caseReuse   float32
	roundState  float32

	helperRelativePositionX   float32
	helperRelativePositionY   float32
	helperXVelocityComparison float32
	helperYVelocityComparison float32

	enemyXVelocityComparison float32
	enemyYVelocityComparison float32
	enemyAirborneState       float32
	enemyLyingDownState      float32
	enemyHitState            float32
	enemyBlockState          float32
	enemyAttackState         float32
	enemyMoveID              float32
	enemyPressureMoveID      float32

	enemyHelperXVelocityComparison float32
	enemyHelperYVelocityComparison float32
	enemyHelperRelativePositionX   float32
	enemyHelperRelativePositionY   float32
}
