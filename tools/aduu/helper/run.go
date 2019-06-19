package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Run(cmd *exec.Cmd, name string, settings ...ExecuteSetting) {
	if err := RunE(cmd, name, settings...); err != nil {
		os.Exit(1)
	}
}

func RunE(cmd *exec.Cmd, name string, settings ...ExecuteSetting) (err error) {
	setting := getSetting(settings)
	applySettings(cmd, setting)

	 _, err = runWithSettings(cmd, setting)
	if err != nil {
		fmt.Printf("%s failed: %v\n", name, err)
		return err
	}
	return nil
}

func getSetting(settings []ExecuteSetting) (setting ExecuteSetting) {
	if len(settings) > 1 {
		panic("settings only read for len=1")
	}
	if len(settings) == 0 {
		setting = ExecuteSetting{}
	} else {
		setting = settings[0]
	}
	return
}



func applySettings(cmd *exec.Cmd, setting ExecuteSetting) {
	if setting.Dir != "" {
		cmd.Dir = os.ExpandEnv(setting.Dir)
	}

	if setting.output {
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
}

func runWithSettings(cmd *exec.Cmd, setting ExecuteSetting) (out string, err error) {
	var timer *time.Timer
	exited := false
	var errFromTimeout error

	if setting.timeout != 0 {
		timer = time.AfterFunc(setting.timeout, func() {
			if !exited {
				errFromTimeout = cmd.Process.Kill()
			}
		})
	}

	if setting.start && setting.output {
		panic("can't do start with output")
	} else if setting.start && !setting.output {
		err = cmd.Start()
	} else if !setting.start && setting.output {
		var byteOut []byte
		byteOut, err = cmd.Output()
		if byteOut != nil {
			out = strings.TrimSpace(string(byteOut))
		}
	} else  {
		err = cmd.Run()
	}

	if errFromTimeout != nil {
		if err != nil {
			err = fmt.Errorf("error from timeout %v wraps err: %v", errFromTimeout, err)
		} else {
			err = fmt.Errorf("error from timeout: %v", errFromTimeout)
		}
	}

	exited = true
	if timer != nil {
		timer.Stop()
	}
	return
}

func RunWithOutput(cmd *exec.Cmd, name string, settings ...ExecuteSetting) string {
	out, err := RunWithOutputE(cmd, name, settings...)
	if err != nil {
		os.Exit(1)
	}
	return out
}

func RunWithOutputE(cmd *exec.Cmd, name string, settings ...ExecuteSetting) (out string, err error) {
	setting := getSetting(settings).withOutput()
	applySettings(cmd, setting)

	out, err = runWithSettings(cmd, setting)
	if err != nil {
		fmt.Printf("%s - %v failed with %v\n", name, cmd, err)
		return out, err
	}
	return
}

func first(s string) string {
	return string([]rune(s)[0])
}

func search(text []rune, what string) int {
	whatRunes := []rune(what)

	for i := range text {
		found := true
		for j := range whatRunes {
			if text[i+j] != whatRunes[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

func SplitCommand(s string) []string {
	runes := []rune(s)
	var values []string

	for len(runes) != 0 {

		// TODO: ' and " can be in the middle of a token, so checking the start of a string is not enough.
		if runes[0] == '"' {
			next := search(runes[1:], "\"")
			if next == -1 {
				panic("no matching \": " + s)
			}

			values = append(values, string(runes[1:next+1]))
			runes = runes[next+2:]
		} else {
			// TODO: If ' before next space, do skip this space and go instead to the next ' to continue from there.
			next := search(runes, " ")
			if next == -1 {
				values = append(values, string(runes))
				runes = []rune{}
				break
			}

			nextSpecial := search(runes, "'")

			// Early return if nextSpecial is not so.
			if nextSpecial == -1 || nextSpecial > next {
				values = append(values, string(runes[:next]))
				runes = runes[next+1:]
				continue
			}

			var current []rune

			for next != -1 && nextSpecial != -1 && nextSpecial < next {
				// Skip newlines enclosed by specials '  '.

				// Add to current.
				current = append(current, runes[:nextSpecial]...)
				// Move ahead.
				runes = runes[nextSpecial+1:]

				// Find the closing special character.
				nextSpecial = search(runes, "'")
				if nextSpecial == -1 {
					panic("no matching ': " + s)
				}
				// Add up to the closing parameter to current.
				current = append(current, runes[:nextSpecial]...)
				// Mve ahead.
				runes = runes[nextSpecial+1:]

				nextSpecial = search(runes, "'")
				next = search(runes, " ")
			}
			nextSpecial = search(runes, "'")
			next = search(runes, " ")
			if next == -1 {
				// Add current to last part.
				values = append(values, string(append(current, runes...)))
				// Move ahead.
				runes = []rune{}
				break
			}

			// Add current to closing part.
			values = append(values, string(append(current, runes[:next]...)))
			// Move ahead.
			runes = runes[next+1:]
		}
	}

	out := make([]string, len(values))
	for i := range values {
		out[i] = string(values[i])
	}
	return out
}

func SplitCommand2(s string) []string {
	input := s
	var values []string

	for len(s) != 0 {
		if s[0] == '"' {
			next := strings.Index(s[1:], "\"")
			if next == -1 {
				panic("no matching \": " + input)
			}

			values = append(values, s[1:next+1])
			s = s[next+2:]
		} else {
			next := strings.Index(s, " ")
			if next == -1 {
				values = append(values, s)
				s = ""
				break
			}

			values = append(values, s[:next])
			s = s[next+1:]
		}
	}

	return values
}

func Execute(s string, obj interface{}, settings ...ExecuteSetting) *exec.Cmd {
	cmd, _ := ExecuteE(s, obj, settings...)
	return cmd
}

func ExecuteE(s string, obj interface{}, settings ...ExecuteSetting) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(s))

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, obj); err != nil {
		panic(err)
	}

	fmt.Println(tpl.String())

	args := SplitCommand(os.ExpandEnv(tpl.String()))
	//args := strings.Split(tpl.String(), " ")

	if len(args) == 0 {
		panic("template length was zero.")
	}

	cmd := exec.Command(args[0], args[1:]...)
	return cmd, RunE(cmd, args[0], settings...)
}

func ExecuteWithOutput(s string, obj interface{}, settings ...ExecuteSetting) string {
	out, err := ExecuteWithOutputE(s, obj, settings...)
	if err != nil {
		fmt.Printf("Command failed with %v\n", err)
		os.Exit(1)
	}
	return out
}

func ExecuteWithOutputE(s string, obj interface{}, settings ...ExecuteSetting) (out string, err error) {
	t := template.Must(template.New("").Parse(s))

	var tpl bytes.Buffer

	if err := t.Execute(&tpl, obj); err != nil {
		panic(err)
	}

	fmt.Println(tpl.String())

	args := SplitCommand(os.ExpandEnv(tpl.String()))

	if len(args) == 0 {
		panic("template length was zero.")
	}

	cmd := exec.Command(args[0], args[1:]...)
	return RunWithOutputE(cmd, args[0], settings...)
}

type ExecuteSetting struct {
	Dir string
	start bool
	output bool
	timeout time.Duration
}

func (s ExecuteSetting) WithStart() ExecuteSetting {
	s.start = true
	return s
}

var WithStart = ExecuteSetting{
	start: true,
}

func WithDir(path string) ExecuteSetting {
	return ExecuteSetting{
		Dir: path,
	}
}

func (s ExecuteSetting) withOutput() ExecuteSetting {
	s.output = true
	return s
}

func (s ExecuteSetting) WithTimeout(duration time.Duration) ExecuteSetting {
	s.timeout = duration
	return s
}

func WithTimeout(time time.Duration) ExecuteSetting {
	return ExecuteSetting{
		timeout: time,
	}
}


func Getwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func Printwd() {
	fmt.Println(Getwd())
}