package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"regexp"
	"ssstatistics/internal/domain"
	"strconv"
	"strings"
	"time"
)

// Utility regex
var (
	ckeyRegex = regexp.MustCompile(`[';()\\"&^?:$#№@_\s%]`)
	IsDrone   = regexp.MustCompile(`maintenance drone \(\d+\)`)
	IsMobName = regexp.MustCompile(`\w+ \(\d+\)`)
)

type ServerCheckbox struct {
	Checkboxes []string `form:"server[]"`
}

func SetCheckboxStates(c *gin.Context) {
	var serversForm ServerCheckbox
	err := c.ShouldBind(&serversForm)
	if err != nil {
		println(err)
		return
	}
	c.SetCookie("serverCheckboxes", strings.Join(serversForm.Checkboxes, "|"), 10000, "", "", true, true)
	c.Set("serverCheckboxes", serversForm)
}

func GetCheckboxStates(c *gin.Context) map[string]string {
	bool2Checked := func(b bool) string {
		if b {
			return "checked"
		}
		return ""
	}

	// If POST request
	formInterface, exist := c.Get("serverCheckboxes")
	if exist {
		serversForm := formInterface.(ServerCheckbox)
		return map[string]string{
			Alpha: bool2Checked(slices.Contains(serversForm.Checkboxes, Alpha)),
			Beta:  bool2Checked(slices.Contains(serversForm.Checkboxes, Beta)),
			Gamma: bool2Checked(slices.Contains(serversForm.Checkboxes, Gamma)),
		}
	}

	// If GET request
	cookie, err := c.Cookie("serverCheckboxes")
	if err != nil {
		return map[string]string{
			Alpha: "checked",
			Beta:  "checked",
			Gamma: "checked",
		}
	}

	checkboxes := strings.Split(cookie, "|")
	checkboxesStates := map[string]string{
		Alpha: bool2Checked(slices.Contains(checkboxes, Alpha)),
		Beta:  bool2Checked(slices.Contains(checkboxes, Beta)),
		Gamma: bool2Checked(slices.Contains(checkboxes, Gamma)),
	}
	return checkboxesStates
}

func ApplyDBQueryByDate(db *gorm.DB, c *gin.Context) {
	startDate, endDate, err := GetValidDates(c)
	if err != nil {
		return
	}
	db.Where("date BETWEEN ? AND ?", startDate, endDate)
}

func GetValidDates(c *gin.Context) (string, string, error) {
	startDate := c.DefaultPostForm("date_start", "2022-02-27")
	endDate := c.DefaultPostForm("date_end", time.Now().Format("2006-01-02"))
	if startDate != "" && endDate != "" {
		startStatisticsDateTime, _ := time.Parse("2006-01-02", CurrentStatisticsDate)
		startDateTime, _ := time.Parse("2006-01-02", startDate)
		if startDateTime.After(startStatisticsDateTime) || startDateTime.Equal(startStatisticsDateTime) {
			return startDate, endDate, nil
		}
	}
	return startDate, endDate, fmt.Errorf("not valid date or dates")
}

func GetChosenServers(c *gin.Context) (string, string, string) {
	var checkboxes = GetCheckboxStates(c)
	// cringe and cringe orm
	checkboxesKeys := make([]string, 0, len(checkboxes))
	for k, v := range checkboxes {
		if v != "" {
			checkboxesKeys = append(checkboxesKeys, k)
		} else {
			checkboxesKeys = append(checkboxesKeys, "")
		}
	}
	return ServerByAddress[checkboxesKeys[0]], ServerByAddress[checkboxesKeys[1]], ServerByAddress[checkboxesKeys[2]]
}

func ApplyDBQueryByServers(db *gorm.DB, c *gin.Context) {
	s1, s2, s3 := GetChosenServers(c)
	db.Where("server_address = ? OR server_address = ? OR server_address = ?", s1, s2, s3)
}

