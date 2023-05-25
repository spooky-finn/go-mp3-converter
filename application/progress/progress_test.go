package progress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPercentage(t *testing.T) {
	p := New()

	p.Send(DownloadStage, 10)

	progEvent := <-p.Ch

	assert.Equal(t, 15, progEvent)

}
