package domain

type Repos struct {
	Accounts AccountRepo
}

type Caches struct{}

type Logger struct{}

type Handlers struct {
	Accounts AccountHandlers
}

type Adapters struct {
	Repos  Repos
	Logger Logger
	Caches Caches
}
