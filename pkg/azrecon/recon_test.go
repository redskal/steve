/*
Released under YOLO licence. Idgaf what you do.
*/

package azrecon

import (
	"reflect"
	"testing"
)

func TestCheckAzureCnames(t *testing.T) {
	domain, err := CheckAzureCnames("learn.microsoft.com", "1.1.1.1:53")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(domain)
}

func TestCheckResourceExists(t *testing.T) {
	type args struct {
		resource string
		resolver string
	}
	tests := []struct {
		name    string
		args    args
		want    []Resource
		wantErr bool
	}{
		{
			"learn-public",
			args{resource: "learn-public", resolver: "8.8.8.8:53"},
			[]Resource{
				{Domain: "learn-public.trafficmanager.net", Type: "Traffic Manger (Load Balancer)"},
			},
			false,
		},
		{
			"npplatformportal (may fail to ordering)",
			args{resource: "npplatformportal", resolver: "8.8.8.8:53"},
			[]Resource{
				{Domain: "npplatformportal.scm.azurewebsites.net", Type: "App Services (Management)"},
				{Domain: "npplatformportal.azurewebsites.net", Type: "App Services"},
			},
			false,
		},
		{
			"button mash",
			args{resource: "hsdbikhvsikbkjhsdbvfkhsdjbksfjdhebf", resolver: "8.8.8.8:53"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckResourceExists(tt.args.resource, tt.args.resolver)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckResourceExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckResourceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
