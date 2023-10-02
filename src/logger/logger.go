package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/logger"
)

type Logger struct {
	*logger.Logger
}

var ColorFuncMap = map[string]*color.Color{
	"red":     color.New(color.FgRed),
	"green":   color.New(color.FgGreen),
	"yellow":  color.New(color.FgYellow),
	"blue":    color.New(color.FgBlue),
	"magenta": color.New(color.FgMagenta),
	"cyan":    color.New(color.FgCyan),
	"white":   color.New(color.FgWhite),

	"bred":     color.New(color.FgRed, color.Bold),
	"bgreen":   color.New(color.FgGreen, color.Bold),
	"byellow":  color.New(color.FgYellow, color.Bold),
	"bblue":    color.New(color.FgBlue, color.Bold),
	"bmagenta": color.New(color.FgMagenta, color.Bold),
	"bcyan":    color.New(color.FgCyan, color.Bold),
	"bwhite":   color.New(color.FgWhite, color.Bold),

	"ured":     color.New(color.FgRed, color.Underline),
	"ugreen":   color.New(color.FgGreen, color.Underline),
	"uyellow":  color.New(color.FgYellow, color.Underline),
	"ublue":    color.New(color.FgBlue, color.Underline),
	"umagenta": color.New(color.FgMagenta, color.Underline),
	"ucyan":    color.New(color.FgCyan, color.Underline),
	"uwhite":   color.New(color.FgWhite, color.Underline),

	"bured":     color.New(color.FgRed, color.Bold, color.Underline),
	"bugreen":   color.New(color.FgGreen, color.Bold, color.Underline),
	"buyellow":  color.New(color.FgYellow, color.Bold, color.Underline),
	"bublue":    color.New(color.FgBlue, color.Bold, color.Underline),
	"bumagenta": color.New(color.FgMagenta, color.Bold, color.Underline),
	"bucyan":    color.New(color.FgCyan, color.Bold, color.Underline),
	"buwhite":   color.New(color.FgWhite, color.Bold, color.Underline),

	"ired":     color.New(color.FgRed, color.Italic),
	"igreen":   color.New(color.FgGreen, color.Italic),
	"iyellow":  color.New(color.FgYellow, color.Italic),
	"iblue":    color.New(color.FgBlue, color.Italic),
	"imagenta": color.New(color.FgMagenta, color.Italic),
	"icyan":    color.New(color.FgCyan, color.Italic),
	"iwhite":   color.New(color.FgWhite, color.Italic),

	"bired":     color.New(color.FgRed, color.Bold, color.Italic),
	"bigreen":   color.New(color.FgGreen, color.Bold, color.Italic),
	"biyellow":  color.New(color.FgYellow, color.Bold, color.Italic),
	"biblue":    color.New(color.FgBlue, color.Bold, color.Italic),
	"bimagenta": color.New(color.FgMagenta, color.Bold, color.Italic),
	"bicyan":    color.New(color.FgCyan, color.Bold, color.Italic),
	"biwhite":   color.New(color.FgWhite, color.Bold, color.Italic),

	"ibured":     color.New(color.FgRed, color.Italic, color.Underline, color.Bold),
	"ibugreen":   color.New(color.FgGreen, color.Italic, color.Underline, color.Bold),
	"ibuyellow":  color.New(color.FgYellow, color.Italic, color.Underline, color.Bold),
	"ibublue":    color.New(color.FgBlue, color.Italic, color.Underline, color.Bold),
	"ibumagenta": color.New(color.FgMagenta, color.Italic, color.Underline, color.Bold),
	"ibucyan":    color.New(color.FgCyan, color.Italic, color.Underline, color.Bold),
	"ibuwhite":   color.New(color.FgWhite, color.Italic, color.Underline, color.Bold),
}

