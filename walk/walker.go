package walk

import (
	"context"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/stat"
	"github.com/charmbracelet/log"
	po "github.com/leonelquinteros/gotext"
)

func ProcessFile(file *stat.File) {
	if !flags.Progress && !file.IsDir && flags.NoWalk {
		return
	}
	var warns []error

	state := &progressState{}
	ctx, cancel := context.WithCancel(context.Background())
	progressCtx, progressCancel := context.WithCancel(context.Background())

	defer progressCancel()
	defer cancel()

	go func() {
		if flags.NoProgress || flags.PrintOnWalk {
			return
		}
		select {
		case <-time.After(flags.ProgressDelay):
			progressUpdater(ctx, progressCancel, state)
		case <-ctx.Done():
			progressCancel()
			return
		}
	}()

	err := filepath.Walk(file.Path,
		func(path string, info fs.FileInfo, err error) error {
			err = filter(path, info, err)
			if err != nil {
				switch err {
				case filepath.SkipDir, filepath.SkipAll:
					return err
				default:
					warns = append(warns, err)
					state.Lock()
					state.nwarns = len(warns)
					state.Unlock()
					return nil
				}
			}
			file.FilesNumber++
			file.Size += info.Size()

			state.Lock()
			state.nfiles = file.FilesNumber
			state.size = uint64(file.Size)
			state.Unlock()

			if flags.PrintOnWalk && !flags.NoProgress {
				log.Info(po.Get("Reading \"%s\"...", path))
			}

			return nil
		},
	)
	cancel()
	if !flags.NoProgress && !flags.PrintOnWalk {
		<-progressCtx.Done()
	}

	if err != nil {
		warns = append(warns, err)
	}

	for i, err2 := range warns {
		if i >= flags.WarnLimit {
			log.Warn(po.Get("too many warns (%d)", len(warns)))
			break
		}
		log.Warn(err2.Error())
	}
}
