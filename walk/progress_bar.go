package walk

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Tom5521/fsize/flags"
	"github.com/dustin/go-humanize"
	po "github.com/leonelquinteros/gotext"
	"github.com/schollz/progressbar/v3"
)

type progressState struct {
	sync.RWMutex
	nwarns int
	nfiles int64
	size   int64
}

func progressUpdater(
	ctx context.Context,
	cancel context.CancelFunc,
	state *progressState,
) {
	bar := progressbar.NewOptions64(-1,
		progressbar.OptionFullWidth(),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetRenderBlankState(false),
		progressbar.OptionSetSpinnerChangeInterval(time.Millisecond*50),
	)
	startedTime := time.Now().Add(-flags.ProgressDelay)
	descriptionFormat := "%s [%v]%s"

	lastDotUpdate := time.Now()

	var dots string

	set := func() {
		state.RLock()
		nfiles := state.nfiles
		nwarns := state.nwarns
		size := state.size
		state.RUnlock()

		text := fmt.Sprintf(descriptionFormat,
			po.Get(
				"%d files, %d errors, total size: %s",
				nfiles,
				nwarns,
				humanize.IBytes(uint64(size)),
			),
			time.Since(startedTime).Truncate(time.Second),
			dots,
		)
		bar.Describe(text)
	}
	for {
		select {
		case <-time.After(time.Second / 3):
			if time.Since(lastDotUpdate) >= time.Second {
				if len(dots) >= 6 {
					dots = ""
				}
				dots += "."
				lastDotUpdate = time.Now()
			}
			set()
		case <-ctx.Done():
			dots = ""
			set()
			bar.Finish()
			if !flags.NotClearBar {
				bar.Clear()
			}
			fmt.Fprintf(os.Stderr, "\n")
			cancel()
			return
		}
	}
}
