package extract

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xo/dburl"

	"makeit.imfr.cgi.com/lino/pkg/dataconnector"
	"makeit.imfr.cgi.com/lino/pkg/extract"
	"makeit.imfr.cgi.com/lino/pkg/id"
	"makeit.imfr.cgi.com/lino/pkg/relation"
	"makeit.imfr.cgi.com/lino/pkg/table"
)

var dataconnectorStorage dataconnector.Storage
var relStorage relation.Storage
var tabStorage table.Storage
var idStorage id.Storage
var dataSourceFactories map[string]extract.DataSourceFactory
var extractExporter extract.RowExporter

// local flags
var limit uint

var logger extract.Logger

// SetLogger if needed, default no logger
func SetLogger(l extract.Logger) {
	logger = l
	extract.SetLogger(l)
}

// Inject dependencies
func Inject(dbas dataconnector.Storage, rs relation.Storage, ts table.Storage, ids id.Storage, dsfmap map[string]extract.DataSourceFactory, rowExporter extract.RowExporter) {
	dataconnectorStorage = dbas
	relStorage = rs
	tabStorage = ts
	idStorage = ids
	dataSourceFactories = dsfmap
	extractExporter = rowExporter
}

// NewCommand implements the cli extract command
func NewCommand(fullName string, err *os.File, out *os.File, in *os.File) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "extract [DB Alias Name]",
		Short:   "Extract data from a database",
		Long:    "",
		Example: fmt.Sprintf("  %[1]s extract mydatabase --limit 1", fullName),
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			datasource, e1 := getDataSource(args[0])
			if e1 != nil {
				fmt.Fprintln(err, e1.Error())
				os.Exit(1)
			}

			plan, e2 := getExtractionPlan()
			if e2 != nil {
				fmt.Fprintln(err, e2.Error())
				os.Exit(1)
			}

			e3 := extract.Extract(plan, datasource, extractExporter)
			if e3 != nil {
				fmt.Fprintln(err, e3.Error())
				os.Exit(1)
			}
		},
	}
	cmd.Flags().UintVarP(&limit, "limit", "l", 1, "limit the number of results")
	cmd.SetOut(out)
	cmd.SetErr(err)
	cmd.SetIn(in)
	return cmd
}

func getDataSource(dataconnectorName string) (extract.DataSource, *extract.Error) {
	alias, e1 := dataconnector.Get(dataconnectorStorage, dataconnectorName)
	if e1 != nil {
		return nil, &extract.Error{Description: e1.Error()}
	}

	u, e2 := dburl.Parse(alias.URL)
	if e2 != nil {
		return nil, &extract.Error{Description: e2.Error()}
	}

	datasourceFactory, ok := dataSourceFactories[u.Unaliased]
	if !ok {
		return nil, &extract.Error{Description: "no datasource found for database type"}
	}

	return datasourceFactory.New(alias.URL), nil
}

func getExtractionPlan() (extract.Plan, *extract.Error) {
	ep, err1 := id.GetExtractionPlan(idStorage)
	if err1 != nil {
		return nil, &extract.Error{Description: err1.Error()}
	}

	relations, err2 := relStorage.List()
	if err2 != nil {
		return nil, &extract.Error{Description: err2.Error()}
	}

	tables, err3 := tabStorage.List()
	if err3 != nil {
		return nil, &extract.Error{Description: err3.Error()}
	}

	filter := extract.NewFilter(limit, extract.Row{})
	stepList := getStepList(ep, relations, tables)

	return extract.NewPlan(filter, stepList), nil
}

func getStepList(ep id.ExtractionPlan, relations []relation.Relation, tables []table.Table) extract.StepList {
	rmap := map[string]relation.Relation{}
	for _, relation := range relations {
		rmap[relation.Name] = relation
	}

	tmap := map[string]table.Table{}
	for _, table := range tables {
		tmap[table.Name] = table
	}

	smap := []id.Step{}
	for idx := uint(0); idx < ep.Len(); idx++ {
		smap = append(smap, ep.Step(idx))
	}

	logger.Debug(fmt.Sprintf("there is %v step(s) to build", ep.Len()))
	defer logger.Debug(fmt.Sprintf("finished building %v step(s) with success", ep.Len()))

	converter := epToStepListConverter{
		rmap:   rmap,
		tmap:   tmap,
		smap:   smap,
		exrmap: map[string]extract.Relation{},
		extmap: map[string]extract.Table{},
		exsmap: map[uint]extract.Step{},
	}

	return converter.getSteps()
}

