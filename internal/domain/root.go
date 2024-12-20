package domain

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Root struct {
	Date      string     `gorm:"type:date"`
	RoundID   int32      `json:"round_id" gorm:"primaryKey;uniqueIndex"`
	Version   int32      `json:"version"`
	Mode      string     `json:"mode" gorm:"size:128"`
	Bundle    string     `json:"bundle" gorm:"size:128"`
	StartTime string     `json:"start_time" gorm:"size:16"`
	Map       string     `json:"map" gorm:"size:128"`
	Duration  string     `json:"duration" gorm:"size:16"`
	EndTime   string     `json:"end_time"  gorm:"size:16"`
	Factions  []Factions `json:"factions"`
	// Totally broken, and I'm too lazy to fix it
	//OrphanedRoles     []Role              `json:"orphaned_roles" gorm:"foreignKey:OwnerID"`
	ModeResult        string              `json:"mode_result" gorm:"size:128"`
	MinimapImage      string              `json:"minimap_image" gorm:"size:256"`
	ServerAddress     string              `json:"server_address" gorm:"size:256"`
	BaseCommitSha     string              `json:"base_commit_sha" gorm:"size:256"`
	TestMerges        string              `json:"test_merges" gorm:"size:256"`
	CompletionHTML    string              `json:"completion_html"`
	Score             Score               `json:"score"`
	Achievements      []Achievement       `json:"achievements"`
	Medals            []Medal             `json:"medals"`
	CommunicationLogs []CommunicationLogs `json:"communication_logs"`
	Deaths            []Deaths            `json:"deaths"`
	Explosions        []Explosions        `json:"explosions"`
	EMPExplosions     []EMPExplosions     `json:"emps"`
	ManifestEntries   []ManifestEntries   `json:"manifest_entries"`
	LeaveStats        []LeaveStats        `json:"leave_stats"`
	Rating            Rating              `json:"rating"`
	Vote              []Vote              `json:"completed_votes"`
}

type Rating struct {
	ID          int32              `gorm:"uniqueIndex"`
	RootID      int32              `gorm:"index"`
	RatingsTemp map[string]float32 `json:"ratings" gorm:"-"`
	Ratings     []RatingValues
}

type RatingValues struct {
	ID       int32  `gorm:"uniqueIndex"`
	RatingID int32  `gorm:"index"`
	Key      string `gorm:"size:256"`
	Value    float32
}

func (r *Rating) BeforeSave(tx *gorm.DB) error {
	for category, value := range r.RatingsTemp {
		r.Ratings = append(r.Ratings, RatingValues{
			Key:   category,
			Value: value,
		})
	}
	return nil
}

type Vote struct {
	ID          int32          `gorm:"uniqueIndex"`
	RootID      int32          `gorm:"index"`
	Name        string         `json:"name"`
	TotalVotes  int            `json:"total_votes"`
	TotalVoters int            `json:"total_voters"`
	Winner      string         `json:"winner"`
	ResultsTemp map[string]int `json:"results" gorm:"-"`
	Results     []VoteValues
}

type VoteValues struct {
	ID     int32  `gorm:"uniqueIndex"`
	VoteID int32  `gorm:"index"`
	Key    string `gorm:"size:256"`
	Value  int
}

func (r *Vote) BeforeSave(tx *gorm.DB) error {
	for category, value := range r.ResultsTemp {
		r.Results = append(r.Results, VoteValues{
			Key:   category,
			Value: value,
		})
	}
	return nil
}

type Achievement struct {
	ID     int32  `gorm:"uniqueIndex"`
	RootID int32  `gorm:"index"`
	Key    string `json:"key" gorm:"size:256"`
	Name   string `json:"name" gorm:"size:256"`
	Title  string `json:"title" gorm:"size:256"`
	Desc   string `json:"desc" gorm:"size:256"`
}

