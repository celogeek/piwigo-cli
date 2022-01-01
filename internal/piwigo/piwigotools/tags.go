package piwigotools

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Tags []*Tag

type Tag struct {
	Id           int        `json:"id,string"`
	Name         string     `json:"name"`
	LastModified TimeResult `json:"lastmodified"`
	Birthdate    TimeResult `json:"birthdate"`
	Url          string     `json:"url"`
	UrlName      string     `json:"url_name"`
	ImageUrl     string     `json:"page_url"`
}

func (t Tags) NamesWithAgeAt(createdAt *TimeResult) []string {
	names := make([]string, len(t))
	for i, tag := range t {
		names[i] = tag.NameWithAgeAt(createdAt)
	}
	return names
}

func (t *Tag) NameWithAgeAt(createdAt *TimeResult) string {
	bd := t.Birthdate.AgeAt(createdAt)
	if bd != "" {
		return fmt.Sprintf("%s (%s)", t.Name, bd)
	}
	return t.Name
}

func (t Tags) Select(exclude *regexp.Regexp) ([]*Tag, error) {
	fzf := "fzf --multi --height=30% --border --layout=reverse -e --bind=esc:clear-query --with-nth 2 --delimiter=\"\t\""
	cmd := exec.Command("sh", "-c", fzf)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		defer in.Close()
		for i, tag := range t {
			if exclude != nil && exclude.MatchString(tag.Name) {
				continue
			}
			in.Write([]byte(fmt.Sprintf("%d\t%s\n", i, tag.Name)))
		}
	}()
	out, _ := cmd.Output()
	rows := strings.Split(string(out), "\n")
	selections := make([]*Tag, 0, len(rows))
	for _, row := range rows {
		i, err := strconv.Atoi(strings.Split(row, "\t")[0])
		if err != nil {
			continue
		}
		selections = append(selections, t[i])
	}
	return selections, nil
}
