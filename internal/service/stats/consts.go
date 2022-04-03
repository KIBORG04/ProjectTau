package stats

const (
	ServerAlphaAddress = "game.taucetistation.org:2506"
	ServerBetaAddress  = "game.taucetistation.org:2507"
	ServerGammaAddress = "game.taucetistation.org:2508"
)

const (
	Alpha = "Alpha"
	Beta  = "Beta"
	Gamma = "Gamma"
)

const (
	ObjectiveWIN   = "SUCCESS"
	ObjectiveHALF  = "HALF"
	ObjectiveLFAIL = "FAIL"
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
