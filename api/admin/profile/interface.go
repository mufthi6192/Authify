package profile

import responseFormatter "SMM-PPOB/helper/formatter"

type RepositoryInterface interface {
	GetProfileRepository(userId string) responseFormatter.QueryData
	GetLoginHistoryRepository(userId string) responseFormatter.QueryData
	GetLatestLoginHistoryRepository(userId string) responseFormatter.QueryData
	ChangePasswordRepository(data ChangePasswordDto) responseFormatter.QueryData
}

type ServiceInterface interface {
	GetProfileService() responseFormatter.HttpData
	GetLoginHistoryService() responseFormatter.HttpData
	GetLatestLoginHistoryService() responseFormatter.HttpData
	ChangePasswordService() responseFormatter.HttpData
}
