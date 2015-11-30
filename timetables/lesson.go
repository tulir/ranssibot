package timetables

import (
	"fmt"
	"maunium.net/go/ranssibot/lang"
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

// ToString converts a lesson to a string.
func (lesson Lesson) ToString() string {
	if lesson.Course == 0 || lesson.Lesson == 0 {
		return fmt.Sprintf(lang.Translate("lesson-format.noncoursed"), lesson.Subject.Name, lesson.Subject.ShortName)
	}
	return fmt.Sprintf(lang.Translate("lesson-format.coursed"), lesson.Subject.Name, lesson.Subject.ShortName, lesson.Course, lesson.Lesson)
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

	/*state := 0
	var subj string
	var courseID, lessonID int
	for index, char := range str {
		charStr := string(char)
		if state == 0 {
			if index < 2 {
				if !unicode.IsLetter(char) {
					return Lesson{}, errors.New("Error: The first two charcaters of a subject string are always letters.")
				}
				subj += charStr
			} else {
				if unicode.IsLetter(char) {
					continue
				} else if unicode.IsSpace(char) {
					state++
				} else if unicode.IsDigit(char) {
					i, err := strconv.Atoi(charStr)
					if err != nil {
						return Lesson{}, err
					}
					courseID = i
					state += 2
				} else {
					return Lesson{}, errors.New("Error")
				}
			}
		} else if state == 1 || state == 2 || state == 3 {
			if unicode.IsSpace(char) {
				continue
			} else if unicode.IsDigit(char) {
				i, err := strconv.Atoi(charStr)
				if err != nil {
					return Lesson{}, err
				}
				if state == 1 {
					courseID = i
				} else {
					if lessonID != 0 {
						lessonID = lessonID*10 + i
					} else {
						lessonID = i
					}
				}
				state++
			} else if unicode.IsLetter(char) {
				if strings.EqualFold(charStr, "a") {
					state = 5
				}
			} else {
				break
			}
		} else if state >= 5 {
			var checkfor rune
			switch state {
			case 5:
				checkfor = 'l'
			case 6:
				checkfor = 'k'
			case 7:
				fallthrough
			case 8:
				checkfor = 'a'
			case 9:
				lessonID = 1
			}
			if char == unicode.ToLower(checkfor) || char == unicode.ToUpper(checkfor) {
				state++
			} else {
				break
			}
		} else if state == 4 {
			break
		}
	}
	print(str)
	print(" -> ")
	print(subj)
	print(" - ")
	print(lessonID)
	print(" - ")
	print(courseID)
	print("\n")
	if lessonID == 0 && courseID != 0 {
		lessonID = courseID
		courseID = 1
	} else if courseID == 0 {
		courseID = 1
	}

	for name, value := range subjects {
		if strings.EqualFold(name, subj) {
			return Lesson{Subject: value, Course: courseID, Lesson: lessonID}, nil
		}
	}
	return Lesson{Subject: Subject{Name: subj, ShortName: subj}, Course: courseID, Lesson: lessonID}, nil*/
}
