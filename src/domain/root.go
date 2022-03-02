package domain

import (
	"github.com/lib/pq"
)

var Models = []interface{}{
	&Root{}, &Factions{}, &Role{},
	&Score{}, &Achievement{}, &CommunicationLogs{},
	&Deaths{}, &Explosions{}, &ManifestEntries{},
	&LeaveStats{}, &Damage{},
	&Objectives{}, &CultInfo{}, &UplinkInfo{},
	&UplinkPurchases{}, &Aspects{}, &RitenameByCount{},
}

type Root struct {
	RoundID           uint                `json:"round_id" gorm:"primaryKey"`
	Version           uint                `json:"version"`
	Mode              string              `json:"mode"`
	StartTime         string              `json:"start_time"`
	Map               string              `json:"map"`
	Duration          string              `json:"duration"`
	EndTime           string              `json:"end_time"`
	Factions          []Factions          `json:"factions"`
	OrphanedRoles     []Role              `json:"orphaned_roles" gorm:"foreignKey:OwnerID"`
	ModeResult        string              `json:"mode_result"`
	MinimapImage      string              `json:"minimap_image"`
	ServerAddress     string              `json:"server_address"`
	BaseCommitSha     string              `json:"base_commit_sha"`
	TestMerges        string              `json:"test_merges"`
	CompletionHTML    string              `json:"completion_html"`
	Score             Score               `json:"score"`
	Achievements      []Achievement       `json:"achievements"`
	CommunicationLogs []CommunicationLogs `json:"communication_logs"`
	Deaths            []Deaths            `json:"deaths"`
	Explosions        []Explosions        `json:"explosions"`
	ManifestEntries   []ManifestEntries   `json:"manifest_entries"`
	LeaveStats        []LeaveStats        `json:"leave_stats"`
}

type Achievement struct {
	ID     uint
	RootID uint
	Key    string `json:"key"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
}

type Score struct {
	ID             uint
	RootID         uint
	Crewscore      int            `json:"crewscore"`
	Rating         string         `json:"rating"`
	Stuffshipped   int            `json:"stuffshipped"`
	Stuffharvested int            `json:"stuffharvested"`
	Oremined       int            `json:"oremined"`
	Researchdone   int            `json:"researchdone"`
	Eventsendured  int            `json:"eventsendured"`
	Powerloss      int            `json:"powerloss"`
	Mess           int            `json:"mess"`
	Meals          int            `json:"meals"`
	Disease        int            `json:"disease"`
	Deadcommand    int            `json:"deadcommand"`
	Arrested       int            `json:"arrested"`
	Traitorswon    int            `json:"traitorswon"`
	Roleswon       int            `json:"roleswon"`
	Allarrested    int            `json:"allarrested"`
	Opkilled       int            `json:"opkilled"`
	Disc           int            `json:"disc"`
	Nuked          int            `json:"nuked"`
	Destranomaly   int            `json:"destranomaly"`
	RecAntags      int            `json:"rec_antags"`
	CrewEscaped    int            `json:"crew_escaped"`
	CrewDead       int            `json:"crew_dead"`
	CrewTotal      int            `json:"crew_total"`
	CrewSurvived   int            `json:"crew_survived"`
	Captain        pq.StringArray `json:"captain" gorm:"type:text[]"`
	Powerbonus     int            `json:"powerbonus"`
	Messbonus      int            `json:"messbonus"`
	Deadaipenalty  int            `json:"deadaipenalty"`
	Foodeaten      int            `json:"foodeaten"`
	Clownabuse     int            `json:"clownabuse"`
	Richestname    int            `json:"richestname"`
	Richestjob     int            `json:"richestjob"`
	Richestcash    int            `json:"richestcash"`
	Richestkey     int            `json:"richestkey"`
	Dmgestname     int            `json:"dmgestname"`
	Dmgestjob      int            `json:"dmgestjob"`
	Dmgestdamage   int            `json:"dmgestdamage"`
	Dmgestkey      int            `json:"dmgestkey"`
}

type CommunicationLogs struct {
	ID      uint
	RootID  uint
	Time    string `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Type    string `json:"type"`
}

type Damage struct {
	ID       uint
	DeathsID uint
	Brute    float64 `json:"BRUTE"`
	Fire     float64 `json:"FIRE"`
	Toxin    float64 `json:"TOXIN"`
	Oxy      float64 `json:"OXY"`
	Clone    float64 `json:"CLONE"`
	Brain    float64 `json:"BRAIN"`
}

type Deaths struct {
	ID               uint
	RootID           uint
	Name             string  `json:"name"`
	AssignedRole     string  `json:"assigned_role"`
	SpecialRole      string  `json:"special_role"`
	Damage           Damage  `json:"damage"`
	RealName         string  `json:"real_name"`
	MindName         string  `json:"mind_name"`
	DeathX           int     `json:"death_x"`
	DeathY           int     `json:"death_y"`
	DeathZ           int     `json:"death_z"`
	TimeOfDeath      float64 `json:"time_of_death"`
	FromSuicide      int     `json:"from_suicide"`
	LastAttackerName string  `json:"last_attacker_name"`
}

type Explosions struct {
	ID               uint
	RootID           uint
	EpicenterX       int `json:"epicenter_x"`
	EpicenterY       int `json:"epicenter_y"`
	EpicenterZ       int `json:"epicenter_z"`
	DevastationRange int `json:"devastation_range"`
	HeavyImpactRange int `json:"heavy_impact_range"`
	LightImpactRange int `json:"light_impact_range"`
	FlashRange       int `json:"flash_range"`
}

type ManifestEntries struct {
	ID           uint
	RootID       uint
	Name         string `json:"name"`
	AssignedRole string `json:"assigned_role"`
	SpecialRole  string `json:"special_role"`
	AntagRoles   string `json:"antag_roles"`
}

type LeaveStats struct {
	ID           uint
	RootID       uint
	Name         string `json:"name"`
	StartTime    string `json:"start_time"`
	AssignedRole string `json:"assigned_role"`
	SpecialRole  string `json:"special_role"`
	AntagRoles   string `json:"antag_roles"`
	LeaveType    string `json:"leave_type"`
	LeaveTime    string `json:"leave_time"`
}
