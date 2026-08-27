package service

import (
	"fmt"
	"strings"
)

type Command struct {
	Name  string
	Flags map[string]string
	Args  []string
}

func Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, fmt.Errorf("missing command")
	}
	cmd := Command{Name: strings.TrimSpace(args[0]), Flags: map[string]string{}}
	for i := 1; i < len(args); i++ {
		tok := args[i]
		if !strings.HasPrefix(tok, "--") {
			cmd.Args = append(cmd.Args, tok)
			continue
		}
		key, value, hasValue := strings.Cut(strings.TrimPrefix(tok, "--"), "=")
		if !hasValue {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				value = args[i+1]
				i++
			} else {
				value = "true"
			}
		}
		cmd.Flags[key] = value
	}
	return cmd, nil
}

func (c Command) String(name string) string {
	if v, ok := c.Flags[name]; ok {
		return v
	}
	return ""
}

func (c Command) MustString(name string) (string, error) {
	if v := c.String(name); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing --%s", name)
}

func (c Command) Bool(name string) bool {
	v := strings.ToLower(c.String(name))
	return v == "1" || v == "true" || v == "yes"
}

