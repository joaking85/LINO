package id

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"makeit.imfr.cgi.com/lino/pkg/id"
	"makeit.imfr.cgi.com/lino/pkg/relation"
)

var idStorage id.Storage
var relStorage relation.Storage
var idExporter id.Exporter

// SetLogger if needed, default no logger
func SetLogger(l id.Logger) {
	id.SetLogger(l)
}

// Inject dependencies
func Inject(ids id.Storage, rels relation.Storage, ex id.Exporter) {
	idStorage = ids
	relStorage = rels
	idExporter = ex
}

// NewCommand implements the cli id command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "id {create,display-plan,show-graph} [arguments ...]",
		Short:   "Manage ingress descriptor",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s id create mydatabase public.customer", fullName),
	}
	cmd.AddCommand(newCreateCommand(fullName, err, out, in))
	cmd.AddCommand(newDisplayPlanCommand(fullName, err, out, in))
	cmd.AddCommand(newShowGraphCommand(fullName, err, out, in))
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}