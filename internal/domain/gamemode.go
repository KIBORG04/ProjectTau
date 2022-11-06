package domain

type Factions struct {
	RootID            int32               `gorm:"index"`
	ID                int32               `gorm:"uniqueIndex"`
	Name              string              `json:"name" gorm:"size:256"`
	FactionName       string              `json:"id" gorm:"size:256"`
	MinorVictory      int32               `json:"minor_victory"`
	FactionObjectives []FactionObjectives `json:"objectives" gorm:"foreignKey:OwnerID"`
	Members           []Role              `json:"members" gorm:"foreignKey:OwnerID"`
	Type              string              `json:"type" gorm:"size:256"`
	Victory           int32               `json:"victory"`
	CultInfo          CultInfo            `json:"cult_info"`
}

type Role struct {
	OwnerID          int32            `gorm:"index"`
	ID               int32            `gorm:"uniqueIndex"`
	Name             string           `json:"name" gorm:"size:256"`
	RoleName         string           `json:"id" gorm:"size:256"`
	IsRoundstartRole int32            `json:"is_roundstart_role"`
	RoleObjectives   []RoleObjectives `json:"objectives" gorm:"foreignKey:OwnerID"`
	Type             string           `json:"type"`
	Victory          int32            `json:"victory"`
	FactionName      string           `json:"faction_id" gorm:"size:256"`
	MindName         string           `json:"mind_name" gorm:"size:256"`
	MindCkey         string           `json:"mind_ckey" gorm:"size:256"`
	UplinkInfo       UplinkInfo       `json:"uplink_info" gorm:"foreignKey:RoleID"`
	ChangelingInfo   ChangelingInfo   `json:"changeling_info" gorm:"foreignKey:RoleID"`
}

type FactionObjectives Objectives
type RoleObjectives Objectives
type Objectives struct {
	ID                 int32  `gorm:"uniqueIndex"`
	OwnerID            int32  `gorm:"index"`
	Owner              string `json:"owner" gorm:"size:256"`
	ExplanationText    string `json:"explanation_text"`
	Completed          string `json:"completed" gorm:"size:16"`
	TargetName         string `json:"target_name" gorm:"size:256"`
	Type               string `json:"type" gorm:"size:256"`
	TargetAssignedRole string `json:"target_assigned_role" gorm:"size:256"`
	TargetSpecialRole  string `json:"target_special_role" gorm:"size:256"`
}

type UplinkInfo struct {
	ID              int32             `gorm:"uniqueIndex"`
	RoleID          int32             `gorm:"index"`
	TotalTC         int32             `json:"total_TC"`
	SpentTC         int32             `json:"spent_TC"`
	UplinkPurchases []UplinkPurchases `json:"uplink_purchases"`
}

type UplinkPurchases struct {
	ID           int32  `gorm:"uniqueIndex"`
	UplinkInfoID int32  `gorm:"index"`
	Cost         int32  `json:"cost"`
	Bundlename   string `json:"bundlename" gorm:"size:256"`
	ItemType     string `json:"item_type" gorm:"size:256"`
}

type ChangelingInfo struct {
	ID                 int32                `gorm:"uniqueIndex"`
	RoleID             int32                `gorm:"index"`
	VictimsNumber      int32                `json:"victims_number"`
	ChangelingPurchase []ChangelingPurchase `json:"changeling_purchase"`
}

type ChangelingPurchase struct {
	ID               int32  `gorm:"uniqueIndex"`
	ChangelingInfoID int32  `gorm:"index"`
	Cost             int32  `json:"cost"`
	PowerType        string `json:"power_type"`
	PowerName        string `json:"power_name" gorm:"size:256"`
}

type CultInfo struct {
	ID                 int32           `gorm:"uniqueIndex"`
	FactionsID         int32           `gorm:"index"`
	Aspects            Aspects         `json:"aspects"`
	RitenameByCount    RitenameByCount `json:"ritename_by_count"`
	RealNumberMembers  int32           `json:"real_number_members"`
	CapturedAreas      int32           `json:"captured_areas"`
	EndFavor           float64         `json:"end_favor"`
	EndPiety           float64         `json:"end_piety"`
	RunesOnStation     int32           `json:"runes_on_station"`
	AnomaliesDestroyed int32           `json:"anomalies_destroyed"`
}

