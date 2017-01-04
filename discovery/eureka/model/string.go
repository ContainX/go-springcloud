package model

import (
	"bytes"
	"strconv"
)

func (a *Application) String() string {
	b := &FormattedBuffer{}
	b.append("")
	b.append("Name      : ", a.Name)
	b.append("Instances : [")
	for _, i := range a.Instances {
		b.append("")
		b.writeInstance(i, "  ")
	}
	b.append("]")
	return b.String()
}

func (i *Instance) String() string {
	b := &FormattedBuffer{}
	b.writeInstance(i, "")
	return b.String()
}

func (f *FormattedBuffer) writeInstance(i *Instance, indent string) {
	f.append(indent, "InstanceId   : ", i.InstanceId)
	f.append(indent, "Hostname     : ", i.HostName)
	f.append(indent, "IpAddr       : ", i.IpAddr)
	f.append(indent, "Port         : ", strconv.Itoa(i.Port.Number))
	f.append(indent, "SecurePort   : ", strconv.Itoa(i.SecurePort.Number))
	f.append(indent, "Status       : ", string(i.Status))
}

type FormattedBuffer struct {
	bytes.Buffer
}

func (f *FormattedBuffer) append(elements ...string) {
	for _, e := range elements {
		f.WriteString(e)
	}
	f.WriteString("\n")
}
