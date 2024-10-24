package settings

var OP = []string{
	"FETCHING",
	"PARSING",
	"RETRY FETCHING",
	"RETRY PARSING",
	"PRINTING",
	"FOUND",
	"NOT FOUND",
	"CACHED",
	"CANCELLED",
	"DELETED",
}

var DICTS = []string{"CAMBRIDGE", "WEBSTER", "SOHA"}

const (
	CAMBRIDGE = iota
	WEBSTER
	SOHA
)
