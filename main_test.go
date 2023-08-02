package main

import (
    "testing"
    "regexp"
)

func TestQuit(t *testing.T) {
    command := "quit"
    want := regexp.MustCompile(`\b`+command+`\b`)
    msg, err := input("quit")
    if !want.MatchString(msg) || err != nil {
        t.Fatalf(`input("quit") = %q, %v, want match for %#q, nil`, msg, err, want)
    }
}
