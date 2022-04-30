package stats

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"ssstatistics/internal/domain"
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
func getRootsByCheckboxes(db *gorm.DB, checkboxes map[string]string) ([]domain.Root, []*domain.Root, []*domain.Root, []*domain.Root, []*domain.Root) {
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

type BaseInfo struct {
	Name  string
	Count uint
}

func (b BaseInfo) GetName() string {
	return b.Name
}

func (b BaseInfo) GetCount() uint {
	return b.Count
}
