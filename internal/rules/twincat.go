package rules

// TwinCAT returns the TwinCAT 3 rule set.
func TwinCAT() RuleSet {
	return RuleSet{
		Vendor:      "twincat",
		DisplayName: "TwinCAT 3",
		DetectGlobs: []string{"*.tsproj", "*.plcproj", "*.tspproj", "*.TcPOU", "*.TcDUT", "*.TcGVL"},
		Sections: []Section{
			{
				Title:   "Build artifacts",
				Comment: "Compiled PLC output, regenerated on every build.",
				Patterns: []Pattern{
					{Glob: "_Boot/"},
					{Glob: "_CompileInfo/"},
					{Glob: "*.compiled"},
					{Glob: "*.compileinfo"},
					{Glob: "*.tpy", Comment: "TwinCAT 2 PLC project export — build artifact"},
				},
			},
			{
				Title: "LineIDs",
				Comment: "Breakpoint tracking metadata. Changes on every save,\n" +
					"pollutes diffs, and is user-specific. Safe to ignore.\n" +
					"TwinCAT 3.1.4026+ allows disabling LineID generation entirely.",
				Patterns: []Pattern{
					{Glob: "*.TcLIDs"},
				},
			},
			{
				Title: "Type system cache",
				Comment: "TMC files are regenerated from source on compile.\n" +
					"Do NOT ignore if your project uses them for ADS symbol access\n" +
					"and you don't recompile on deploy. Comment out the line below if so.",
				Patterns: []Pattern{
					{Glob: "*.tmc"},
				},
			},
			{
				Title:   "User settings",
				Comment: "Per-user Visual Studio and TwinCAT settings.",
				Patterns: []Pattern{
					{Glob: "*.suo"},
					{Glob: "*.user"},
					{Glob: ".vs/"},
					{Glob: "*.dbmdl"},
					{Glob: "*.jfm"},
					{Glob: "TwinCAT RT (x64)*/"},
				},
			},
			{
				Title:   "Hardware scan cache",
				Comment: "EtherCAT device scan results, machine-specific.",
				Patterns: []Pattern{
					{Glob: "_ModulesScan/"},
					{Glob: "DeviceDescription/"},
				},
			},
			{
				Title:   "License files",
				Comment: "Vendor and machine-bound license blobs — never commit.",
				Patterns: []Pattern{
					{Glob: "*.tclrs"},
					{Glob: "*.reg"},
				},
			},
		},
		Attributes: []AttrSection{
			{
				Title: "Merge strategy",
				Comment: "TwinCAT XML files should use union merge to reduce\n" +
					"false conflicts from concurrent edits to different POUs.",
				Rules: []AttrRule{
					{Pattern: "*.tsproj", Attrs: "merge=union"},
					{Pattern: "*.plcproj", Attrs: "merge=union"},
					{Pattern: "*.tspproj", Attrs: "merge=union"},
				},
			},
			{
				Title:   "Binary files",
				Comment: "These should not be diffed as text.",
				Rules: []AttrRule{
					{Pattern: "*.tmc", Attrs: "binary"},
					{Pattern: "*.compiled", Attrs: "binary"},
				},
			},
			{
				Title: "Diff driver",
				Comment: "If plc-project-diff is installed, use it for semantic diffs.\n" +
					"Uncomment after installing plc-project-diff:\n" +
					"# *.TcPOU  diff=plc\n" +
					"# *.TcDUT  diff=plc\n" +
					"# *.TcGVL  diff=plc",
				Rules: nil,
			},
			{
				Title:   "Line endings",
				Comment: "Force LF for all text files (consistency across Windows/Linux).",
				Rules: []AttrRule{
					{Pattern: "*", Attrs: "text=auto eol=lf"},
					{Pattern: "*.st", Attrs: "text eol=lf"},
					{Pattern: "*.ST", Attrs: "text eol=lf"},
				},
			},
		},
		HMISections: []Section{
			{
				Title:   "TwinCAT HMI",
				Comment: "TwinCAT HMI build output and source maps.",
				Patterns: []Pattern{
					{Glob: "*.hmi.usercontrol.js.map"},
					{Glob: "TcHmi/www/"},
					{Glob: ".tchmibuild/"},
					{Glob: "**/PublishedContent/"},
				},
			},
		},
		LFSAttrs: []AttrSection{
			{
				Title:   "Git LFS",
				Comment: "Large binary blobs — track via LFS instead of committing directly.",
				Rules: []AttrRule{
					{Pattern: "*.cam", Attrs: "filter=lfs diff=lfs merge=lfs -text", Comment: "CAM tables"},
					{Pattern: "*.png", Attrs: "filter=lfs diff=lfs merge=lfs -text", Comment: "HMI images"},
					{Pattern: "*.jpg", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.bin", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
				},
			},
		},
		HookWarnings: []HookWarning{
			{Pattern: "_Boot/", Reason: "boot project — should be built from source"},
			{Pattern: "_CompileInfo/", Reason: "build artifact"},
			{Pattern: "*.TcLIDs", Reason: "LineID file — pollutes diffs"},
			{Pattern: "*.tmc", Reason: "type cache — regenerated on build"},
			{Pattern: "*.compiled", Reason: "build artifact"},
		},
	}
}
