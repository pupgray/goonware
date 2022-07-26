package types

type Config struct {
	WorkingDirectory string
	// Mode
	Mode int32 
	HibernateMinWaitMinutes int32 
	HibernateMaxWaitMinutes int32 
	HibernateActivityLength int32 
	LoadedPackage     string
	LoadedPackageData *EdgewarePackage
	StartOnBoot bool 
	RunOnExit   bool

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

	// Drive Filler
	DriveFiller bool
	DriveFillerDelay int32
	DriveFillerBase string
	DriveFillerTags []string
	DriveFillerImageSource int32
	DriveFillerImageUseTags bool
	DriveFillerDownloadMinimumScoreToggle bool
	DriveFillerDownloadMinimumScoreThreshold int32
}