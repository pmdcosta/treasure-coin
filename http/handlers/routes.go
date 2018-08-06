package handlers

// default pages.
const (
	IndexPage   = "index.html"
	ProfilePage = "profile.html"
	AboutPage   = "about.html"
)

// default routes.
const (
	IndexRoute   = "/"
	ProfileRoute = "/me"
	AboutRoute   = "/about"
)

// auth pages.
const (
	SignInPage = "signin.html"
	SignUpPage = "signup.html"
)

// auth routes.
const (
	SignInRoute  = "/signin"
	SignUpRoute  = "/signup"
	SignOutRoute = "/signout"
)

// game pages.
const (
	CreateGamePage       = "create_game.html"
	DescribeGamePage     = "describe_game.html"
	ListGamePage         = "list_game.html"
	DescribeTreasurePage = "describe_treasure.html"
)

// game routes.
const (
	CreateGameRoute       = "/create"
	DescribeGameRoute     = "/describe/:game"
	ListGameRoute         = "/list"
	DescribeTreasureRoute = "/describe/:game/treasure/:treasure"
	FoundTreasureRoute    = "/found/:game/:treasure"
)
