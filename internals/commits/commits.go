package commits

import "fmt"

type CommitType struct {
	Type    string
	Meaning string
}

func (t CommitType) String() string {
	return fmt.Sprintf("%s: %s", t.Type, t.Meaning)
}

type Commit struct {
	Type                       string
	Scope                      string
	BreakingChanges            bool
	BreakingChangesDescription string
	Title                      string
	Description                string
}

func (c Commit) String() string {
	myString := c.Type
	if c.Scope != "" {
		myString += fmt.Sprintf("(%s)", c.Scope)
	}
	myString += ": " + c.Title
	if c.Description != "" {
		myString += "\n\n" + c.Description
	}

	if c.BreakingChanges {
		myString += "\n\nBREAKING CHANGE: " + c.BreakingChangesDescription
	}

	return myString

}
