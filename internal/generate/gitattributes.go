package generate

import (
	"fmt"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// Gitattributes returns the rendered .gitattributes content.
func Gitattributes(r rules.RuleSet, opts Options) string {
	var b strings.Builder

	writeBanner(&b, fmt.Sprintf("%s .gitattributes", r.DisplayName), opts.Version)

	sections := append([]rules.AttrSection(nil), r.Attributes...)
	if opts.WithLFS {
		sections = append(sections, r.LFSAttrs...)
	}

	for _, s := range sections {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("# --- %s ---\n", s.Title))
		if s.Comment != "" {
			for _, line := range strings.Split(s.Comment, "\n") {
				b.WriteString("# ")
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
		for _, rule := range s.Rules {
			if rule.Comment != "" {
				b.WriteString("# ")
				b.WriteString(rule.Comment)
				b.WriteString("\n")
			}
			b.WriteString(fmt.Sprintf("%s %s\n", padRight(rule.Pattern, 16), rule.Attrs))
		}
	}
	return b.String()
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
