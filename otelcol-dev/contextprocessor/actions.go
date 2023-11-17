package contextprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/client"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// The ExeContext the current context and the attributes
type ExeContext struct {
	ctx           context.Context
	cliInfo       client.Info
	resourceAttrs pcommon.Map
	newMetadata   map[string][]string
}

// `NewExeContext` constructs an empty ExeContext
func NewExeContext() *ExeContext {
	newCtx := context.Background()
	return &ExeContext{
		ctx:           newCtx,
		cliInfo:       client.FromContext(newCtx),
		resourceAttrs: pcommon.NewMap(),
		newMetadata:   make(map[string][]string),
	}
}

// Sets a new context
func (exc *ExeContext) SetContext(ctx context.Context, attrs pcommon.Map) {
	exc.ctx = ctx
	exc.cliInfo = client.FromContext(ctx)
	exc.resourceAttrs = attrs
	// initialize the context with the keys of the current
	exc.newMetadata = make(map[string][]string)
}

func (exc *ExeContext) GetContext() context.Context {
	return client.NewContext(exc.ctx,
		client.Info{
			Metadata: client.NewMetadata(exc.newMetadata),
		})
}

func (exc *ExeContext) GenerateAction(action ActionConfig) (Action, error) {
	base := baseAction{
		exeContext: exc,
		key:        *action.Key,
	}
	valueDefault := ""
	if action.ValueDefault != nil {
		valueDefault = *action.ValueDefault
	}
	fromAttribute := ""
	if action.FromAttribute != nil {
		fromAttribute = *action.FromAttribute
	}
	switch action.Action {
	case INSERT:
		return &ActionInsert{
			baseAction: base,
			value:      valueDefault,
			fromAttr:   fromAttribute,
		}, nil
	case UPSERT:
		return &ActionUpsert{
			baseAction: base,
			value:      valueDefault,
			fromAttr:   fromAttribute,
		}, nil
	case UPDATE:
		return &ActionUpdate{
			baseAction: base,
			value:      valueDefault,
			fromAttr:   fromAttribute,
		}, nil
	case DELETE:
		return &ActionDelete{
			baseAction: base,
		}, nil
	default:
		return nil, fmt.Errorf("unknown action type")
	}
}

// Actions
type Action interface {
	execute()
}

type baseAction struct {
	exeContext *ExeContext
	key        string
}

func (a *baseAction) getAttrKey(key, def string) (string, bool) {
	value := def
	v, exists := a.exeContext.resourceAttrs.Get(key)
	if exists {
		switch v.Type() {
		case pcommon.ValueTypeStr:
			value = v.Str()
		default:
			value = v.AsString()
		}
	}
	return value, exists
}

func (a *baseAction) getCxtKey() ([]string, bool) {
	if v, exists := a.exeContext.newMetadata[a.key]; exists {
		return v, exists
	} else {
		value := a.exeContext.cliInfo.Metadata.Get(a.key)
		return value, len(value) != 0
	}
}

// Concrete actions

type ActionInsert struct {
	baseAction
	value    string
	fromAttr string
}

func (a *ActionInsert) execute() {
	value := []string{a.value}
	if len(a.fromAttr) > 0 {
		value[0], _ = a.getAttrKey(a.fromAttr, a.value)
	}
	if _, exists := a.getCxtKey(); !exists {
		a.exeContext.newMetadata[a.key] = value
	}
}

type ActionUpsert struct {
	baseAction
	value    string
	fromAttr string
}

func (a *ActionUpsert) execute() {
	value := []string{a.value}
	if len(a.fromAttr) > 0 {
		value[0], _ = a.getAttrKey(a.fromAttr, a.value)
	}
	a.exeContext.newMetadata[a.key] = value
}

type ActionUpdate struct {
	baseAction
	value    string
	fromAttr string
}

func (a *ActionUpdate) execute() {
	value := a.value
	if len(a.fromAttr) > 0 {
		value, _ = a.getAttrKey(a.fromAttr, a.value)
	}
	if v, exists := a.getCxtKey(); exists {
		a.exeContext.newMetadata[a.key] = append(v, value)
	}
}

type ActionDelete struct {
	baseAction
}

func (a *ActionDelete) execute() {
	delete(a.exeContext.newMetadata, a.key)
}

///

type ExeActionsRunner struct {
	actions    []Action
	exeContext *ExeContext
}

func NewExeActionsRunner() *ExeActionsRunner {
	return &ExeActionsRunner{
		actions:    make([]Action, 0),
		exeContext: NewExeContext(),
	}
}

func (exr *ExeActionsRunner) AddAction(action ActionConfig) error {
	a, err := exr.exeContext.GenerateAction(action)
	if err == nil {
		exr.actions = append(exr.actions, a)
	}
	return err
}

// The executeCommands method executes all the commands
// one by one
func (exr *ExeActionsRunner) Apply(ctx context.Context, attrs pcommon.Map) context.Context {
	exr.exeContext.SetContext(ctx, attrs)
	for _, a := range exr.actions {
		a.execute()
	}
	return exr.exeContext.GetContext()
}

/////////////

// exeRunner = NewExeRunner()
// for _, action := range cfg.Actions {
// 	if err := exeRunner.AddAction(action); err != nil {
// 		return err
// 	}
// }

// // On the processor
// newCtx = exeRunner.Apply(ctx, attrs)
