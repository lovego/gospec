package gitignore

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestGlob(t *testing.T) {
	cases := []struct {
		pat, name string
		ok        bool
	}{
		{"*.log", "a.log", true},
		{"*.log", "a.log.txt", false},
		{"*.log", ".log", true},
		{"q*x", "qx", true},
		{"q*x", "qxx", true},
		{"q*x", "q.x", true},
		{"q*x", "q.txt", false},
		{"[abc].txt", "b.txt", true},
		{"[abc].txt", "d.txt", false},
		{"[!abc].txt", "d.txt", true},
		{"[!abc].txt", "a.txt", false},
		{"[a-c].txt", "b.txt", true},
		{"[a-c].txt", "d.txt", false},
		{"?.c", "x.c", true},
		{"?.c", "xx.c", false},
		{"?.c", ".c", false},
		{`\*.star`, "*.star", true},
		{`\*.star`, "a.star", false},
		{"*", ".gitignore", true},
		{"foo", "foo", true},
		{"foo", "food", false},
		{`space\ `, "space ", true},
		{`\#hash`, "#hash", true},
		{`\!bang`, "!bang", true},
	}
	for _, c := range cases {
		if got := glob(c.pat, c.name); got != c.ok {
			t.Errorf("glob(%q, %q) = %v, want %v", c.pat, c.name, got, c.ok)
		}
	}
}

func TestIgnorePatterns(t *testing.T) {
	root := newRepo(t)
	defer os.RemoveAll(root)
	writeTree(t, root, map[string]string{
		".gitignore": `
# comment

*.log
!keep.log
foo/
/rooted.txt
doc/frotz
abc/**
!abc/keep.txt
vendor/
!vendor/keep.go
dir/**
*.swp
build
/only_here.txt
sub/special.go
space\ 
\#hash
\!bang
trailing   
`,
		"sub/.gitignore":            "nested_only.go\n!keep.log\n",
		"gen/.gitignore":            "*.gen\n",
		"a.log":                     "",
		"keep.log":                  "",
		"sub/a.log":                 "",
		"sub/keep.log":              "",
		"sub/keep2.log":             "",
		"foo/a.go":                  "",
		"food":                      "",
		"rooted.txt":                "",
		"sub/rooted.txt":            "",
		"doc/frotz/x.go":            "",
		"doc/other.txt":             "",
		"a/doc/frotz":               "",
		"abc/keep.txt":              "",
		"abc/skip.txt":              "",
		"abc/d/x.txt":               "",
		"vendor/keep.go":            "",
		"vendor/x.go":               "",
		"dir/file.txt":              "",
		"dir/nested/x.go":           "",
		"file.swp":                  "",
		"a/file.swp":                "",
		"readme.md":                 "",
		"build/a.go":                "",
		"a/build/b.go":              "",
		"only_here.txt":             "",
		"sub/only_here.txt":         "",
		"sub/special.go":            "",
		"special.go":                "",
		"a/sub/special.go":          "",
		"sub/nested_only.go":        "",
		"sub/nested/nested_only.go": "",
		"sub/other.go":              "",
		"gen/a.gen":                 "",
		"gen/b.txt":                 "",
		"a.gen":                     "",
		"space ":                    "",
		"#hash":                     "",
		"!bang":                     "",
		"trailing":                  "",
	})

	assertIgnored(t, root, map[string]bool{
		"a.log":                     true,
		"keep.log":                  false,
		"sub/a.log":                 true,
		"sub/keep.log":              false,
		"sub/keep2.log":             true,
		"foo":                       true,
		"foo/a.go":                  true,
		"food":                      false,
		"rooted.txt":                true,
		"sub/rooted.txt":            false,
		"doc/frotz":                 true,
		"doc/frotz/x.go":            true,
		"doc/other.txt":             false,
		"a/doc/frotz":               false,
		"abc":                       false,
		"abc/keep.txt":              false,
		"abc/skip.txt":              true,
		"abc/d/x.txt":               true,
		"vendor":                    true,
		"vendor/keep.go":            true,
		"vendor/x.go":               true,
		"dir":                       false,
		"dir/file.txt":              true,
		"dir/nested":                true,
		"dir/nested/x.go":           true,
		"file.swp":                  true,
		"a/file.swp":                true,
		"readme.md":                 false,
		"build":                     true,
		"build/a.go":                true,
		"a/build":                   true,
		"a/build/b.go":              true,
		"only_here.txt":             true,
		"sub/only_here.txt":         false,
		"sub/special.go":            true,
		"special.go":                false,
		"a/sub/special.go":          false,
		"sub/nested_only.go":        true,
		"sub/nested/nested_only.go": true,
		"sub/other.go":              false,
		"gen/a.gen":                 true,
		"gen/b.txt":                 false,
		"a.gen":                     false,
		"space ":                    true,
		"#hash":                     true,
		"!bang":                     true,
		"trailing":                  true,
	})
	assertSameAsGit(t, root)
}

