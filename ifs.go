package main

// Interface defining a source of open requisition report
type ReqSource interface {
	GetReqReport() (*ReqReport, error)
}

// Concrete implementation of ReqSource interface
type IFS struct {
}

func (i *IFSSource) GetReqReport() (*ReqReport, error) {
	return nil, nil
}

type ReqLine struct {
	ReqNo     string
	LineNo    int
	ReleaseNo int
	BuyerName string
}

type ReqReport []ReqLine
