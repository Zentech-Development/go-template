package bindings

type ApplicationURLs struct {
	Home     string
	Login    string
	Register string
}

var URLs = ApplicationURLs{
	Home:     "/",
	Login:    "/login",
	Register: "/register",
}
