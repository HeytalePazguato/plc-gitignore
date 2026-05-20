package rules

// Rockwell returns the Rockwell Studio 5000 Logix Designer rule set.
func Rockwell() RuleSet {
	return RuleSet{
		Vendor:      "rockwell",
		DisplayName: "Rockwell Studio 5000",
		DetectGlobs: []string{"*.ACD", "*.L5X", "*.L5K"},
		Sections: []Section{
			{
				Title:   "Backup files",
				Comment: "Studio 5000 auto-creates .BAK on save.",
				Patterns: []Pattern{
					{Glob: "*.BAK"},
					{Glob: "*.bak"},
				},
			},
			{
				Title:   "CRC / checksum files",
				Comment: "Generated, machine-specific.",
				Patterns: []Pattern{
					{Glob: "*.CRC"},
					{Glob: "*.crc"},
				},
			},
			{
				Title:   "Compiled runtime",
				Comment: "RSS is the compiled runtime export — built from the ACD.",
				Patterns: []Pattern{
					{Glob: "*.RSS"},
					{Glob: "*.rss"},
				},
			},
			{
				Title:   "Auto-save files",
				Comment: "Studio 5000 auto-save scratch.",
				Patterns: []Pattern{
					{Glob: "*.WRK"},
					{Glob: "*.SEM"},
					{Glob: "*.SEMOLD"},
					{Glob: "*.ACD.lock"},
				},
			},
			{
				Title:   "Project archives",
				Comment: "ACDs are binary; their archives shouldn't be diffed.",
				Patterns: []Pattern{
					{Glob: "*.ACD.bak"},
				},
			},
		},
		Attributes: []AttrSection{
			{
				Title:   "Binary files",
				Comment: "ACD is a proprietary binary format; L5X/L5K are XML/text.",
				Rules: []AttrRule{
					{Pattern: "*.ACD", Attrs: "binary"},
					{Pattern: "*.acd", Attrs: "binary"},
				},
			},
			{
				Title:   "Merge strategy",
				Comment: "L5X is XML — union merge reduces conflicts on parallel edits.",
				Rules: []AttrRule{
					{Pattern: "*.L5X", Attrs: "merge=union"},
					{Pattern: "*.l5x", Attrs: "merge=union"},
				},
			},
			{
				Title:   "Line endings",
				Comment: "Force LF for text exports.",
				Rules: []AttrRule{
					{Pattern: "*", Attrs: "text=auto eol=lf"},
					{Pattern: "*.L5X", Attrs: "text eol=lf"},
					{Pattern: "*.L5K", Attrs: "text eol=lf"},
				},
			},
		},
		HMISections: []Section{
			{
				Title:   "FactoryTalk View HMI",
				Comment: "FT View ME/SE build output.",
				Patterns: []Pattern{
					{Glob: "**/CacheRoot/"},
					{Glob: "**/HMICache/"},
					{Glob: "*.MER"},
				},
			},
		},
		LFSAttrs: []AttrSection{
			{
				Title:   "Git LFS",
				Comment: "ACD files are large binaries — track via LFS.",
				Rules: []AttrRule{
					{Pattern: "*.ACD", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.acd", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
				},
			},
		},
		HookWarnings: []HookWarning{
			{Pattern: "*.BAK", Reason: "backup file — auto-generated"},
			{Pattern: "*.CRC", Reason: "checksum — machine-specific"},
			{Pattern: "*.WRK", Reason: "auto-save scratch"},
		},
	}
}
