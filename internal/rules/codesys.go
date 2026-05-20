package rules

// Codesys returns the Codesys rule set.
func Codesys() RuleSet {
	return RuleSet{
		Vendor:      "codesys",
		DisplayName: "Codesys",
		DetectGlobs: []string{"*.project", "*.library", "*.projectarchive"},
		Sections: []Section{
			{
				Title:   "Build artifacts",
				Comment: "Codesys build output, regenerated on every compile.",
				Patterns: []Pattern{
					{Glob: "*.compileinfo"},
					{Glob: "*.object"},
					{Glob: "*.compiled-library"},
					{Glob: "*.app", Comment: "deployable application — build from source"},
				},
			},
			{
				Title:   "Workspace metadata",
				Comment: "Eclipse/Codesys workspace state, user-specific.",
				Patterns: []Pattern{
					{Glob: ".metadata/"},
					{Glob: ".settings/"},
					{Glob: ".project.user"},
				},
			},
			{
				Title:   "Codesys Python scripting",
				Comment: "Codesys exposes a Python scripting interface; ignore caches.",
				Patterns: []Pattern{
					{Glob: "__pycache__/"},
					{Glob: "*.pyc"},
				},
			},
			{
				Title:   "Device build cache",
				Comment: "Device-specific build caches and temporary files.",
				Patterns: []Pattern{
					{Glob: "_*.tmp/"},
					{Glob: "PlcLogic/", Comment: "compiled logic — regenerated"},
					{Glob: "VisuFiles/", Comment: "Visu build output"},
				},
			},
			{
				Title:   "Backup files",
				Comment: "Codesys auto-save backups.",
				Patterns: []Pattern{
					{Glob: "*.project.~u"},
					{Glob: "*.library.~u"},
					{Glob: "Backup/"},
				},
			},
		},
		Attributes: []AttrSection{
			{
				Title:   "Binary files",
				Comment: "Codesys project archives are binary and should not be diffed.",
				Rules: []AttrRule{
					{Pattern: "*.projectarchive", Attrs: "binary"},
					{Pattern: "*.library", Attrs: "binary"},
					{Pattern: "*.compiled-library", Attrs: "binary"},
				},
			},
			{
				Title:   "Line endings",
				Comment: "Force LF for IEC 61131-3 source text.",
				Rules: []AttrRule{
					{Pattern: "*", Attrs: "text=auto eol=lf"},
					{Pattern: "*.st", Attrs: "text eol=lf"},
					{Pattern: "*.iecst", Attrs: "text eol=lf"},
				},
			},
		},
		HMISections: []Section{
			{
				Title:   "Codesys Visu",
				Comment: "Visu (HMI) build output and generated assets.",
				Patterns: []Pattern{
					{Glob: "VisuFiles/"},
					{Glob: "Visu/*.compiled"},
					{Glob: "WebVisu/"},
				},
			},
		},
		LFSAttrs: []AttrSection{
			{
				Title:   "Git LFS",
				Comment: "Large binary assets — track via LFS.",
				Rules: []AttrRule{
					{Pattern: "*.png", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.jpg", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.bmp", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
					{Pattern: "*.bin", Attrs: "filter=lfs diff=lfs merge=lfs -text"},
				},
			},
		},
		HookWarnings: []HookWarning{
			{Pattern: "*.compileinfo", Reason: "build artifact"},
			{Pattern: "*.object", Reason: "build artifact"},
			{Pattern: "*.app", Reason: "deployable — build from source"},
			{Pattern: ".metadata/", Reason: "workspace state — user-specific"},
		},
	}
}
