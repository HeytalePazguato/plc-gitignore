package rules

// Siemens returns the Siemens TIA Portal rule set.
func Siemens() RuleSet {
	return RuleSet{
		Vendor:      "siemens",
		DisplayName: "Siemens TIA Portal",
		DetectGlobs: []string{"*.ap17", "*.ap18", "*.ap19", "*.ap20"},
		Sections: []Section{
			{
				Title:   "TIA Portal frame storage",
				Comment: "TIA's internal object-frame storage — machine-specific.",
				Patterns: []Pattern{
					{Glob: ".Siemens.Automation.Objectframe.FileStorage/"},
				},
			},
			{
				Title:   "User-specific files",
				Comment: "Per-user TIA settings and additional files.",
				Patterns: []Pattern{
					{Glob: "UserFiles/"},
					{Glob: "AdditionalFiles/"},
					{Glob: "Logs/"},
				},
			},
			{
				Title:   "Archive files",
				Comment: "TIA project archives — built from source, not committed.",
				Patterns: []Pattern{
					{Glob: "*.ap*_*"},
					{Glob: "*.zap*"},
				},
			},
			{
				Title:   "Compiled blocks",
				Comment: "Compiled OB/FB/FC output, regenerated on build.",
				Patterns: []Pattern{
					{Glob: "**/*.compiled"},
					{Glob: "**/Stations/"},
				},
			},
			{
				Title:   "Temporary files",
				Comment: "TIA scratch space.",
				Patterns: []Pattern{
					{Glob: "*.tmp"},
					{Glob: "TempFiles/"},
				},
			},
		},
		Attributes: []AttrSection{
			{
				Title:   "Binary files",
				Comment: "TIA project files are binary and cannot be merged.",
				Rules: []AttrRule{
					{Pattern: "*.ap17", Attrs: "binary"},
					{Pattern: "*.ap18", Attrs: "binary"},
					{Pattern: "*.ap19", Attrs: "binary"},
					{Pattern: "*.ap20", Attrs: "binary"},
					{Pattern: "*.zap*", Attrs: "binary"},
				},
			},
			{
				Title:   "Line endings",
				Comment: "Force LF for any text files (SCL, AWL exports).",
				Rules: []AttrRule{
					{Pattern: "*", Attrs: "text=auto eol=lf"},
					{Pattern: "*.scl", Attrs: "text eol=lf"},
					{Pattern: "*.awl", Attrs: "text eol=lf"},
				},
			},
		},
		HMISections: []Section{
			{
				Title:   "WinCC HMI",
				Comment: "WinCC Unified / Comfort HMI build artifacts.",
				Patterns: []Pattern{
					{Glob: "**/HmiTargets/"},
					{Glob: "**/RuntimeBinaries/"},
				},
			},
		},
		LFSAttrs: []AttrSection{
			{
				Title:   "Git LFS",
				Comment: "TIA archives can be hundreds of MB — track via LFS.",
				Rules: []AttrRule{
					{Pattern: "*.ap17", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.ap18", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.ap19", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.ap20", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
				},
			},
		},
		HookWarnings: []HookWarning{
			{Pattern: ".Siemens.Automation.Objectframe.FileStorage/", Reason: "TIA frame storage"},
			{Pattern: "UserFiles/", Reason: "user-specific files"},
			{Pattern: "Logs/", Reason: "log directory"},
		},
	}
}
