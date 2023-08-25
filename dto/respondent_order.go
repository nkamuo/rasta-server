package dto

import (
	"github.com/nkamuo/rasta-server/model"
	"github.com/nkamuo/rasta-server/utils/geo"
)

type RespondentOrderEntryIOutput struct {
	Order   model.Order               `json:"order"`
	Routing geo.DistanceMatrixElement `json:"routing"`
}
