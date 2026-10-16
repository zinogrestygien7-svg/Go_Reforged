package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runCLI(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errb bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errb)
	return code, out.String(), errb.String()
}

func TestRunStdin(t *testing.T) {
	tests := []struct {
		name       string
		stdin      string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{"exemple du sujet", "hello (up) [cap]big world[/cap]\n", []string{"-", "-"},
			0, "HELLO Big World\n", ""},
		{"fins de ligne Windows et lignes vides", "a (up)\r\n\r\nb\r\n", []string{"-", "-"},
			0, "A\r\n\r\nb\r\n", ""},
		{"sans saut de ligne final", "a (up)", []string{"-", "-"},
			0, "A", ""},
		{"entrée vide", "", []string{"-", "-"}, 0, "", ""},
		{"avertissement : succès", "go (up, -1)\n", []string{"-", "-"},
			0, "go (up, -1)\n", "<stdin>:1:4: warning: invalid count in modifier \"(up, -1)\"\n"},
		{"strict : rien n'est écrit", "go (up, -1)\n", []string{"--strict", "-", "-"},
			2, "", "<stdin>:1:4: error: invalid count in modifier \"(up, -1)\"\n"},
		{"bloc non fermé", "[up]hello\n", []string{"-", "-"},
			2, "", "<stdin>:1:1: error: block \"[up]\" is never closed\n"},
		{"check avec avertissement", "go (up, -1)\n", []string{"--check", "-"},
			3, "", "<stdin>:1:4: warning: invalid count in modifier \"(up, -1)\"\n"},
		{"check sans avertissement", "hello (up)\n", []string{"--check", "-"}, 0, "", ""},
		{"majuscule de phrase d'une ligne à l'autre", "first line.\nsecond line\n",
			[]string{"--sentence-case", "-", "-"}, 0, "First line.\nSecond line\n", ""},
		{"numéro de ligne dans le diagnostic", "ok\nZZ (hex)\n", []string{"-", "-"},
			0, "ok\nZZ (hex)\n", "<stdin>:2:4: warning: invalid hexadecimal number \"ZZ\"\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, stderr := runCLI(tc.stdin, tc.args...)
			if code != tc.wantCode {
				t.Errorf("code: got %d, want %d", code, tc.wantCode)
			}
			if stdout != tc.wantStdout {
				t.Errorf("stdout\n got: %q\nwant: %q", stdout, tc.wantStdout)
			}
			if stderr != tc.wantStderr {
				t.Errorf("stderr\n got: %q\nwant: %q", stderr, tc.wantStderr)
			}
		})
	}
}

func TestRunUsageErrors(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.txt")

	tests := []struct {
		name string
		args []string
	}{
		{"aucun argument", nil},
		{"un seul fichier", []string{"in.txt"}},
		{"trois fichiers", []string{"a", "b", "c"}},
		{"option inconnue", []string{"--nope", "-", "-"}},
		{"workers invalide", []string{"--workers", "0", "-", "-"}},
		{"fichier d'entrée introuvable", []string{missing, "-"}},
		{"fichier d'articles introuvable", []string{"--articles", missing, "-", "-"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if code, _, _ := runCLI("", tc.args...); code != 1 {
				t.Errorf("code: got %d, want 1", code)
			}
		})
	}
}

func TestRunFiles(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "sample.txt")
	out := filepath.Join(dir, "result.txt")
	if err := os.WriteFile(in, []byte("Simply add 42 (hex) and 10 (bin) ."), 0o644); err != nil {
		t.Fatal(err)
	}

	if code, _, _ := runCLI("", in, out); code != 0 {
		t.Fatalf("code: got %d, want 0", code)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if want := "Simply add 66 and 2."; string(got) != want {
		t.Errorf("result\n got: %q\nwant: %q", got, want)
	}
}

func TestRunSameInputAndOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "same.txt")
	if err := os.WriteFile(path, []byte("hello (up)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code, _, _ := runCLI("", path, path); code != 0 {
		t.Fatalf("code: got %d, want 0", code)
	}
	got, _ := os.ReadFile(path)
	if want := "HELLO\n"; string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunNoOutputOnError(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	out := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(in, []byte("[up]hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if code, _, _ := runCLI("", in, out); code != 2 {
		t.Fatalf("code: got %d, want 2", code)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Errorf("le fichier de sortie ne doit pas exister (err = %v)", err)
	}
}

func TestRunArticlesFile(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "articles.txt")
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(good, []byte("university a\nhouse a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bad, []byte("university\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, stdout, _ := runCLI("It is an university next to an house .\n", "--articles", good, "-", "-")
	if code != 0 || stdout != "It is a university next to a house.\n" {
		t.Errorf("code %d, stdout %q", code, stdout)
	}

	code, stdout, stderr := runCLI("hello\n", "--articles", bad, "-", "-")
	if code != 2 || stdout != "" {
		t.Errorf("code %d, stdout %q (want 2 et rien)", code, stdout)
	}
	if want := bad + ":1:1: error: malformed exception line \"university\"\n"; stderr != want {
		t.Errorf("stderr\n got: %q\nwant: %q", stderr, want)
	}
}
