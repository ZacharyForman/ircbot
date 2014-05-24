package main

import "github.com/ZacharyForman/bot"            //the actual bot
import "github.com/ZacharyForman/modules/fmtio"  //input/output funcs.

import "github.com/ZacharyForman/modules/admin"  //admin, e.g. kick, ban
import "github.com/ZacharyForman/modules/regex"  //regex funcs. s/blah/blahblah/
import "github.com/ZacharyForman/modules/npost"  //New posts to specified subs
import "github.com/ZacharyForman/modules/timer"  //timer module

func main() {
    b := bot.NewBot(".bot")
    b.Register(input.NewFmtIOModule(".fmtio"))
    b.Register(admin.NewAdminModule(".admin"))
    b.Register(regex.NewRegexModule(".regex"))
    b.Register(npost.NewNpostModule(".npost"))
    b.Register(timer.NewTimerModule(".timer"))
    b.Run()
}