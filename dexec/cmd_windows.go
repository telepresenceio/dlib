package dexec

import "golang.org/x/sys/windows"

func (c *Cmd) canInterrupt() bool {
	return c != nil &&
		c.Cmd != nil &&
		c.SysProcAttr != nil &&
		(c.SysProcAttr.CreationFlags&windows.CREATE_NEW_PROCESS_GROUP) != 0
}
