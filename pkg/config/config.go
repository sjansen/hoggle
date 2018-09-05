package config

type GitConfig interface {
	Get(name string) (string, error)
	Set(name, value string) error
	Unset(name, value_regex string) error
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

func Uninstall(git GitConfig) (err error) {
	/* Set standalonetransferagent first so that, if an error occurs,
	   we're less likely to leave git in a bad state.
	*/
	if err = git.Unset("lfs.standalonetransferagent", "^hoggle$"); err != nil {
		return
	}

	if err = git.Unset("lfs.customtransfer.hoggle.args", ""); err != nil {
		return
	}

	if err = git.Unset("lfs.customtransfer.hoggle.path", ""); err != nil {
		return
	}

	return git.Unset("lfs.customtransfer.hoggle.concurrent", "")
}
