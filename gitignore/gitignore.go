// Package gitignore matches paths against the .gitignore files of a project.
//
// Patterns are read from the Git repository root down to the directory being
// scanned. Nested .gitignore files apply inside their own directory. Outside a
// repository, rules start at the directory passed to New.
package gitignore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Set is an ordered list of gitignore patterns. The last matching pattern wins.
// A path whose parent directory is ignored stays ignored.
type Set struct {
	root     string
	patterns []pattern
}

type pattern struct {
	negate   bool
	dirOnly  bool
	domain   []string
	segments []seg
}

type seg struct {
	any  bool
	glob string
}

// New loads the .gitignore rules that apply inside dir.
// dir itself is included, so patterns in dir/.gitignore are visible.
func New(dir string) *Set {
	abs := absPath(dir)
	s := ancestors(abs)
	s.load(abs)
	return s
}

// Enter returns a copy of s with dir's .gitignore appended.
// dir must be a directory under the same root.
func (s *Set) Enter(dir string) *Set {
	n := &Set{
		root:     s.root,
		patterns: append([]pattern(nil), s.patterns...),
	}
	n.load(absPath(dir))
	return n
}

// Ignored reports whether path is excluded by the loaded patterns.
// isDir must be true when path is a directory: directory-only patterns
// (a trailing slash in .gitignore) match only directories.
func (s *Set) Ignored(path string, isDir bool) bool {
	parts, ok := s.pathParts(path)
	if !ok || len(parts) == 0 {
		return false
	}
	for i := 1; i < len(parts); i++ {
		if s.excluded(parts[:i], true) {
			return true
		}
	}
	return s.excluded(parts, isDir)
}

func ancestors(abs string) *Set {
	root := findRoot(abs)
	s := &Set{root: root}
	var chain []string
	for d := abs; d != root; d = filepath.Dir(d) {
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		chain = append(chain, parent)
		if parent == root {
			break
		}
	}
	for i := len(chain) - 1; i >= 0; i-- {
		s.load(chain[i])
	}
	return s
}

func findRoot(abs string) string {
	d := abs
	for {
		if st, err := os.Stat(filepath.Join(d, ".git")); err == nil && (st.IsDir() || st.Mode().IsRegular()) {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return abs
		}
		d = parent
	}
}

func (s *Set) load(dir string) {
	data, err := ioutil.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	domain, ok := relParts(s.root, dir)
	if !ok {
		return
	}
	text := strings.Replace(string(data), "\r\n", "\n", -1)
	text = strings.Replace(text, "\r", "\n", -1)
	for _, line := range strings.Split(text, "\n") {
		if p, ok := parseLine(line, domain); ok {
			s.patterns = append(s.patterns, p)
		}
	}
}

func (s *Set) pathParts(path string) ([]string, bool) {
	return relParts(s.root, absPath(path))
}

func (s *Set) excluded(parts []string, isDir bool) bool {
	ignored := false
	for i := range s.patterns {
		if !s.patterns[i].match(parts, isDir) {
			continue
		}
		ignored = !s.patterns[i].negate
	}
	return ignored
}

func (p pattern) match(parts []string, isDir bool) bool {
	if p.dirOnly && !isDir {
		return false
	}
	if len(parts) < len(p.domain) {
		return false
	}
	for i, name := range p.domain {
		if parts[i] != name {
			return false
		}
	}
	return matchSegs(p.segments, parts[len(p.domain):])
}

func parseLine(line string, domain []string) (pattern, bool) {
	line = trimTrailingSpace(line)
	if line == "" || line[0] == '#' {
		return pattern{}, false
	}
	p := pattern{domain: append([]string(nil), domain...)}
	if line[0] == '!' {
		p.negate = true
		line = line[1:]
	}
	if line == "" {
		return pattern{}, false
	}
	if hasTrailingSlash(line) {
		p.dirOnly = true
		line = line[:len(line)-1]
	}
	if line == "" {
		return pattern{}, false
	}
	anchored := false
	if line[0] == '/' {
		anchored = true
		line = line[1:]
	}
	if line == "" {
		return pattern{}, false
	}
	raw := splitSlash(line)
	if len(raw) > 1 {
		anchored = true
	}
	// A pattern with no slash matches in any directory.
	if !anchored {
		raw = append([]string{"**"}, raw...)
	}
	for _, part := range raw {
		if part == "**" {
			p.segments = append(p.segments, seg{any: true})
			continue
		}
		p.segments = append(p.segments, seg{glob: part})
	}
	if len(p.segments) == 0 {
		return pattern{}, false
	}
	return p, true
}

// trimTrailingSpace drops trailing spaces unless they are escaped with '\'.
func trimTrailingSpace(s string) string {
	i := len(s)
	for i > 0 && s[i-1] == ' ' {
		slashes := 0
		for j := i - 2; j >= 0 && s[j] == '\\'; j-- {
			slashes++
		}
		if slashes%2 == 1 {
			break
		}
		i--
	}
	return s[:i]
}

