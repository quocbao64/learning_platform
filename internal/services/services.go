package services

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAuthService, wire.Bind(new(AuthService), new(*authService)),
	NewUserService, wire.Bind(new(UserService), new(*userService)),
	NewCourseService, wire.Bind(new(CourseService), new(*courseService)),
)
