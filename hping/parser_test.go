package hping

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    *HpingData
		wantErr bool
	}{
		{
			"Expected output",
			args{`HPING server (eth0 172.29.0.2): icmp mode set, 28 headers + 0 data bytes
len=40 ip=172.29.0.2 ttl=63 id=61414 icmp_seq=0 rtt=118.8 ms
ICMP timestamp: Originate=86023274 Receive=86023387 Transmit=86023387
ICMP timestamp RTT tsrtt=119


--- server hping statistic ---
1 packets tramitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 118.8/118.8/118.8 ms`},
			&HpingData{ForwardDelay: 113, ReverseDelay: 6},
			false,
		},
		{
			"Spurious sample",
			args{`HPING router (eth0 172.28.0.3): icmp mode set, 28 headers + 0 data bytes
len=40 ip=172.28.0.3 ttl=64 id=60915 icmp_seq=0 rtt=20.1 ms
ICMP timestamp: Originate=73989970 Receive=73989970 Transmit=73989970
ICMP timestamp RTT tsrtt=20


--- router hping statistic ---
1 packets tramitted, 1 packets received, 0% packet loss
round-trip min/avg/max = 20.1/20.1/20.1 ms`},
			nil,
			true,
		},
		{
			"Lost packet",
			args{`HPING router (eth0 172.28.0.3): icmp mode set, 28 headers + 0 data bytes

--- router hping statistic ---
1 packets tramitted, 0 packets received, 100% packet loss
round-trip min/avg/max = 0.0/0.0/0.0 ms`},
			nil,
			true,
		},
		{
			"Empty",
			args{``},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
