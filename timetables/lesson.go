package timetables

import (
	"strconv"
	"strings"
	"unicode"
)

// Subject stores the long and short names of a subject.
type Subject struct {
	ShortName string
	Name      string
}

var subjects = map[string]Subject{
	"liikunta": Subject{Name: "Liikunta", ShortName: "LI"},
	"ai":       Subject{Name: "Äidinkieli", ShortName: "AI"},
	"ma":       Subject{Name: "Matematiikka", ShortName: "MA"},
	"fy":       Subject{Name: "Fysiikka", ShortName: "FY"},
	"ke":       Subject{Name: "Kemia", ShortName: "KE"},
	"en":       Subject{Name: "Englanti", ShortName: "EN"},
	"fi":       Subject{Name: "Filosofia", ShortName: "FI"},
	"ps":       Subject{Name: "Psykologia", ShortName: "PS"},
	"tt":       Subject{Name: "Terveystieto", ShortName: "TT"},
	"yh":       Subject{Name: "Yhteiskuntaoppi", ShortName: "YH"},
	"ue":       Subject{Name: "Uskonto", ShortName: "UE"},
	"et":       Subject{Name: "Elämänkatsomustieto", ShortName: "ET"},
	"ru":       Subject{Name: "Ruotsi", ShortName: "RU"},
	"ge":       Subject{Name: "Maantieto", ShortName: "GE"},
	"bi":       Subject{Name: "Biologia", ShortName: "BI"},
}

var ignore = []string{
	"psil",
	"mat valm",
	"yhkvlh",
	"vlh",
	"alv",
	"tyhjää",
}

// Lesson stores a subject, the course number and the lesson number.
type Lesson struct {
	Subject        Subject
	Course, Lesson int
}

// ParseLesson attempts to parse a lesson from the given string.
func ParseLesson(str string) *Lesson {
	str = strings.ToLower(str)

	if strings.HasPrefix(str, "rt+") {
		str = str[3:]
	}

	for _, ign := range ignore {
		if strings.HasPrefix(str, ign) {
			return nil
		}
	}

	var subject = Subject{}
	var courseID, lessonID int

	for name, value := range subjects {
		if strings.HasPrefix(str, name) {
			subject = value
			break
		}
	}
	if subject.Name == "" {
		return nil
	}

	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			break
		}
		str = str[1:]
	}

	if len(str) == 0 {
		return &Lesson{Subject: subject, Course: 0, Lesson: 0}
	}

	if len(str) == 2 {
		if str == "11" {
			courseID = 1
			lessonID = 11
		} else {
			courseID, _ = strconv.Atoi(string(str[0]))
			lessonID, _ = strconv.Atoi(string(str[1]))
		}
	} else if str[1] == ' ' {
		courseID, _ = strconv.Atoi(string(str[0]))
		if strings.EqualFold(str[2:], "alkaa") {
			lessonID = 1
		} else {
			lessonID, _ = strconv.Atoi(str[2:])
		}
	} else if len(str) == 3 {
		courseID, _ = strconv.Atoi(string(str[0]))
		lessonID, _ = strconv.Atoi(str[1:])
	} else if strings.EqualFold(str, "alkaa") {
		courseID = 1
		lessonID = 1
	} else {
		return &Lesson{Subject: subject, Course: 0, Lesson: 0}
	}

	return &Lesson{Subject: subject, Course: courseID, Lesson: lessonID}
}
