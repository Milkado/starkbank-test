package common

import "runtime"

var Reset = "\033[0m"
var Red = "\033[38;5;196m"    // Pure Red
var Green = "\033[38;5;46m"   // Pure Green
var Yellow = "\033[38;5;226m" // Pure Yellow
var Blue = "\033[38;5;27m"    // Pure Blue
var Purple = "\033[38;5;164m" // Pure Purple
var Cyan = "\033[38;5;51m"    // Pure Cyan
var Gray = "\033[38;5;245m"   // Medium Gray
var White = "\033[38;5;255m"  // Pure White

func Init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}
