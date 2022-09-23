package bots

type Bot interface {
	Initialize()
	Send(string)
}
