package movo

import (
    "os"
    )

func init() {
}

type C1Person struct {
}

type C1ActDetail struct {
}

func GetNomRoll(actID uint, session string) ([]C1Person, os.Error)

func GetActDetail(actID uint, session string) (C1ActDetail, os.Error)

//TODO - use correct return type
func GetPersDOB(serviceNo, session string) (int, os.Error)

func C1Login(usr, pwd string) (string, os.Error)
