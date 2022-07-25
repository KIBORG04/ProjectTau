package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/utils"
	"strconv"
	"strings"
)

func setCheckboxStates(c *gin.Context) {
	var serversForm ServerCheckbox
	err := c.ShouldBind(&serversForm)
	if err != nil {
		println(err)
		return
	}
	c.SetCookie("serverCheckboxes", strings.Join(serversForm.Checkboxes, "|"), 10000, "", "", true, true)
	c.Set("serverCheckboxes", serversForm)
}

func getCheckboxStates(c *gin.Context) map[string]string {
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

// db is configured
func getRootsByCheckboxes(db *gorm.DB, c *gin.Context) ([]domain.Root, []*domain.Root, []*domain.Root, []*domain.Root, []*domain.Root) {
	checkboxes := getCheckboxStates(c)

	var roots []domain.Root

	db.Find(&roots)

	var processRoots []*domain.Root

	var alphaRoots []*domain.Root
	var betaRoots []*domain.Root
	var gammaRoots []*domain.Root
	for _, rr := range roots {
		root := rr
		switch root.ServerAddress {
		case ServerAlphaAddress:
			alphaRoots = append(alphaRoots, &root)
			if checkboxes[Alpha] != "" {
				processRoots = append(processRoots, &root)
			}
		case ServerBetaAddress:
			betaRoots = append(betaRoots, &root)
			if checkboxes[Beta] != "" {
				processRoots = append(processRoots, &root)
			}
		case ServerGammaAddress:
			gammaRoots = append(gammaRoots, &root)
			if checkboxes[Gamma] != "" {
				processRoots = append(processRoots, &root)
			}
		}
	}

	return roots, processRoots, alphaRoots, betaRoots, gammaRoots
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
func (info InfoSlice) hasName(name string) (*StatInfo, bool) {
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

func isStationPlayer(assignment, name string) bool {
	return slices.Contains(stationPositions, assignment) && utils.IsDrone.FindString(name) == ""
}

type RoundTime struct {
	Hour uint
	Min  uint
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

func isRoundStartLeaver(stat domain.LeaveStats) bool {
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
