// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly
// prohibited. Proprietary and confidential
//
// @doc
//
// @end
// Created :   16 September 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package routes

import (
	"bitbucket.org/nanobox/na-api"
)

func Init() {
	api.Router.Get("/files", api.TraceRequest(listFileNames))
	api.Router.Get("/files/{file}", api.TraceRequest(getFileContents))
	api.Router.Get("/branches", api.TraceRequest(showBranches))
	api.Router.Get("/branches/{branch}", api.TraceRequest(showBranchDetails))
	api.Router.Get("/commits", api.TraceRequest(showCommits))
	api.Router.Get("/commits/commit", api.TraceRequest(showCommitDetails))
}
