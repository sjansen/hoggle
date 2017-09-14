package config

type GitConfig interface {
	Get(name string) (string, error)
	Set(name, value string) error
}

func Init(git GitConfig, uri string) (err error) {
	// TODO don't reset standalonetransferagent without --force

	if err = git.Set("lfs.customtransfer.hoggle.args", uri); err != nil {
		return
	}

	if err = git.Set("lfs.customtransfer.hoggle.path", "hoggle"); err != nil {
		return
	}

	if err = git.Set("lfs.customtransfer.hoggle.concurrent", "false"); err != nil {
		return
	}

	/* Set standalonetransferagent last so that, if an error occurs,
	   we're less likely to leave git in a bad state.
	*/
	return git.Set("lfs.standalonetransferagent", "hoggle")
}
