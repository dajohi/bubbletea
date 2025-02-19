package tea

import (
	"context"
	"io"

	"github.com/muesli/termenv"
)

// ProgramOption is used to set options when initializing a Program. Program can
// accept a variable number of options.
//
// Example usage:
//
//	p := NewProgram(model, WithInput(someInput), WithOutput(someOutput))
type ProgramOption func(*Program)

// WithContext lets you specify a context in which to run the Program. This is
// useful if you want to cancel the execution from outside. When a Program gets
// cancelled it will exit with an error ErrProgramKilled.
func WithContext(ctx context.Context) ProgramOption {
	return func(p *Program) {
		p.ctx = ctx
	}
}

// WithOutput sets the output which, by default, is stdout. In most cases you
// won't need to use this.
func WithOutput(output io.Writer) ProgramOption {
	return func(p *Program) {
		if o, ok := output.(*termenv.Output); ok {
			p.output = o
		} else {
			p.output = termenv.NewOutput(output, termenv.WithColorCache(true))
		}
	}
}

// WithInput sets the input which, by default, is stdin. In most cases you
// won't need to use this.
func WithInput(input io.Reader) ProgramOption {
	return func(p *Program) {
		p.input = input
		p.startupOptions |= withCustomInput
	}
}

// WithInputTTY opens a new TTY for input (or console input device on Windows).
func WithInputTTY() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withInputTTY
	}
}

// WithoutSignalHandler disables the signal handler that Bubble Tea sets up for
// Programs. This is useful if you want to handle signals yourself.
func WithoutSignalHandler() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withoutSignalHandler
	}
}

// WithoutCatchPanics disables the panic catching that Bubble Tea does by
// default. If panic catching is disabled the terminal will be in a fairly
// unusable state after a panic because Bubble Tea will not perform its usual
// cleanup on exit.
func WithoutCatchPanics() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withoutCatchPanics
	}
}

// WithAltScreen starts the program with the alternate screen buffer enabled
// (i.e. the program starts in full window mode). Note that the altscreen will
// be automatically exited when the program quits.
//
// Example:
//
//	p := tea.NewProgram(Model{}, tea.WithAltScreen())
//	if _, err := p.Run(); err != nil {
//	    fmt.Println("Error running program:", err)
//	    os.Exit(1)
//	}
//
// To enter the altscreen once the program has already started running use the
// EnterAltScreen command.
func WithAltScreen() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withAltScreen
	}
}

// WithMouseCellMotion starts the program with the mouse enabled in "cell
// motion" mode.
//
// Cell motion mode enables mouse click, release, and wheel events. Mouse
// movement events are also captured if a mouse button is pressed (i.e., drag
// events). Cell motion mode is better supported than all motion mode.
//
// To enable mouse cell motion once the program has already started running use
// the EnableMouseCellMotion command. To disable the mouse when the program is
// running use the DisableMouse command.
//
// The mouse will be automatically disabled when the program exits.
func WithMouseCellMotion() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withMouseCellMotion // set
		p.startupOptions &^= withMouseAllMotion // clear
	}
}

// WithMouseAllMotion starts the program with the mouse enabled in "all motion"
// mode.
//
// EnableMouseAllMotion is a special command that enables mouse click, release,
// wheel, and motion events, which are delivered regardless of whether a mouse
// button is pressed, effectively enabling support for hover interactions.
//
// Many modern terminals support this, but not all. If in doubt, use
// EnableMouseCellMotion instead.
//
// To enable the mouse once the program has already started running use the
// EnableMouseAllMotion command. To disable the mouse when the program is
// running use the DisableMouse command.
//
// The mouse will be automatically disabled when the program exits.
func WithMouseAllMotion() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withMouseAllMotion   // set
		p.startupOptions &^= withMouseCellMotion // clear
	}
}

// WithoutRenderer disables the renderer. When this is set output and log
// statements will be plainly sent to stdout (or another output if one is set)
// without any rendering and redrawing logic. In other words, printing and
// logging will behave the same way it would in a non-TUI commandline tool.
// This can be useful if you want to use the Bubble Tea framework for a non-TUI
// application, or to provide an additional non-TUI mode to your Bubble Tea
// programs. For example, your program could behave like a daemon if output is
// not a TTY.
func WithoutRenderer() ProgramOption {
	return func(p *Program) {
		p.renderer = &nilRenderer{}
	}
}

// WithANSICompressor removes redundant ANSI sequences to produce potentially
// smaller output, at the cost of some processing overhead.
//
// This feature is provisional, and may be changed or removed in a future version
// of this package.
func WithANSICompressor() ProgramOption {
	return func(p *Program) {
		p.startupOptions |= withANSICompressor
	}
}

// WithFilter supplies an event filter that will be invoked before Bubble Tea
// processes a tea.Msg. The event filter can return any tea.Msg which will then
// get handled by Bubble Tea instead of the original event. If the event filter
// returns nil, the event will be ignored and Bubble Tea will not process it.
//
// As an example, this could be used to prevent a program from shutting down if
// there are unsaved changes.
//
// Example:
//
//	func filter(m tea.Model, msg tea.Msg) tea.Msg {
//		if _, ok := msg.(tea.QuitMsg); !ok {
//			return msg
//		}
//
//		model := m.(myModel)
//		if model.hasChanges {
//			return nil
//		}
//
//		return msg
//	}
//
//	p := tea.NewProgram(Model{}, tea.WithFilter(filter));
//
//	if _,err := p.Run(); err != nil {
//		fmt.Println("Error running program:", err)
//		os.Exit(1)
//	}
func WithFilter(filter func(Model, Msg) Msg) ProgramOption {
	return func(p *Program) {
		p.filter = filter
	}
}
