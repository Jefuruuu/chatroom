package utils

import (
	"chatroom/services/storage/tokens"
	"chatroom/services/storage/users"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		tokenR   tokens.ITokenRepo
		userName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateToken(tt.args.tokenR, tt.args.userName); got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUserNameExist(t *testing.T) {
	type args struct {
		userLoginR users.IUserLoginRepo
		userName   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUserNameExist(tt.args.userLoginR, tt.args.userName); got != tt.want {
				t.Errorf("CheckUserNameExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
