package types

type Config struct {
	WorkingDirectory string
	// Mode
	Mode int32 
	HibernateMinWait int32 
	HibernateMaxWait int32 
	HibernateActivityLength int32 
	// Annoyances
	Annoyances   bool  
	TimerDelay   int32

	AnnoyancePopups  bool
	PopupChance  int32 
	PopupOpacity int32 
	PopupDenialMode bool 
	PopupDenialChance int32 
	PopupMitosis bool 
	PopupMitosisStrength int32 
	PopupTimeout bool 
	PopupTimeoutDelay int32 

	AnnoyanceVideos  bool 
	VideoChance      int32 
	VideoVolume      int32 

	AnnoyancePrompts  bool
	PromptChance      int32 
	MaxMistakesToggle bool 
	MaxMistakes       int32

	AnnoyanceAudio   bool 
	AudioChance      int32 
	AudioVolume      int32 
	// Package
	LoadedPackage     string
	LoadedPackageData *EdgewarePackage
	// Other
	StartOnBoot bool 
	RunOnExit   bool
}