type Aspects struct {
	ID         int32 `gorm:"uniqueIndex"`
	CultInfoID int32 `gorm:"index"`
	Mortem     int32 `json:"Mortem"`
	Cruciatu   int32 `json:"cruciatu"`
	Progressus int32 `json:"Progressus"`
	Fames      int32 `json:"Fames"`
	Telum      int32 `json:"Telum"`
	Metallum   int32 `json:"Metallum"`
	Partum     int32 `json:"Partum"`
	Salutis    int32 `json:"Salutis"`
	Spiritus   int32 `json:"Spiritus"`
	Arsus      int32 `json:"Arsus"`
	Chaos      int32 `json:"Chaos"`
	Rabidus    int32 `json:"Rabidus"`
	Obscurum   int32 `json:"Obscurum"`
	Lux        int32 `json:"Lux"`
	Lucrum     int32 `json:"Lucrum"`
	Turbam     int32 `json:"Turbam"`
}

type RitenameByCount struct {
	ID              int32 `gorm:"uniqueIndex"`
	CultInfoID      int32 `gorm:"index"`
	Deathalarm      int32 `json:"Ангел-хранитель"`
	Sacrifice       int32 `json:"Жертвоприношение"`
	Convert         int32 `json:"Обращение"`
	Emp             int32 `json:"ЭМИ"`
	Draint32orture  int32 `json:"Высасывание Жизни"`
	RaiseTorture    int32 `json:"Воскрешение"`
	CreateSlave     int32 `json:"Создание Гомункула"`
	SummonAcolyt    int32 `json:"Призыв Аколита"`
	Brainswap       int32 `json:"Обмен Разумов"`
	GiveForcearmor  int32 `json:"Создание Силовой Ауры"`
	UpgradeTome     int32 `json:"Улучшение Тома"`
	ImposeBlind     int32 `json:"Наложить Ослепление"`
	ImposeDeaf      int32 `json:"Наложить Глухоту"`
	ImposeStun      int32 `json:"Наложить Оглушение"`
	Communicate     int32 `json:"Общение"`
	Talisman        int32 `json:"Призыв Талисмана"`
	Soulstone       int32 `json:"Призыв Камня Душ"`
	Constructshell  int32 `json:"Призыв Оболочки"`
	Narsie          int32 `json:"Призыв Нар-Си"`
	CultPortal      int32 `json:"Призыв Портала"`
	MakeSkeleton    int32 `json:"Скелетофикация"`
	Synthconversion int32 `json:"Синтетическое Возвышение"`
	Freesacrifice   int32 `json:"Добровольное Жертвоприношение"`
	Clownconversion int32 `json:"Клоунконверсия"`
	Invite          int32 `json:"Божественное Приглашение"`
	Charge          int32 `json:"Беспроводная Зарядка"`
	Food            int32 `json:"Создание Еды"`
	Pray            int32 `json:"Молитва"`
	Honk            int32 `json:"Клоунский Крик"`
	Animation       int32 `json:"Анимация"`
	Spook           int32 `json:"Испуг"`
	Illuminate      int32 `json:"Озарение"`
	ReviveAnimal    int32 `json:"Возрождение Животного"`
	Banana          int32 `json:"Атомная Реконструкция Молекулярной Решётки Целого Благословлённого Банана."`
	BananaOre       int32 `json:"Обогащение Молекул Кислорода Атомами Банана"`
	CallAnimal      int32 `json:"Призыв Животного"`
	CreateSword     int32 `json:"Создание Меча"`
	CreateTalisman  int32 `json:"Создание Талисмана"`
	Devaluation     int32 `json:"Девальвация"`
	Upgrade         int32 `json:"Улучшение"`
}
