package main

import (
	"fmt"
	"github.com/docker/libcontainer/netlink"
	"os"
)

func usage() {
	fmt.Printf("tc qdisc show dev <device>")
	os.Exit(1)
}

func listQdisc(iface *net.Interface) error {
	s, err := getNetlinkSocket()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	wb := newNetlinkRequest(syscall.RTM_GETQDISC, syscall.NLM_F_DUMP)

	msg := newIfInfomsg(syscall.AF_UNSPEC)
	msg.Index = int32(iface.Index)
	wb.AddData(msg)

	if err := s.Send(wb); err != nil {
		return nil, err
	}

	pid, err := s.GetPid()
	if err != nil {
		return nil, err
	}

	res := make([]Route, 0)

outer:
	for {
		msgs, err := s.Receive()
		if err != nil {
			return nil, err
		}
		for _, m := range msgs {
			if err := s.CheckMessage(m, wb.Seq, pid); err != nil {
				if err == io.EOF {
					break outer
				}
				return nil, err
			}
			if m.Header.Type != syscall.RTM_NEWROUTE {
				continue
			}

			//handle individual messages
		}
	}

	return res, nil
}

func main() {

	args := os.Args
	if len(args) != 5 || args[1] != "qdisc" || args[2] != "show" || args[3] != "dev" {
		usage()
	}

	device := args[4]

	fmt.Printf("\ndevice = %s\n", device)
}
