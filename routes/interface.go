package routes

type ApiRouteInterface interface {
	Api() error
}

type WebRouteInterface interface {
	Web() error
}
