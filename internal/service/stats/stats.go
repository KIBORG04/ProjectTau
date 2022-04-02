package stats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"sort"
	"ssstatistics/internal/domain"
	"ssstatistics/internal/keymap"
	r "ssstatistics/internal/repository"
	"ssstatistics/internal/service/charts"
	"ssstatistics/internal/utils"
	"strings"
)

type ServerCheckbox struct {
	Checkboxes []string `form:"server[]"`
}

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
			Alpha: "",
			Beta:  "",
			Gamma: "",
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

func RootPOST(c *gin.Context) {
	setCheckboxStates(c)
	RootGET(c)
}

func RootGET(c *gin.Context) {
	var roots []domain.Root
	r.Database.Preload("Deaths").Find(&roots)

	var links domain.Link
	str := fmt.Sprintf("%%%d%%", roots[len(roots)-1].RoundID)
	r.Database.Where("link LIKE ?", str).First(&links)

	crewDeathsCount := make(keymap.MyMap[string, uint], 0)
	crewDeathsSum := 0

	roleDeathsCount := make(keymap.MyMap[string, uint], 0)
	roleDeathsSum := 0

	modesCount := make(keymap.MyMap[string, uint], 0)
	modesSum := 0

	checkboxStates := getCheckboxStates(c)

	var processRoots []*domain.Root

	var alphaRoots []*domain.Root
	var betaRoots []*domain.Root
	var gammaRoots []*domain.Root
	for _, rr := range roots {
		root := rr
		switch root.ServerAddress {
		case ServerAlphaAddress:
			alphaRoots = append(alphaRoots, &root)
			if checkboxStates[Alpha] != "" {
				processRoots = append(processRoots, &root)
			}
		case ServerBetaAddress:
			betaRoots = append(betaRoots, &root)
			if checkboxStates[Beta] != "" {
				processRoots = append(processRoots, &root)
			}
		case ServerGammaAddress:
			gammaRoots = append(gammaRoots, &root)
			if checkboxStates[Gamma] != "" {
				processRoots = append(processRoots, &root)
			}
		}
	}

	for _, root := range processRoots {
		modesCount = keymap.AddElem(modesCount, root.Mode, 1)
		modesSum++
		for _, death := range root.Deaths {
			if slices.Contains(stationPositions, death.AssignedRole) && utils.IsDrone.FindString(death.Name) == "" {
				crewDeathsCount = keymap.AddElem(crewDeathsCount, death.AssignedRole, 1)
				crewDeathsSum++
			}
			if death.SpecialRole != "" {
				roleDeathsCount = keymap.AddElem(roleDeathsCount, death.SpecialRole, 1)
				roleDeathsSum++
			}
		}
	}

	sort.Stable(sort.Reverse(modesCount))
	sort.Stable(sort.Reverse(crewDeathsCount))
	sort.Stable(sort.Reverse(roleDeathsCount))

	c.HTML(200, "index.html", gin.H{
		"totalRounds":      len(roots),
		"version":          roots[len(roots)-1].Version,
		"lastDate":         links.Date,
		"alphaRounds":      len(alphaRoots),
		"betaRounds":       len(betaRoots),
		"gammaRounds":      len(gammaRoots),
		"modesCount":       modesCount,
		"modesSum":         modesSum,
		"crewDeathsCount":  crewDeathsCount,
		"crewDeathsSum":    crewDeathsSum,
		"roleDeathsCount":  roleDeathsCount,
		"roleDeathsSum":    roleDeathsSum,
		"serverCheckboxes": checkboxStates,
	})
}

func Gamemode(c *gin.Context) {
	render := charts.Gamemode()

	c.HTML(200, "chart.html", gin.H{
		"charts": render,
	})
}

func Cult(c *gin.Context) {
	render := charts.Cult()

	c.HTML(200, "chart.html", gin.H{
		"charts": render,
	})
}