// GetRoots db is configured
func GetRoots(db *gorm.DB, c *gin.Context) []*domain.Root {
	checkboxes := GetCheckboxStates(c)

	var roots []*domain.Root

	ApplyDBQueryByDate(db, c)
	ApplyDBQueryByServers(db, c)

	db.Omit("CompletionHTML").
		Find(&roots)

	var processRoots []*domain.Root
	for _, rr := range roots {
		root := rr
		switch root.ServerAddress {
		case ServerAlphaAddress:
			if checkboxes[Alpha] != "" {
				processRoots = append(processRoots, root)
			}
		case ServerBetaAddress:
			if checkboxes[Beta] != "" {
				processRoots = append(processRoots, root)
			}
		case ServerGammaAddress:
			if checkboxes[Gamma] != "" {
				processRoots = append(processRoots, root)
			}
		}
	}

	return processRoots
}

type InfoSlice []StatInfo

func (info InfoSlice) Len() int {
	return len(info)
}
func (info InfoSlice) Less(i, j int) bool {
	return info[i].GetCount() < info[j].GetCount()
}
func (info InfoSlice) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}
func (info InfoSlice) HasName(name string) (*StatInfo, bool) {
	for i := 0; i < len(info); i++ {
		if info[i].GetName() == name {
			return &info[i], true
		}
	}
	return nil, false
}

type StatInfo interface {
	GetName() string
	GetCount() uint
}

func IsStationPlayer(assignment, name string) bool {
	return slices.Contains(StationPositions, assignment) && IsDrone.FindString(name) == ""
}

type RoundTime struct {
	Hour uint
	Min  uint
}

func (r RoundTime) ToSeconds() uint {
	return r.Hour*3600 + r.Min*60
}

func ParseRoundTime(time string) (RoundTime, error) {
	if time == "" {
		return RoundTime{}, fmt.Errorf("empty time")
	}
	strs := strings.Split(time, ":")

	hour, err := strconv.Atoi(strs[0])
	if err != nil {
		return RoundTime{}, err
	}
	min, err := strconv.Atoi(strs[1])
	if err != nil {
		return RoundTime{}, err
	}
	return RoundTime{Hour: uint(hour), Min: uint(min)}, nil
}

func IsRoundStartLeaver(stat domain.LeaveStats) bool {
	if stat.LeaveType == "" {
		return false
	}

	roundTime, err := ParseRoundTime(stat.LeaveTime)
	if err != nil {
		return false
	}
	if stat.LeaveType == Cryo && roundTime.Min < 15 {
		return false
	}
	if 5 < roundTime.Min && roundTime.Min < 30 {
		return true
	}
	if stat.LeaveType == Cryo && roundTime.Min < 45 { // body is in cryo for 15 minutes
		return true
	}
	return false
}

// Ckey byond functions which remove all spaces and not word symbols
func Ckey(str string) string {
	return strings.ToLower(ckeyRegex.ReplaceAllString(str, ""))
}

type VictoryTypes int

const (
	FactionVictory VictoryTypes = iota
	RoleVictory
	ErrorVictory
)

// GetAntagonistWinType Реальную победу в режиме сложно определить, эта функция должна облегчить задачу
func GetAntagonistWinType(roleName, factionName string) (VictoryTypes, error) {
	if slices.Contains(TeamlRoles, factionName) {
		return FactionVictory, nil
	} else if slices.Contains(SoloRoles, roleName) {
		return RoleVictory, nil
	}
	return ErrorVictory, fmt.Errorf("role or faction not allowed")
}

type Player struct {
	Ckey string `form:"ckey"`
}

func GetValidatePlayer(c *gin.Context) (*Player, error) {
	var player Player
	err := c.BindQuery(&player)
	if err != nil {
		return nil, err
	}
	if player.Ckey == "" {
		return nil, fmt.Errorf("ckey not entered in query")
	}
	player.Ckey = Ckey(player.Ckey)
	return &player, nil
}

func NormalizeByondBase64(str string) string {
	re := regexp.MustCompile(`'data:image/png;base64,[^'><]*`)
	newStr := re.ReplaceAllStringFunc(str, func(base64Str string) string {
		lastChar := base64Str[len(base64Str)-1]
		if lastChar == '=' || (len(base64Str)%4) == 0 {
			return base64Str
		}
		return base64Str[:len(base64Str)-1]
	})
	return newStr
}
