package rules

// BR returns the B&R Automation Studio rule set.
func BR() RuleSet {
	return RuleSet{
		Vendor:      "br",
		DisplayName: "B&R Automation Studio",
		DetectGlobs: []string{"*.apj", "Physical.pkg", "Logical.pkg"},
		Sections: []Section{
			{
				Title:   "Build artifacts",
				Comment: "Binaries and intermediate output, regenerated on every build.",
				Patterns: []Pattern{
					{Glob: "Binaries/"},
					{Glob: "Temp/"},
					{Glob: "Diagnosis/"},
					{Glob: "**/Cpu/Cache/"},
				},
			},
			{
				Title:   "Lock files",
				Comment: "B&R uses .isopen as a project-locked-by-someone marker.",
				Patterns: []Pattern{
					{Glob: "*.isopen"},
				},
			},
			{
				Title:   "User settings",
				Comment: "Per-user workspace state.",
				Patterns: []Pattern{
					{Glob: "LastUser.set"},
					{Glob: "User.set"},
					{Glob: "*.set.bak"},
				},
			},
			{
				Title:   "CPU/hardware build cache",
				Comment: "Generated from Physical.pkg + Logical.pkg on build.",
				Patterns: []Pattern{
					{Glob: "**/RUC.br"},
					{Glob: "**/*.br"},
					{Glob: "**/SG4/"},
					{Glob: "**/X86/"},
				},
			},
		},
		Attributes: []AttrSection{
			{
				Title:   "Binary files",
				Comment: "B&R compiled output is binary.",
				Rules: []AttrRule{
					{Pattern: "*.br", Attrs: "binary"},
					{Pattern: "*.zip", Attrs: "binary"},
				},
			},
			{
				Title:   "Line endings",
				Comment: "Force LF for source text.",
				Rules: []AttrRule{
					{Pattern: "*", Attrs: "text=auto eol=lf"},
					{Pattern: "*.st", Attrs: "text eol=lf"},
				},
			},
		},
		HMISections: []Section{
			{
				Title:   "B&R mappView HMI",
				Comment: "mappView build output.",
				Patterns: []Pattern{
					{Glob: "**/mappView/Output/"},
					{Glob: "**/mappView/Tmp/"},
				},
			},
		},
		LFSAttrs: []AttrSection{
			{
				Title:   "Git LFS",
				Comment: "Large media — track via LFS.",
				Rules: []AttrRule{
					{Pattern: "*.png", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.bin", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
				},
			},
		},
		HookWarnings: []HookWarning{
			{Pattern: "Binaries/", Reason: "build artifact"},
			{Pattern: "Temp/", Reason: "temp directory"},
			{Pattern: "*.isopen", Reason: "lock file — user-specific"},
		},
	}
}
