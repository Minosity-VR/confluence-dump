package saver

import (
	"testing"

	"github.com/Minosity-VR/confdump/internal/client"
)

func TestSaveConfluencePage(t *testing.T) {
	type args struct {
		rootPath string
		page     client.ConfluencePage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test SaveConfluencePage",
			args: args{
				rootPath: "./test_confExport",
				page: client.ConfluencePage{

					OwnerId:  "someId",
					Position: 123456,
					Id:       "1111",
					SpaceId:  "2222",
					Body: client.ConfluentPageBody{
						Storage: client.ConfluentPageBodyStorage{
							Representation: "storage",
							Value:          "<p /><p>Hello;</p>",
						},
					},
					Status: "current",
					Title:  "About Me",
					Links: client.ConfluentPageLinks{
						Editui: "/pages/resumedraft.action?draftId=1111",
						Tinyui: "/x/hoc",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := SaveConfluencePage(tt.args.rootPath, tt.args.page); (err != nil) != tt.wantErr {
				t.Errorf("SaveConfluencePage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
