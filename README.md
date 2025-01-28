# termflow 

**a TUI for task management**

termflow is a terminal-based task manager built in Go using bubbletea for the interface. It was inspired by [clikan](https://github.com/kitplummer/clikan) & incorporates additional functionality that I developed in [my fork of clikan](https://github.com/tomcotter7/clikan), but aims to improve the experience by providing a persistent TUI rather than requiring repeated command usage.

## Installation

### Quick Installation (installs to ~/.local/bin)
 
```bash
git clone https://github.com/yourusername/termflow.git
cd termflow
make install
```
This creates a `termflow` binary in your ~/.local/bin folder.

## Usage

There are a variety of keyboard shortcuts available for use:

- (`a`)dd: Add a new task to the current column
- (`p`)romote: Move task forward 1 column
- (`r`)egress: Move task backward 1 column
- (`d`)elete: Delete task from the list
- (`e`)dit: Edit the underlying task details
- (`s`)how: Show all details about a task
- (`t`)oday: Set a task to be due today
- (`b`)locked: Set a task to be blocked

`q` / `esc` is always to exit a screen.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
