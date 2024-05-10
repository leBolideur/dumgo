package users

type UserCommArgs struct {
	Nick string
}
type UserCommResponse struct {
	Status bool
	Token  string
	Msg    string
}
