package bindings

type ApplicationURLs struct {
	LandingPage  string
	HomePage     string
	LoginPage    string
	RegisterPage string
}

var URLs = ApplicationURLs{
	LandingPage:  "/",
	HomePage:     "/home",
	LoginPage:    "/login",
	RegisterPage: "/register",
}
