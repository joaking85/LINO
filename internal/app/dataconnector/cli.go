// Copyright (C) 2021 CGI France
//
// This file is part of LINO.
//
// LINO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// LINO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with LINO.  If not, see <http://www.gnu.org/licenses/>.

package dataconnector

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"makeit.imfr.cgi.com/lino/pkg/dataconnector"
)

var storage dataconnector.Storage
var dataPingerFactory map[string]dataconnector.DataPingerFactory

var logger dataconnector.Logger = dataconnector.Nologger{}

// SetLogger if needed, default no logger
func SetLogger(l dataconnector.Logger) {
	logger = l
	dataconnector.SetLogger(l)
}

// Inject dependencies
func Inject(dbas dataconnector.Storage, dpf map[string]dataconnector.DataPingerFactory) {
	storage = dbas
	dataPingerFactory = dpf
}

// NewCommand implements the cli dataconnector command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dataconnector {add,list} [arguments ...]",
		Short:   "Manage database aliases",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s dataconnector add mydatabase postgresql://postgres:sakila@localhost:5432/postgres?sslmode=disable", fullName),
		Aliases: []string{"dc"},
	}
	cmd.AddCommand(newAddCommand(fullName, err, out, in))
	cmd.AddCommand(newListCommand(fullName, err, out, in))
	cmd.AddCommand(newPingCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}
