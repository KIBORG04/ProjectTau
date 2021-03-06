package stats

const (
	ServerAlphaAddress = "game.taucetistation.org:2506"
	ServerBetaAddress  = "game.taucetistation.org:2507"
	ServerGammaAddress = "game.taucetistation.org:2508"
)

const CurrentStatistics = "2022-02-27"

const (
	Alpha = "Alpha"
	Beta  = "Beta"
	Gamma = "Gamma"
)

var ServerByAddress = map[string]string{
	Alpha: ServerAlphaAddress,
	Beta:  ServerBetaAddress,
	Gamma: ServerGammaAddress,
}

const (
	ObjectiveWIN   = "SUCCESS"
	ObjectiveHALF  = "HALF"
	ObjectiveLFAIL = "FAIL"
)

const (
	GhostedInCryo = "Ghosted in Cryopod"
	Cryo          = "Cryopod"
	Ghosted       = "Ghosted"
	Disconnected  = "Disconnected"
)

var (
	commandPositions = []string{
		"Captain",
		"Head of Personnel",
		"Head of Security",
		"Chief Engineer",
		"Research Director",
		"Chief Medical Officer",
	}

	engineeringPositions = []string{
		"Chief Engineer",
		"Station Engineer",
		"Atmospheric Technician",
		"Technical Assistant",
	}

	medicalPositions = []string{
		"Chief Medical Officer",
		"Medical Doctor",
		"Geneticist",
		"Psychiatrist",
		"Chemist",
		"Virologist",
		"Paramedic",
		"Medical Intern",
	}

	sciencePositions = []string{
		"Research Director",
		"Scientist",
		"Geneticist", //Part of both medical and science
		"Roboticist",
		"Xenobiologist",
		"Xenoarchaeologist",
		"Research Assistant",
	}

	civilianPositions = []string{
		"Head of Personnel",
		"Barber",
		"Bartender",
		"Botanist",
		"Chef",
		"Janitor",
		"Librarian",
		"Quartermaster",
		"Cargo Technician",
		"Shaft Miner",
		"Recycler",
		"Internal Affairs Agent",
		"Chaplain",
		"Test Subject",
		"Clown",
		"Mime",
	}

	securityPositions = []string{
		"Head of Security",
		"Warden",
		"Detective",
		"Security Officer",
		"Forensic Technician",
		"Security Cadet",
	}

	nonhumanPositions = []string{
		"AI",
		"Cyborg",
		"pAI",
	}

	stationPositions []string

	SoloRoles = []string{
		"TraitorChan",
		"Traitor",
		"Wizard",
		"Changeling",
		"Cortical Borer",
		"Space Ninja",
		// tau ceti.........
		"Shadowling",
		"Thrall",
	}

	TeamlRoles = []string{
		"Organized Crimes Department",
		"Cult Of Blood",
		"Revolution",
		"Syndicate Operatives",
		"Blob Conglomerate",
		"Abductor Team",
		"Alien Hivemind",
		"Vox Shoal",
		"Zobmies",
		"Families",
	}

	ShortModeName = map[string]string{
		"Organized Crimes Department": "OCD",
		"Cult Of Blood":               "Cult",
		"Blob Conglomerate":           "Blobs",
		"Abductor Team":               "Abductors",
		"Alien Hivemind":              "Aliens",
		"Syndicate Operatives":        "The Nuke",
		"fwafaw":                      "fwaf",
	}
)

func PopulatePositions() {
	stationPositions = append(stationPositions, commandPositions...)
	stationPositions = append(stationPositions, engineeringPositions...)
	stationPositions = append(stationPositions, medicalPositions...)
	stationPositions = append(stationPositions, sciencePositions...)
	stationPositions = append(stationPositions, civilianPositions...)
	stationPositions = append(stationPositions, securityPositions...)
	stationPositions = append(stationPositions, nonhumanPositions...)
}
