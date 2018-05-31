package pcap

import (
	"fmt"
	"testing"
	"time"
)

func TestPcap(t *testing.T) {
	fmt.Println(PcapLookupDev())
	h, err := OpenLive("eno1", 512, true, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(h.DatalinkType())
	fmt.Println(h.SetFilter("tcp port 80"))
}
