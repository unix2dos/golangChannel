package main

func makeThumb(files []string) (thumbfiles []string, err error) {

	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(files))

	for _, f := range files {
		go func(f string) {
			var it item
			//it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range files {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return
}
