package walk

import (
	"context"
	"os"
	"time"

	"github.com/Tom5521/fsize/flags"
	"github.com/Tom5521/fsize/stat"
	"github.com/charmbracelet/log"
	"github.com/karrick/godirwalk"
	po "github.com/leonelquinteros/gotext"
)

func ProcessFile(file *stat.File) {
	if (!flags.Progress || flags.NoProgress) || !file.IsDir || flags.NoWalk {
		return
	}
	var warns []error

	state := &progressState{}
	ctx, cancel := context.WithCancel(context.Background())
	progressCtx, progressCancel := context.WithCancel(context.Background())

	defer progressCancel()
	defer cancel()

	go func() {
		if (flags.NoProgress || !flags.Progress) || flags.PrintOnWalk {
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

	err := godirwalk.Walk(file.Path, &godirwalk.Options{
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			err := filter(osPathname, directoryEntry)
			if err != nil {
				return err
			}
			if flags.PrintOnWalk && !flags.NoProgress {
				log.Info(po.Get("Reading \"%s\"...", osPathname))
			}

			defer func() {
				file.FilesNumber++

				state.Lock()
				state.nfiles = file.FilesNumber
				state.size = file.Size
				state.Unlock()
			}()

			stat, err := os.Stat(osPathname)
			if err != nil {
				return err
			}

			file.Size += stat.Size()

			return nil
		},
		ErrorCallback: func(s string, err error) godirwalk.ErrorAction {
			warns = append(warns, err)
			state.Lock()
			state.nwarns = len(warns)
			state.Unlock()
			return godirwalk.SkipNode
		},
		FollowSymbolicLinks: flags.FollowSymlinks,
	})
	if err != nil {
		warns = append(warns, err)
		return
	}

	cancel()
	if (!flags.NoProgress || flags.Progress) && !flags.PrintOnWalk {
		<-progressCtx.Done()
	}

	for i, err := range warns {
		if i >= flags.WarnLimit {
			log.Warn(po.Get("too many warns (%d)", len(warns)))
			break
		}
		log.Warn(err)
	}
}
