package event

var columnNames = map[string]string{
	"Date":        "Date",
	"Day":         "Day",
	"Day Code":    "DayCode",
	"Time":        "Time",
	"Event Code":  "EventCode",
	"Event":       "EventName",
	"Type":        "Type",
	"RoundHeat":   "Session",
	"R/H":         "Session",
	"Prize Level": "Prizes",
	"Class":       "Class",
	"Style":       "Style",
	"Format":      "Format",
	"Duration":    "Duration",
	"Location":    "Location",
	"GM":          "GM",
	"Category":    "Category",
}

var validCategories = []string{
	"Century",
	"Legacy",
	"Trial",
	"Sponsor",
	"Vendor",
	"Admin",
	"Auction",
	"Demo",
	"Juniors",
	"Meeting",
	"Open Gaming",
	"Registration",
	"Seminar",
	"Service",
	"Shuttle",
	"Vendors",
}

var validClasses = []string{
	"A",
	"B",
	"C",
}

var validDayCodes = []string{
	"FFr",
	"FSa",
	"FSu",
	"Mo",
	"Tu",
	"We",
	"Th",
	"Fr",
	"Sa",
	"Su",
	"SMo",
}

var validDays = []string{
	"First Friday",
	"First Saturday",
	"First Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
	"Sunday",
	"Second Monday",
}

var validEventTypes = []string{
	"Tournament",
	"Auction",
	"Demo",
	"Juniors",
	"Meeting",
	"Open Gaming",
	"Registration",
	"Seminar",
	"Service",
	"Shuttle",
	"Vendors",
}

var validFormats = []string{
	"HE",
	"HMSE",
	"HMW-P",
	"HMW-T",
	"HWO",
	"Jr SE",
	"SE",
	"SEM",
	"SW",
	"SwEl",
}

var validMultipleSessionTypes = []string{
	"Demo",
	"Heat",
	"Round",
}

var validUniqueSessionTypes = []string{
	"Draft",
	"Mulligan",
	"Quarterfinal",
	"Semifinal",
	"Final",
}

var validStyles = []string{
	"Continuous",
	"Freeform",
	"Scheduled",
}
