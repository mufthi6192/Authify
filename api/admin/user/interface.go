package user

import responseFormatter "SMM-PPOB/helper/formatter"

type RepositoryInterface interface {
	AddUserRepository(data AddUserDto) responseFormatter.QueryData
	GetUserRepository() responseFormatter.QueryData
	GetUserDetailRepository() responseFormatter.QueryData
	UpdateUserRepository() responseFormatter.QueryData
	DeleteUserRepository(userId string) responseFormatter.QueryData
}

type ServiceInterface interface {
	AddUserService() responseFormatter.HttpData
	GetUserService() responseFormatter.HttpData
	GetUserDetailService() responseFormatter.HttpData
	UpdateUserService() responseFormatter.HttpData
	DeleteUserService() responseFormatter.HttpData
}
