package stat

import "os/user"

type FileIDs struct {
	User  *user.User
	Group *user.Group

	SupportFileIDs bool
}