func TestReincludeAndGlobs(t *testing.T) {
	root := newRepo(t)
	defer os.RemoveAll(root)
	writeTree(t, root, map[string]string{
		".gitignore": `
*.go
!important.go
[abc].txt
[!abc].txt
?.c
a/**/b
bar/**
**/generated/**
\*.star
q*x
`,
		"x.go":                "",
		"important.go":        "",
		"sub/important.go":    "",
		"sub/x.go":            "",
		"a.txt":               "",
		"b.txt":               "",
		"d.txt":               "",
		"x.c":                 "",
		"xx.c":                "",
		".c":                  "",
		"a/b/c.txt":           "",
		"a/x/b/c.txt":         "",
		"a/x/y/b":             "",
		"a/x/y/keep.md":       "",
		"bar/z.txt":           "",
		"bar/z/f.txt":         "",
		"x/generated/a.go":    "",
		"x/generated/y/z.txt": "",
		"generated/keep.txt":  "",
		"*.star":              "",
		"a.star":              "",
		"qx":                  "",
		"qxx":                 "",
		"q.x":                 "",
		"q.dat":               "",
	})
	assertIgnored(t, root, map[string]bool{
		"x.go":                true,
		"important.go":        false,
		"sub/important.go":    false,
		"sub/x.go":            true,
		"a.txt":               true,
		"b.txt":               true,
		"d.txt":               true,
		"x.c":                 true,
		"xx.c":                false,
		".c":                  false,
		"a/b":                 true,
		"a/b/c.txt":           true,
		"a/x/b":               true,
		"a/x/b/c.txt":         true,
		"a/x/y/b":             true,
		"a/x/y/keep.md":       false,
		"bar":                 false,
		"bar/z.txt":           true,
		"bar/z/f.txt":         true,
		"generated":           false,
		"generated/keep.txt":  true,
		"x/generated":         false,
		"x/generated/a.go":    true,
		"x/generated/y/z.txt": true,
		"*.star":              true,
		"a.star":              false,
		"qx":                  true,
		"qxx":                 true,
		"q.x":                 true,
		"q.dat":               false,
	})
	assertSameAsGit(t, root)
}

func TestParentReinclude(t *testing.T) {
	root := newRepo(t)
	defer os.RemoveAll(root)
	writeTree(t, root, map[string]string{
		".gitignore":      "/*\n!/src\n/src/*\n!/src/ok\n",
		"other.txt":       "",
		"src/no.go":       "",
		"src/ok/a.go":     "",
		"src/ok/sub/b.go": "",
		"foo/a.go":        "",
	})
	assertIgnored(t, root, map[string]bool{
		"other.txt":       true,
		"src":             false,
		"src/no.go":       true,
		"src/ok":          false,
		"src/ok/a.go":     false,
		"src/ok/sub":      false,
		"src/ok/sub/b.go": false,
		"foo":             true,
		"foo/a.go":        true,
	})
	assertSameAsGit(t, root)
}

func TestDirOnlyDoesNotMatchFile(t *testing.T) {
	root := newRepo(t)
	defer os.RemoveAll(root)
	writeTree(t, root, map[string]string{
		".gitignore": "foo/\n",
		"foo":        "",
	})
	assertIgnored(t, root, map[string]bool{"foo": false})
	assertSameAsGit(t, root)
}

func newRepo(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "gospec-gitignore")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertIgnored(t *testing.T, root string, want map[string]bool) {
	t.Helper()
	for rel, ignored := range want {
		path := filepath.Join(root, filepath.FromSlash(rel))
		fi, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s: %v", rel, err)
		}
		dir := path
		if !fi.IsDir() {
			dir = filepath.Dir(path)
		}
		if got := New(dir).Ignored(path, fi.IsDir()); got != ignored {
			t.Errorf("%s ignored = %v, want %v", rel, got, ignored)
		}
	}
}

func assertSameAsGit(t *testing.T, root string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	cmd := exec.Command("git", "init", "-q", "--template=")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git init: %v: %s", err, out)
	}
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." || rel == ".git" || strings.HasPrefix(rel, ".git"+string(filepath.Separator)) {
			if fi.IsDir() && rel == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		gotDir := path
		if !fi.IsDir() {
			gotDir = filepath.Dir(path)
		}
		got := New(gotDir).Ignored(path, fi.IsDir())
		want, err := gitIgnored(root, filepath.ToSlash(rel))
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("git disagreement on %s (dir=%v): matcher=%v git=%v", rel, fi.IsDir(), got, want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

var checkIgnorePattern = regexp.MustCompile(`^.*:\d+:(.*)$`)

func gitIgnored(root, rel string) (bool, error) {
	cmd := exec.Command("git",
		"-c", "core.excludesFile=/dev/null",
		"-C", root,
		"check-ignore", "-v", "--no-index", "--", rel)
	out, err := cmd.Output()
	text := strings.TrimRight(string(out), "\n")
	if text == "" {
		if err == nil {
			return false, nil
		}
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) == 0 {
			return false, nil
		}
		if ee, ok := err.(*exec.ExitError); ok {
			return false, fmt.Errorf("check-ignore %s: %v: %s", rel, err, ee.Stderr)
		}
		return false, err
	}
	tab := strings.LastIndex(text, "\t")
	if tab < 0 {
		return false, fmt.Errorf("check-ignore %s: unexpected output %q", rel, text)
	}
	m := checkIgnorePattern.FindStringSubmatch(text[:tab])
	if m == nil {
		return false, fmt.Errorf("check-ignore %s: unexpected pattern %q", rel, text[:tab])
	}
	pattern := m[1]
	return pattern == "" || pattern[0] != '!', nil
}