func hasTrailingSlash(s string) bool {
	if !strings.HasSuffix(s, "/") {
		return false
	}
	slashes := 0
	for i := len(s) - 2; i >= 0 && s[i] == '\\'; i-- {
		slashes++
	}
	return slashes%2 == 0
}

func splitSlash(pat string) []string {
	var parts []string
	var b strings.Builder
	for i := 0; i < len(pat); i++ {
		if pat[i] == '\\' && i+1 < len(pat) {
			b.WriteByte('\\')
			i++
			b.WriteByte(pat[i])
			continue
		}
		if pat[i] == '/' {
			parts = append(parts, b.String())
			b.Reset()
			continue
		}
		b.WriteByte(pat[i])
	}
	parts = append(parts, b.String())
	return parts
}

func matchSegs(segs []seg, parts []string) bool {
	return matchSegsAt(segs, parts, 0, 0)
}

func matchSegsAt(segs []seg, parts []string, si, pi int) bool {
	for si < len(segs) {
		if segs[si].any {
			// A trailing "**" matches files inside, so it needs at least one
			// path segment ("dir/**" does not ignore "dir" itself).
			// "**" anywhere else matches zero or more segments ("a/**/b"
			// matches "a/b").
			if si == len(segs)-1 {
				return pi < len(parts)
			}
			for n := pi; n <= len(parts); n++ {
				if matchSegsAt(segs, parts, si+1, n) {
					return true
				}
			}
			return false
		}
		if pi >= len(parts) || !glob(segs[si].glob, parts[pi]) {
			return false
		}
		si++
		pi++
	}
	return pi == len(parts)
}

// glob reports whether name matches a single gitignore path segment.
func glob(pat, s string) bool {
	for {
		if pat == "" {
			return s == ""
		}
		switch pat[0] {
		case '*':
			pat = strings.TrimLeft(pat, "*")
			if pat == "" {
				return true
			}
			for i := 0; ; {
				if glob(pat, s[i:]) {
					return true
				}
				if i >= len(s) {
					return false
				}
				_, size := utf8.DecodeRuneInString(s[i:])
				if size <= 0 {
					return false
				}
				i += size
			}
		case '?':
			if s == "" {
				return false
			}
			_, size := utf8.DecodeRuneInString(s)
			pat = pat[1:]
			s = s[size:]
		case '[':
			end := classEnd(pat)
			if end < 0 {
				if s == "" || s[0] != '[' {
					return false
				}
				pat = pat[1:]
				s = s[1:]
				continue
			}
			if s == "" || !matchClass(pat[1:end], runeAt(s)) {
				return false
			}
			_, size := utf8.DecodeRuneInString(s)
			pat = pat[end+1:]
			s = s[size:]
		case '\\':
			pat = pat[1:]
			if pat == "" || s == "" {
				return false
			}
			pr, psize := utf8.DecodeRuneInString(pat)
			sr, ssize := utf8.DecodeRuneInString(s)
			if sr != pr {
				return false
			}
			pat = pat[psize:]
			s = s[ssize:]
		default:
			if s == "" {
				return false
			}
			pr, psize := utf8.DecodeRuneInString(pat)
			sr, ssize := utf8.DecodeRuneInString(s)
			if sr != pr {
				return false
			}
			pat = pat[psize:]
			s = s[ssize:]
		}
	}
}

func runeAt(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

// classEnd returns the index of the closing ']' of a character class.
// pat must start with '['. A ']' right after '[' or '[!' is literal.
func classEnd(pat string) int {
	i := 1
	if i < len(pat) && (pat[i] == '!' || pat[i] == '^') {
		i++
	}
	if i < len(pat) && pat[i] == ']' {
		i++
	}
	for ; i < len(pat); i++ {
		if pat[i] == ']' {
			return i
		}
	}
	return -1
}

func matchClass(class string, r rune) bool {
	neg := false
	if strings.HasPrefix(class, "!") || strings.HasPrefix(class, "^") {
		neg = true
		class = class[1:]
	}
	matched := false
	rs := []rune(class)
	for i := 0; i < len(rs); i++ {
		if i+2 < len(rs) && rs[i+1] == '-' {
			if r >= rs[i] && r <= rs[i+2] {
				matched = true
				break
			}
			i += 2
			continue
		}
		if r == rs[i] {
			matched = true
			break
		}
	}
	if neg {
		return !matched
	}
	return matched
}

func relParts(root, path string) ([]string, bool) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return nil, false
	}
	if rel == "." {
		return nil, true
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, false
	}
	return strings.Split(filepath.ToSlash(rel), "/"), true
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return abs
}
