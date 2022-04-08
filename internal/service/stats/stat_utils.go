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

//type Test struct {
//	TheName  string
//	TheCount uint
//	Some     string
//	AndHZ    int
//}
//
//func (t Test) GetName() string {
//	return t.TheName
//}
//
//func (t Test) GetCountCount() uint {
//	return t.TheCount
//}

//
//type Test2 struct {
//	TheName    string
//	TheCount   uint
//	Afawfaw    string
//	AndfwafaHZ int
//}
//
//type Amogus struct {
//	Test2
//	fafkoaw string
//}
//
//func (t Test2) GetName() string {
//	return t.TheName
//}
//
//func (t Test2) GetCountCount() uint {
//	return t.TheCount
//}
//
//func Some() {
//	t2 := Test2{
//		Afawfaw:    Alpha,
//		AndfwafaHZ: 1,
//		TheCount:   1,
//		TheName:    "afwa",
//	}
//
//	t := Test{
//		Some:     Alpha,
//		AndHZ:    1,
//		TheCount: 1,
//		TheName:  "afwa",
//	}
//
//	a := Amogus{}
//
//	some := make(InfoSlice, 0)
//	some = append(some, t)
//	some = append(some, t2)
//	some = append(some, a)
//
//	ab := Some11(some)
//}
