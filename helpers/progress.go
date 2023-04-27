package helpers

type Stage string

const DownloadStage Stage = "downloading"
const ConvertStage Stage = "converting"

type Progress struct {
	Ch chan Payload
}

type Payload struct {
	Stage  Stage  `json:"stage"`
	Prog   int    `json:"progress"`
	Status string `json:"status"`
}

func (p *Progress) Send(stage Stage, prog int) {

	p.Ch <- Payload{
		Stage:  stage,
		Prog:   prog,
		Status: "ok",
	}
}

func NewProg() *Progress {
	return &Progress{
		Ch: make(chan Payload),
	}
}
