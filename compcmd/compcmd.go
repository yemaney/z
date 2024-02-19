/*
Package compcmd is a completion driver for Bonzai command trees and
fulfills the bonzai.Completer package interface. See Complete method for
details.

	    Expand All

	@@ -24,35 +23,42 @@ type comp struct{}
*/
package compcmd

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/structs/set/text/set"
)

// New returns a private struct that fulfills the bonzai.Completer
// interface. See Complete method for details.
func New() *comp { return new(comp) }

type comp struct{}

// Complete resolves completion as follows:
//
//  1. If leaf has Comp function, delegate to it
//
//  2. If leaf has no arguments, return all Commands and Params
//
//  3. If first argument is the name of a Command return it, even if it
//     is in the Hidden list
//
//  4. If nargs == 1, append command names, shortcuts and params to the
//     list. Otherwise, only append params.
//
//  5. Remove hidden and duplicate params from the list, if any.
//
//  6. Return all the items in the list that HasPrefix matching the
//     first arg.
//
// See bonzai.Completer.
func (comp) Complete(x bonzai.Command, args ...string) []string {
	// if has completer, delegate
	// if c := x.GetComp(); c != nil {
	// 	return c.Complete(x, args...)
	// }

	nargs := len(args)

	// not sure we've completed the command name itself yet
	if nargs == 0 {
		return []string{x.GetName()}
	}

	//	build list of visible commands and params
	var list []string
	list = append(list, x.GetParams()...)
	if nargs == 1 {
		list = append(list, x.GetCommandNames()...)
		list = append(list, x.GetShortcuts()...)
	}

	// remove hidden and duplicate params
	min := append(x.GetHidden(), args[:nargs]...)
	list = set.Minus[string, string](list, min)

	return filt.HasPrefix(list, args[nargs-1])
}