type Medal struct {
	ID         int32  `gorm:"uniqueIndex"`
	RootID     int32  `gorm:"index"`
	Key        string `json:"key" gorm:"size:256"`
	MedalName  string `json:"medal_name" gorm:"size:256"`
	TargetName string `json:"target_name" gorm:"size:256"`
	ParentName string `json:"parent_name" gorm:"size:256"`
	Reason     string `json:"reason" gorm:"size:2048"`
}

type Score struct {
	ID             int32          `gorm:"uniqueIndex"`
	RootID         int32          `gorm:"index"`
	Crewscore      int32          `json:"crewscore"`
	Rating         string         `json:"rating" gorm:"size:256"`
	Stuffshipped   int32          `json:"stuffshipped"`
	Stuffharvested int32          `json:"stuffharvested"`
	Oremined       int32          `json:"oremined"`
	Researchdone   int32          `json:"researchdone"`
	Eventsendured  int32          `json:"eventsendured"`
	Powerloss      int32          `json:"powerloss"`
	Mess           int32          `json:"mess"`
	Meals          int32          `json:"meals"`
	Disease        int32          `json:"disease"`
	Deadcommand    int32          `json:"deadcommand"`
	Arrested       int32          `json:"arrested"`
	Traitorswon    int32          `json:"traitorswon"`
	Roleswon       int32          `json:"roleswon"`
	Allarrested    int32          `json:"allarrested"`
	Opkilled       int32          `json:"opkilled"`
	Disc           int32          `json:"disc"`
	Nuked          int32          `json:"nuked"`
	Destranomaly   int32          `json:"destranomaly"`
	RecAntags      int32          `json:"rec_antags"`
	CrewEscaped    int32          `json:"crew_escaped"`
	CrewDead       int32          `json:"crew_dead"`
	CrewTotal      int32          `json:"crew_total"`
	CrewSurvived   int32          `json:"crew_survived"`
	Captain        pq.StringArray `json:"captain" gorm:"type:varchar(256)[]"`
	Powerbonus     int32          `json:"powerbonus"`
	Messbonus      int32          `json:"messbonus"`
	Deadaipenalty  int32          `json:"deadaipenalty"`
	Foodeaten      int32          `json:"foodeaten"`
	Clownabuse     int32          `json:"clownabuse"`
	Richestname    string         `json:"richestname"`
	Richestjob     string         `json:"richestjob"`
	Richestcash    float64        `json:"richestcash"`
	Richestkey     string         `json:"richestkey"`
	Dmgestname     string         `json:"dmgestname"`
	Dmgestjob      string         `json:"dmgestjob"`
	Dmgestdamage   int32          `json:"dmgestdamage"`
	Dmgestkey      string         `json:"dmgestkey"`
}

func (d *Score) ColumnsMigration(dx *gorm.DB) {
	dx.Migrator().AlterColumn(&d, "Richestname")
	dx.Migrator().AlterColumn(&d, "Richestjob")
	dx.Migrator().AlterColumn(&d, "Richestkey")
	dx.Migrator().AlterColumn(&d, "Richestcash")
	dx.Migrator().AlterColumn(&d, "Dmgestname")
	dx.Migrator().AlterColumn(&d, "Dmgestjob")
	dx.Migrator().AlterColumn(&d, "Dmgestkey")
}

type CommunicationLogs struct {
	ID      int32  `gorm:"uniqueIndex"`
	RootID  int32  `gorm:"index"`
	Time    string `json:"time" gorm:"size:256"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author" gorm:"size:256"`
	Type    string `json:"type" gorm:"size:256"`
}

func (d *CommunicationLogs) ColumnsMigration(dx *gorm.DB) {
	dx.Migrator().AlterColumn(&d, "Title")
}

type Damage struct {
	ID       int32   `gorm:"uniqueIndex"`
	DeathsID int32   `gorm:"uniqueIndex"`
	Brute    float64 `json:"BRUTE"`
	Fire     float64 `json:"FIRE"`
	Toxin    float64 `json:"TOXIN"`
	Oxy      float64 `json:"OXY"`
	Clone    float64 `json:"CLONE"`
	Brain    float64 `json:"BRAIN"`
}