var (
	RedSprint     = CreateColoredSprint("red")
	YellowSprint  = CreateColoredSprint("yellow")
	GreenSprint   = CreateColoredSprint("green")
	BlueSprint    = CreateColoredSprint("blue")
	MagentaSprint = CreateColoredSprint("magenta")
	CyanSprint    = CreateColoredSprint("cyan")
	WhiteSprint   = CreateColoredSprint("white")

	RedPrintln     = CreateColoredPrintln("red")
	YellowPrintln  = CreateColoredPrintln("yellow")
	GreenPrintln   = CreateColoredPrintln("green")
	BluePrintln    = CreateColoredPrintln("blue")
	MagentaPrintln = CreateColoredPrintln("magenta")
	CyanPrintln    = CreateColoredPrintln("cyan")
	WhitePrintln   = CreateColoredPrintln("white")

	BredSprint     = CreateColoredSprint("bred")
	ByellowSprint  = CreateColoredSprint("byellow")
	BgreenSprint   = CreateColoredSprint("bgreen")
	BblueSprint    = CreateColoredSprint("bblue")
	BmagentaSprint = CreateColoredSprint("bmagenta")
	BcyanSprint    = CreateColoredSprint("bcyan")
	BwhiteSprint   = CreateColoredSprint("bwhite")

	BredPrintln     = CreateColoredPrintln("bred")
	ByellowPrintln  = CreateColoredPrintln("byellow")
	BgreenPrintln   = CreateColoredPrintln("bgreen")
	BbluePrintln    = CreateColoredPrintln("bblue")
	BmagentaPrintln = CreateColoredPrintln("bmagenta")
	BcyanPrintln    = CreateColoredPrintln("bcyan")
	BwhitePrintln   = CreateColoredPrintln("bwhite")

	UredSprint     = CreateColoredSprint("ured")
	UyellowSprint  = CreateColoredSprint("uyellow")
	UgreenSprint   = CreateColoredSprint("ugreen")
	UblueSprint    = CreateColoredSprint("ublue")
	UmagentaSprint = CreateColoredSprint("umagenta")
	UcyanSprint    = CreateColoredSprint("ucyan")
	UwhiteSprint   = CreateColoredSprint("uwhite")

	UredPrintln     = CreateColoredPrintln("ured")
	UyellowPrintln  = CreateColoredPrintln("uyellow")
	UgreenPrintln   = CreateColoredPrintln("ugreen")
	UbluePrintln    = CreateColoredPrintln("ublue")
	UmagentaPrintln = CreateColoredPrintln("umagenta")
	UcyanPrintln    = CreateColoredPrintln("ucyan")
	UwhitePrintln   = CreateColoredPrintln("uwhite")

	BuredSprint     = CreateColoredSprint("bured")
	BuyellowSprint  = CreateColoredSprint("buyellow")
	BugreenSprint   = CreateColoredSprint("bugreen")
	BublueSprint    = CreateColoredSprint("bublue")
	BumagentaSprint = CreateColoredSprint("bumagenta")
	BucyanSprint    = CreateColoredSprint("bucyan")
	BuwhiteSprint   = CreateColoredSprint("buwhite")

	BuredPrintln     = CreateColoredPrintln("bured")
	BuyellowPrintln  = CreateColoredPrintln("buyellow")
	BugreenPrintln   = CreateColoredPrintln("bugreen")
	BubluePrintln    = CreateColoredPrintln("bublue")
	BumagentaPrintln = CreateColoredPrintln("bumagenta")
	BucyanPrintln    = CreateColoredPrintln("bucyan")
	BuwhitePrintln   = CreateColoredPrintln("buwhite")

	IredSprint     = CreateColoredSprint("ired")
	IyellowSprint  = CreateColoredSprint("iyellow")
	IgreenSprint   = CreateColoredSprint("igreen")
	IblueSprint    = CreateColoredSprint("iblue")
	ImagentaSprint = CreateColoredSprint("imagenta")
	IcyanSprint    = CreateColoredSprint("icyan")
	IwhiteSprint   = CreateColoredSprint("iwhite")

	IredPrintln     = CreateColoredPrintln("ired")
	IyellowPrintln  = CreateColoredPrintln("iyellow")
	IgreenPrintln   = CreateColoredPrintln("igreen")
	IbluePrintln    = CreateColoredPrintln("iblue")
	ImagentaPrintln = CreateColoredPrintln("imagenta")
	IcyanPrintln    = CreateColoredPrintln("icyan")
	IwhitePrintln   = CreateColoredPrintln("iwhite")

	BiredSprint     = CreateColoredSprint("bired")
	BiyellowSprint  = CreateColoredSprint("biyellow")
	BigreenSprint   = CreateColoredSprint("bigreen")
	BiblueSprint    = CreateColoredSprint("biblue")
	BimagentaSprint = CreateColoredSprint("bimagenta")
	BicyanSprint    = CreateColoredSprint("bicyan")
	BiwhiteSprint   = CreateColoredSprint("biwhite")

	BiredPrintln     = CreateColoredPrintln("bired")
	BiyellowPrintln  = CreateColoredPrintln("biyellow")
	BigreenPrintln   = CreateColoredPrintln("bigreen")
	BibluePrintln    = CreateColoredPrintln("biblue")
	BimagentaPrintln = CreateColoredPrintln("bimagenta")
	BicyanPrintln    = CreateColoredPrintln("bicyan")
	BiwhitePrintln   = CreateColoredPrintln("biwhite")

	IburedSprint     = CreateColoredSprint("ibured")
	IbuyellowSprint  = CreateColoredSprint("ibuyellow")
	IbugreenSprint   = CreateColoredSprint("ibugreen")
	IbublueSprint    = CreateColoredSprint("ibublue")
	IbumagentaSprint = CreateColoredSprint("ibumagenta")
	IbucyanSprint    = CreateColoredSprint("ibucyan")
	IbuwhiteSprint   = CreateColoredSprint("ibuwhite")

	IburedPrintln     = CreateColoredPrintln("ibured")
	IbuyellowPrintln  = CreateColoredPrintln("ibuyellow")
	IbugreenPrintln   = CreateColoredPrintln("ibugreen")
	IbubluePrintln    = CreateColoredPrintln("ibublue")
	IbumagentaPrintln = CreateColoredPrintln("ibumagenta")
	IbucyanPrintln    = CreateColoredPrintln("ibucyan")
	IbuwhitePrintln   = CreateColoredPrintln("ibuwhite")
)

