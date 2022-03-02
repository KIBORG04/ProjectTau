package domain

type Factions struct {
	Name         string       `json:"name"`
	ID           string       `json:"id"`
	MinorVictory int          `json:"minor_victory"`
	Objectives   []Objectives `json:"objectives"`
	Members      []Role       `json:"members"`
	Type         string       `json:"type"`
	Victory      int          `json:"victory"`
	CultInfo     CultInfo     `json:"cult_info"`
}

type Role struct {
	Name             string       `json:"name"`
	ID               string       `json:"id"`
	IsRoundstartRole int          `json:"is_roundstart_role"`
	Objectives       []Objectives `json:"objectives"`
	Type             string       `json:"type"`
	Victory          int          `json:"victory"`
	FactionID        string       `json:"faction_id"`
	MindName         string       `json:"mind_name"`
	MindCkey         string       `json:"mind_ckey"`
	UplinkInfo       UplinkInfo   `json:"uplink_info"`
}

type Objectives struct {
	Owner              string `json:"owner"`
	ExplanationText    string `json:"explanation_text"`
	Completed          string `json:"completed"`
	TargetName         string `json:"target_name"`
	Type               string `json:"type"`
	TargetAssignedRole string `json:"target_assigned_role"`
	TargetSpecialRole  string `json:"target_special_role"`
}

type UplinkPurchases struct {
	Cost       int    `json:"cost"`
	Bundlename string `json:"bundlename"`
	ItemType   string `json:"item_type"`
}

type UplinkInfo struct {
	TotalTC         int               `json:"total_TC"`
	SpentTC         int               `json:"spent_TC"`
	UplinkPurchases []UplinkPurchases `json:"uplink_purchases"`
}

type Aspects struct {
	Mortem     int `json:"Mortem"`
	Progressus int `json:"Progressus"`
	Fames      int `json:"Fames"`
	Telum      int `json:"Telum"`
	Metallum   int `json:"Metallum"`
	Partum     int `json:"Partum"`
	Cruciatu   int `json:"Cruciatu"`
	Salutis    int `json:"Salutis"`
	Spiritus   int `json:"Spiritus"`
	Arsus      int `json:"Arsus"`
	Chaos      int `json:"Chaos"`
	Rabidus    int `json:"Rabidus"`
	Absentia   int `json:"Absentia"`
	Obscurum   int `json:"Obscurum"`
	Lux        int `json:"Lux"`
	Lucrum     int `json:"Lucrum"`
	Turbam     int `json:"Turbam"`
}

type RitenameByCount struct {
	Deathalarm      int `json:"Ангел-хранитель"`
	Sacrifice       int `json:"Жертвоприношение"`
	Convert         int `json:"Обращение"`
	Emp             int `json:"ЭМИ"`
	DrainTorture    int `json:"Высасывание Жизни"`
	RaiseTorture    int `json:"Воскрешение"`
	CreateSlave     int `json:"Создание Гомункула"`
	SummonAcolyt    int `json:"Призыв Аколита"`
	Brainswap       int `json:"Обмен Разумов"`
	GiveForcearmor  int `json:"Создание Силовой Ауры"`
	UpgradeTome     int `json:"Улучшение Тома"`
	ImposeBlind     int `json:"Наложить Ослепление"`
	ImposeDeaf      int `json:"Наложить Глухоту"`
	ImposeStun      int `json:"Наложить Оглушение"`
	Communicate     int `json:"Общение"`
	Talisman        int `json:"Призыв Талисмана"`
	Soulstone       int `json:"Призыв Камня Душ"`
	Constructshell  int `json:"Призыв Оболочки"`
	Narsie          int `json:"Призыв Нар-Си"`
	CultPortal      int `json:"Призыв Портала"`
	MakeSkeleton    int `json:"Скелетофикация"`
	Synthconversion int `json:"Синтетическое Возвышение"`
	Freesacrifice   int `json:"Добровольное Жертвоприношение"`
	Clownconversion int `json:"Клоунконверсия"`
	Invite          int `json:"Божественное Приглашение"`
	Charge          int `json:"Беспроводная Зарядка"`
	Food            int `json:"Создание Еды"`
	Pray            int `json:"Молитва"`
	Honk            int `json:"Клоунский Крик"`
	Animation       int `json:"Анимация"`
	Spook           int `json:"Испуг"`
	Illuminate      int `json:"Озарение"`
	ReviveAnimal    int `json:"Возрождение Животного"`
	Banana          int `json:"Атомная Реконструкция Молекулярной Решётки Целого Благословлённого Банана."`
	BananaOre       int `json:"Обогащение Молекул Кислорода Атомами Банана"`
	CallAnimal      int `json:"Призыв Животного"`
	CreateSword     int `json:"Создание Меча"`
	CreateTalisman  int `json:"Создание Талисмана"`
	Devaluation     int `json:"Девальвация"`
	Upgrade         int `json:"Улучшение"`
}

type CultInfo struct {
	Aspects            Aspects         `json:"aspects"`
	RitenameByCount    RitenameByCount `json:"ritename_by_count"`
	RealNumberMembers  uint            `json:"real_number_members"`
	CapturedAreas      uint            `json:"captured_areas"`
	EndFavor           float64         `json:"end_favor"`
	EndPiety           float64         `json:"end_piety"`
	RunesOnStation     uint            `json:"runes_on_station"`
	AnomaliesDestroyed uint            `json:"anomalies_destroyed"`
}
