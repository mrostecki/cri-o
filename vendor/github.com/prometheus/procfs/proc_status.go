// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package procfs

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// ProcStat provides status information about the process,
// read from /proc/[pid]/stat.
type ProcStatus struct {
	// The process ID.
	PID int
	// The process name.
	Name string

	// Peak virtual memory size.
	VmPeak uint64
	// Virtual memory size.
	VmSize uint64
	// Locked memory size.
	VmLck uint64
	// Pinned memory size.
	VmPin uint64
	// Peak resident set size.
	VmHWM uint64
	// Resident set size (sum of RssAnnon RssFile and RssShmem).
	VmRSS uint64
	// Size of resident anonymous memory.
	RssAnon uint64
	// Size of resident file mappings.
	RssFile uint64
	// Size of resident shared memory.
	RssShmem uint64
	// Size of data segments.
	VmData uint64
	// Size of stack segments.
	VmStk uint64
	// Size of text segments.
	VmExe uint64
	// Shared library code size.
	VmLib uint64
	// Page table entries size.
	VmPTE uint64
	// Size of second-level page tables.
	VmPMD uint64
	// Swapped-out virtual memory size by anonymous private.
	VmSwap uint64
	// Size of hugetlb memory portions
	HugetlbPages uint64

	// Number of voluntary context switches.
	VoluntaryCtxtSwitches uint64
	// Number of involuntary context switches.
	NonVoluntaryCtxtSwitches uint64
}

// NewStatus returns the current status information of the process.
func (p Proc) NewStatus() (ProcStatus, error) {
	f, err := os.Open(p.path("status"))
	if err != nil {
		return ProcStatus{}, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ProcStatus{}, err
	}

	s := ProcStatus{PID: p.PID}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !bytes.Contains([]byte(line), []byte(":")) {
			continue
		}

		kv := strings.SplitN(line, ":", 2)

		// removes spaces
		k := string(strings.TrimSpace(kv[0]))
		v := string(strings.TrimSpace(kv[1]))
		// removes "kB"
		v = string(bytes.Trim([]byte(v), " kB"))

		// value to int when possible
		// we can skip error check here, 'cause vKBytes is not used when value is a string
		vKBytes, _ := strconv.ParseUint(v, 10, 64)
		// convert kB to B
		vBytes := vKBytes * 1024

		s.fillStatus(k, v, vKBytes, vBytes)
	}

	return s, nil
}

func (s *ProcStatus) fillStatus(k string, vString string, vUint uint64, vUintBytes uint64) {
	switch k {
	case "Name":
		s.Name = vString
	case "VmPeak":
		s.VmPeak = vUintBytes
	case "VmSize":
		s.VmSize = vUintBytes
	case "VmLck":
		s.VmLck = vUintBytes
	case "VmPin":
		s.VmPin = vUintBytes
	case "VmHWM":
		s.VmHWM = vUintBytes
	case "VmRSS":
		s.VmRSS = vUintBytes
	case "RssAnon":
		s.RssAnon = vUintBytes
	case "RssFile":
		s.RssFile = vUintBytes
	case "RssShmem":
		s.RssShmem = vUintBytes
	case "VmData":
		s.VmData = vUintBytes
	case "VmStk":
		s.VmStk = vUintBytes
	case "VmExe":
		s.VmExe = vUintBytes
	case "VmLib":
		s.VmLib = vUintBytes
	case "VmPTE":
		s.VmPTE = vUintBytes
	case "VmPMD":
		s.VmPMD = vUintBytes
	case "VmSwap":
		s.VmSwap = vUintBytes
	case "HugetlbPages":
		s.HugetlbPages = vUintBytes
	case "voluntary_ctxt_switches":
		s.VoluntaryCtxtSwitches = vUint
	case "nonvoluntary_ctxt_switches":
		s.NonVoluntaryCtxtSwitches = vUint
	}
}

// TotalCtxtSwitches returns the total context switch.
func (s ProcStatus) TotalCtxtSwitches() uint64 {
	return s.VoluntaryCtxtSwitches + s.NonVoluntaryCtxtSwitches
}
