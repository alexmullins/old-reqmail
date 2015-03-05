package main

// Interface defining sending updates
type Updater interface {
	Run() // intended to be run in goroutine.
}

// Type implementing Updater interface. Sends emails to buyers
// imforming them of new requisitions.
type EmailUpdater struct {
	Source    *ReqSource
	Repo      *ReqRepository
	EmailAddr string
	EmailPass string
	SMTPAddr  string
}

func (e *EmailUpdater) Run() {
	// TODO
	return
}
