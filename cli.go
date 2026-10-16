package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
	"Go_Reforged/transform"
)

// Codes de sortie.
const (
	exitOK      = 0 // succès (des avertissements ont pu être affichés)
	exitUsage   = 1 // erreur d'utilisation ou de fichier
	exitContent = 2 // erreur de contenu
	exitCheck   = 3 // --check : des avertissements ont été trouvés
)

const usageLine = "usage: go run . [options] <input> <output>"

type config struct {
	input, output string
	sentenceCase  bool
	articlesFile  string
	strict        bool
	check         bool
	workers       int
}

// run exécute le programme et renvoie le code de sortie.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, err := parseArgs(args, stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}
		return exitUsage
	}

	rules, code := buildRules(cfg, stderr)
	if code != exitOK {
		return code
	}

	in, name, err := openInput(cfg.input, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "cannot read input: %v\n", err)
		return exitUsage
	}
	defer in.Close()

	var buf bytes.Buffer // temporaire : remplacé par un fichier temporaire à l'étape 6
	diags, err := processStream(in, &buf, rules)
	if err != nil {
		fmt.Fprintf(stderr, "cannot read %s: %v\n", name, err)
		return exitUsage
	}

	warnings, errs := report(stderr, name, diags, cfg.strict)
	switch {
	case errs > 0:
		return exitContent
	case cfg.check:
		if warnings > 0 {
			return exitCheck
		}
		return exitOK
	}

	if err := writeOutput(cfg.output, buf.Bytes(), stdout); err != nil {
		fmt.Fprintf(stderr, "cannot write output: %v\n", err)
		return exitUsage
	}
	return exitOK
}

// parseArgs lit les options et les deux fichiers. Elle affiche elle-même ses erreurs.
func parseArgs(args []string, stderr io.Writer) (config, error) {
	var cfg config
	fs := flag.NewFlagSet("go-reforged", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.BoolVar(&cfg.sentenceCase, "sentence-case", false, "capitalize the first letter of each sentence")
	fs.StringVar(&cfg.articlesFile, "articles", "", "file of article exceptions")
	fs.BoolVar(&cfg.strict, "strict", false, "turn every warning into an error")
	fs.BoolVar(&cfg.check, "check", false, "report warnings and errors without writing anything")
	fs.IntVar(&cfg.workers, "workers", runtime.NumCPU(), "number of goroutines")

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}

	rest := fs.Args()
	switch {
	case len(rest) == 2:
		cfg.input, cfg.output = rest[0], rest[1]
	case len(rest) == 1 && cfg.check:
		cfg.input = rest[0]
	default:
		fmt.Fprintln(stderr, usageLine)
		return cfg, errors.New("wrong number of arguments")
	}
	if cfg.workers < 1 {
		fmt.Fprintln(stderr, "--workers must be at least 1")
		return cfg, errors.New("invalid workers")
	}
	return cfg, nil
}

// buildRules assemble les règles dans l'ordre imposé par le sujet.
func buildRules(cfg config, stderr io.Writer) ([]transform.Rule, int) {
	var exceptions map[string]string
	if cfg.articlesFile != "" {
		f, err := os.Open(cfg.articlesFile)
		if err != nil {
			fmt.Fprintf(stderr, "cannot open articles file: %v\n", err)
			return nil, exitUsage
		}
		defer f.Close()

		var diags []lexer.Diagnostic
		exceptions, diags, err = transform.ParseExceptions(f)
		if err != nil {
			fmt.Fprintf(stderr, "cannot read articles file: %v\n", err)
			return nil, exitUsage
		}
		if len(diags) > 0 {
			report(stderr, cfg.articlesFile, diags, false)
			return nil, exitContent
		}
	}

	rules := []transform.Rule{
		transform.Blocks{},
		transform.Modifiers{},
		transform.Articles{Exceptions: exceptions},
		transform.Punctuation{},
	}
	if cfg.sentenceCase {
		rules = append(rules, &transform.SentenceCase{})
	}
	return rules, exitOK
}

func openInput(path string, stdin io.Reader) (io.ReadCloser, string, error) {
	if path == "-" {
		return io.NopCloser(stdin), "<stdin>", nil
	}
	f, err := os.Open(path)
	return f, path, err
}

// processStream lit r ligne par ligne, applique les règles et écrit dans w.
func processStream(r io.Reader, w io.Writer, rules []transform.Rule) ([]lexer.Diagnostic, error) {
	reader := bufio.NewReader(r)
	var all []lexer.Diagnostic

	for lineNo := 1; ; lineNo++ {
		line, readErr := reader.ReadString('\n')
		if line == "" {
			if readErr == io.EOF {
				return all, nil
			}
			return all, readErr
		}

		body, ending := splitLineEnding(line)
		tokens, diags := transform.Run(lexer.Lex(body, lineNo), rules)
		all = append(all, diags...)
		if _, err := io.WriteString(w, render.Render(tokens)+ending); err != nil {
			return all, err
		}

		if readErr != nil { // dernière ligne, sans saut de ligne final
			if readErr == io.EOF {
				return all, nil
			}
			return all, readErr
		}
	}
}

// splitLineEnding sépare une ligne de sa fin ("\n", "\r\n" ou rien).
func splitLineEnding(line string) (body, ending string) {
	switch {
	case strings.HasSuffix(line, "\r\n"):
		return line[:len(line)-2], "\r\n"
	case strings.HasSuffix(line, "\n"):
		return line[:len(line)-1], "\n"
	}
	return line, ""
}

// report affiche les diagnostics au format <file>:<line>:<column>: <severity>: <message>.
// En mode strict, un avertissement est affiché et compté comme une erreur.
func report(stderr io.Writer, file string, diags []lexer.Diagnostic, strict bool) (warnings, errs int) {
	for _, d := range diags {
		severity := d.Severity
		if strict && severity == "warning" {
			severity = "error"
		}
		if severity == "error" {
			errs++
		} else {
			warnings++
		}
		fmt.Fprintf(stderr, "%s:%d:%d: %s: %s\n", file, d.Line, d.Col, severity, d.Message)
	}
	return warnings, errs
}

func writeOutput(path string, data []byte, stdout io.Writer) error {
	if path == "-" {
		_, err := stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}