type epToStepListConverter struct {
	rmap map[string]relation.Relation
	tmap map[string]table.Table
	smap []id.Step

	exrmap map[string]extract.Relation
	extmap map[string]extract.Table
	exsmap map[uint]extract.Step
}

func (c epToStepListConverter) getTable(name string) extract.Table {
	if extable, ok := c.extmap[name]; ok {
		return extable
	}

	table, ok := c.tmap[name]
	if !ok {
		logger.Error(fmt.Sprintf("missing table %v in tables.yaml", name))
		return extract.NewTable(name, "")
	}

	logger.Trace(fmt.Sprintf("building table %v", table))

	return extract.NewTable(table.Name, table.Keys[0]) // TODO : support multivalued primary keys
}

func (c epToStepListConverter) getRelation(name string) extract.Relation {
	if exrelation, ok := c.exrmap[name]; ok {
		return exrelation
	}

	relation, ok := c.rmap[name]
	if !ok {
		logger.Error(fmt.Sprintf("missing relation %v in relations.yaml", name))
		return extract.NewRelation(name, nil, nil, "", "")
	}

	logger.Trace(fmt.Sprintf("building relation %v", relation))

	return extract.NewRelation(
		relation.Name,
		c.getTable(relation.Parent.Name),
		c.getTable(relation.Child.Name),
		relation.Parent.Keys[0], // TODO : support multivalued keys
		relation.Child.Keys[0],  // TODO : support multivalued keys
	)
}

func (c epToStepListConverter) getRelationList(relations id.IngressRelationList) extract.RelationList {
	exrelations := []extract.Relation{}
	for idx := uint(0); idx < relations.Len(); idx++ {
		exrelations = append(exrelations, c.getRelation(relations.Relation(idx).Name()))
	}
	return extract.NewRelationList(exrelations)
}

func (c epToStepListConverter) getCycleList(cycles id.CycleList) extract.CycleList {
	excycles := []extract.Cycle{}
	for idx := uint(0); idx < cycles.Len(); idx++ {
		excycles = append(excycles, c.getRelationList(cycles.Cycle(idx)))
	}
	return extract.NewCycleList(excycles)
}

func (c epToStepListConverter) getStepList(previousStep uint) extract.StepList {
	exsteps := []extract.Step{}
	for _, step := range c.smap {
		if step.PreviousStep() == previousStep {
			exsteps = append(exsteps, c.getStep(step.Index()))
		}
	}
	return extract.NewStepList(exsteps)
}

func (c epToStepListConverter) getStep(idx uint) extract.Step {
	if exstep, ok := c.exsmap[idx]; ok {
		return exstep
	}

	step := c.smap[idx-1]

	logger.Trace(fmt.Sprintf("building %v", step))

	var exstep extract.Step
	if step.Index() > 1 {
		exstep = extract.NewStep(
			step.Index(),
			c.getTable(step.Entry().Name()),
			c.getRelation(step.Following().Name()),
			c.getRelationList(step.Relations()),
			c.getCycleList(step.Cycles()),
			c.getStepList(step.Index()),
		)
	} else {
		exstep = extract.NewStep(
			step.Index(),
			c.getTable(step.Entry().Name()),
			nil,
			c.getRelationList(step.Relations()),
			c.getCycleList(step.Cycles()),
			c.getStepList(step.Index()),
		)
	}

	c.exsmap[idx] = exstep

	return exstep
}

func (c epToStepListConverter) getSteps() extract.StepList {
	exsteps := []extract.Step{}
	for _, step := range c.smap {
		exsteps = append(exsteps, c.getStep(step.Index()))
	}
	return extract.NewStepList(exsteps)
}