func IsDebug() bool {
	return os.Getenv("BASHY_DEBUG") == "true"
}

func NewLogger() *Logger {
	var log *logger.Logger
	if IsDebug() {
		logger.SetFlags(1)
		logFile, err := os.OpenFile("logs/bashy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			RedPrintln("Failed to open log file")
			os.Exit(1)
		}
		log = logger.Init("Logger", true, false, logFile)
		ByellowPrintln("Debug mode enabled")
	} else {
		logger.SetFlags(0)
		log = logger.Init("Logger", false, false, os.Stdout)
		logger.SetFlags(0)
	}

	return &Logger{log}
}

func InfoLog(msg string) {
	l := NewLogger()
	if IsDebug() {
		l.Info(msg)
	} else {
		l.Info(GreenSprint(msg))
	}
}

func WarningLog(msg string) {
	l := NewLogger()
	if IsDebug() {
		l.Warning(msg)
	} else {
		l.Warning(YellowSprint(msg))
	}
}

func ErrorLog(msg string) {
	l := NewLogger()
	if IsDebug() {
		l.Error(msg)
	} else {
		l.Error(RedSprint(msg))
	}
}

func FatalLog(msg string) {
	l := NewLogger()
	if IsDebug() {
		l.Fatal(msg)
	} else {
		l.Fatal(RedSprint(msg))
	}
}

func Log(level string, msg string) {
	switch strings.ToLower(level) {
	case "info":
		InfoLog(msg)
	case "warning":
		WarningLog(msg)
	case "error":
		ErrorLog(msg)
	case "fatal":
		FatalLog(msg)
	default:
		ErrorLog("Invalid log level " + level + " for message " + msg)
	}
}

func CreateColoredSprint(colorName string) func(string) string {
	colorFunc, exists := ColorFuncMap[colorName]
	if !exists {
		return func(s string) string {
			return s
		}
	}

	return func(msg string) string {
		return colorFunc.Sprint(msg)
	}
}

func CreateColoredPrintln(colorName string) func(...interface{}) {
	colorFunc, exists := ColorFuncMap[colorName]
	if !exists {
		return func(a ...interface{}) {
			fmt.Println(a...)
		}
	}

	return func(a ...interface{}) {
		colorFunc.Println(a...)
	}
}

func JsonEncode(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		Log("error", "Error encoding json")
		return ""
	}

	return string(j)
}

func JsonDecode(v interface{}) interface{} {
	var result interface{}
	err := json.Unmarshal([]byte(v.(string)), &result)

	if err != nil {
		Log("error", "Error decoding json")
		return nil
	}

	return result
}
