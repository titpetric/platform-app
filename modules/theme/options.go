package theme

type Options struct {
	Header []Link
}

type Link struct {
	Name, URL string
}

var options = &Options{
	Header: []Link{
		{
			Name: "Home",
			URL:  "/",
		},
		{
			Name: "Login",
			URL:  "/login",
		},
		{
			Name: "Logout",
			URL:  "/logout",
		},
	},
}
