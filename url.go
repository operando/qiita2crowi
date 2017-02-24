package main

import (
	"strings"
)

var urlSafe = strings.NewReplacer(
	`^`, `＾`, // for Crowi's regexp
	`$`, `＄`,
	`*`, `＊`,
	`%`, `％`, // query
	`?`, `？`,
	`/`, `／`, // Prevent unexpected stratification
)
