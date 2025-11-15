package sigint

import (
	"os"

	"golang.org/x/sys/windows"
)

func SendInterrupt(proc *os.Process) error {
	err := windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(proc.Pid))
	if err != nil {
		return &os.SyscallError{Syscall: "GenerateConsoleCtrlEvent", Err: err}
	}
	return nil
}
