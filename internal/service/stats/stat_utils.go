package stats

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"ssstatistics/internal/domain"
	r "ssstatistics/internal/repository"
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

func getRootsByCheckboxes(preloads []string, checkboxes map[string]string) ([]domain.Root, []*domain.Root, []*domain.Root, []*domain.Root, []*domain.Root) {
	var roots []domain.Root

	query := r.Database
	for _, v := range preloads {
		query = query.Preload(v)
	}
	query.Find(&roots)

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
