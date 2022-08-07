package shell

import "testing"

func Test_parseUrlDetails(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name       string
		args       args
		wantTenant string
		wantId     string
		wantErr    bool
	}{
		{
			name:       "https with username",
			args:       args{"https://my-name@github.com/myteam/myrepo.git"},
			wantTenant: "myteam", wantId: "myrepo", wantErr: false,
		},
		{
			name:       "https without username",
			args:       args{"https://github.com/myteam/myrepo.git"},
			wantTenant: "myteam", wantId: "myrepo", wantErr: false,
		},
		{
			name:       "ssh with username",
			args:       args{"git@github.com:myteam/myrepo.git"},
			wantTenant: "myteam", wantId: "myrepo", wantErr: false,
		},
		{
			name:       "ssh without username",
			args:       args{"ssh://github.com:myteam/myrepo.git"},
			wantTenant: "myteam", wantId: "myrepo", wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tenant, id, err := parseUrlDetails(tt.args.url)
				if (err != nil) != tt.wantErr {
					t.Errorf("parseUrlDetails() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tenant != tt.wantTenant {
					t.Errorf("parseUrlDetails() tenant = %v, wantTenant %v", tenant, tt.wantTenant)
				}
				if id != tt.wantId {
					t.Errorf("parseUrlDetails() id = %v, wantId %v", id, tt.wantId)
				}
			},
		)
	}
}
