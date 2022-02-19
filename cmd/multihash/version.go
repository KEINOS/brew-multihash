package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/KEINOS/go-utiles/util"
)

// ----------------------------------------------------------------------------
//  Types
// ----------------------------------------------------------------------------

// define JSON data struct. Which is the structure of version.json from the
// upstream repo.
type verData struct {
	Version string `json:"version"`
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

func getVerJSON() string {
	var (
		verJSON verData // JSON data struct.
		verShow string
	)

	// Parse versionJSON to verData struct.
	err := json.Unmarshal(versionJSON, &verJSON)
	if err == nil && verJSON.Version != "" {
		verShow = verJSON.Version
	}

	return verShow
}

func printVer() {
	if isVerShort {
		printVerShort()

		return
	}

	fmt.Println("Application:")
	fmt.Printf("  %s\t", util.GetNameBin())
	printVerShort()

	printModules()
}

// priorities are: version set via ldflags > version in embedded version.json, or "(undefined)".
func printVerShort() {
	var (
		cmtShow = commit  // Commit ID set via ldflags.
		verShow = version // Version set via ldflags
	)

	// get version from embedded json.
	if verShow == "" {
		verShow = getVerJSON()
	}

	if verShow == "" {
		verShow = "v(undefined)"
	}

	// get commit ID via build flags.
	if cmtShow != "" {
		cmtShow = fmt.Sprintf("-%s", cmtShow)
	}

	fmt.Println(uniformVersion(verShow) + cmtShow)
}

func sortModules(mods []map[string]string) (sorted []map[string]string, maxLen int) {
	maxLen = 0

	sort.Slice(mods, func(i, j int) bool {
		if len(mods[i]["name"]) > maxLen {
			maxLen = len(mods[i]["name"])
		}

		return mods[i]["name"] < mods[j]["name"]
	})

	return mods, maxLen
}

func printModules() {
	listMods, maxLen := sortModules(util.GetMods())
	padding := strings.Repeat(" ", maxLen+1)

	fmt.Println("Modules:")

	for _, modInfo := range listMods {
		verRaw := modInfo["version"]
		verUni := uniformVersion(verRaw)
		name := modInfo["name"] + padding

		// Optional version info
		verOpt := ""
		if verRaw != verUni {
			verOpt = fmt.Sprintf("\t(%s)", verRaw)
		}

		fmt.Printf(
			"  %s\t%s\t%s%s\n",
			name[0:maxLen],
			verUni,
			modInfo["path"],
			verOpt,
		)
	}
}

func uniformVersion(verIn string) string {
	ver := strings.ToLower(verIn)

	pv, err := util.ParseVersion(ver)
	if err != nil {
		return verIn
	}

	return fmt.Sprintf("v%v.%v.%v", pv["major"], pv["minor"], pv["patch"])
}
