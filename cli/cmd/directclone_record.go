package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/transport"
)

func populateDirectCloneURLs(rec *model.ScanRecord, url string) {
	rec.Transport = transport.Classify(url)
	if rec.Transport == constants.ScanTransportSSH {
		populateSSHCloneURL(rec, url)

		return
	}

	if rec.Transport == constants.ScanTransportHTTPS {
		populateHTTPSCloneURL(rec, url)

		return
	}

	rec.HTTPSUrl = url
}

func populateSSHCloneURL(rec *model.ScanRecord, url string) {
	rec.SSHUrl = url
	if httpsURL, ok := ConvertURLToHTTPS(url); ok {
		rec.HTTPSUrl = httpsURL
	}
}

func populateHTTPSCloneURL(rec *model.ScanRecord, url string) {
	rec.HTTPSUrl = url
	if sshURL, ok := ConvertURLToSSH(url); ok {
		rec.SSHUrl = sshURL
	}
}