type Deaths struct {
	ID               int32  `gorm:"uniqueIndex"`
	RootID           int32  `gorm:"index"`
	Name             string `json:"name" gorm:"size:256"`
	MobType          string `json:"mob_type" gorm:"size:256"`
	AssignedRole     string `json:"assigned_role" gorm:"size:256"`
	SpecialRole      string `json:"special_role" gorm:"size:256"`
	Damage           Damage `json:"damage"`
	RealName         string `json:"real_name" gorm:"size:256"`
	MindName         string `json:"mind_name" gorm:"size:256"`
	DeathX           int32  `json:"death_x"`
	DeathY           int32  `json:"death_y"`
	DeathZ           int32  `json:"death_z"`
	TimeOfDeath      string `json:"time_of_death" gorm:"size:128"`
	FromSuicide      int32  `json:"from_suicide"`
	LastAttackerName string `json:"last_attacker_name" gorm:"size:256"`
	LastPhrase       string `json:"last_phrase" gorm:"size:1024"`
	LastExaminedName string `json:"last_examined_name" gorm:"size:1024"`
}

func (d *Deaths) ColumnsMigration(dx *gorm.DB) {
	dx.Migrator().AlterColumn(&d, "TimeOfDeath")
}

type Explosions struct {
	ID               int32  `gorm:"uniqueIndex"`
	RootID           int32  `gorm:"index"`
	EpicenterX       int32  `json:"epicenter_x"`
	EpicenterY       int32  `json:"epicenter_y"`
	EpicenterZ       int32  `json:"epicenter_z"`
	DevastationRange int32  `json:"devastation_range"`
	HeavyImpactRange int32  `json:"heavy_impact_range"`
	LightImpactRange int32  `json:"light_impact_range"`
	FlashRange       int32  `json:"flash_range"`
	OccurredTime     string `json:"occurred_time"`
}

type EMPExplosions struct {
	ID               int32  `gorm:"uniqueIndex"`
	RootID           int32  `gorm:"index"`
	EpicenterX       int32  `json:"epicenter_x"`
	EpicenterY       int32  `json:"epicenter_y"`
	EpicenterZ       int32  `json:"epicenter_z"`
	DevastationRange int32  `json:"devastation_range"`
	HeavyImpactRange int32  `json:"heavy_range"`
	LightImpactRange int32  `json:"light_range"`
	OccurredTime     string `json:"occurred_time"`
}

type ManifestEntries struct {
	ID           int32          `gorm:"uniqueIndex"`
	RootID       int32          `gorm:"index"`
	Name         string         `json:"name" gorm:"size:256"`
	Species      string         `json:"species" gorm:"size:256"`
	Gender       string         `json:"gender" gorm:"size:256"`
	Flavor       string         `json:"flavor"`
	Age          uint           `json:"age"`
	AssignedRole string         `json:"assigned_role" gorm:"size:256"`
	SpecialRole  string         `json:"special_role" gorm:"size:256"`
	AntagRoles   pq.StringArray `json:"antag_roles" gorm:"type:varchar(256)[]"`
}

type LeaveStats struct {
	ID           int32          `gorm:"uniqueIndex"`
	RootID       int32          `gorm:"index"`
	Name         string         `json:"name" gorm:"size:256"`
	StartTime    string         `json:"start_time" gorm:"size:256"`
	AssignedRole string         `json:"assigned_role" gorm:"size:256"`
	SpecialRole  string         `json:"special_role" gorm:"size:256"`
	AntagRoles   pq.StringArray `json:"antag_roles" gorm:"type:varchar(256)[]"`
	LeaveType    string         `json:"leave_type" gorm:"size:256"`
	LeaveTime    string         `json:"leave_time" gorm:"size:256"`
}
