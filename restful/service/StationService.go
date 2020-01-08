package service

const stationNameId = ""

func NewStationService() *StationService {
	service := &StationService{
		base: &BaseService{},
	}
	service.base.FindUrl(stationNameId)
	return service
}

type StationService struct {
	base *BaseService
}

func (s *StationService) names(){

}


