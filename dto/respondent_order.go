package dto

import (
	"github.com/nkamuo/rasta-server/utils/geo"
)

type RespondentOrderEntryIOutput struct {
	Order   OrderOutput               `json:"order"`
	Routing geo.DistanceMatrixElement `json:"routing"`
}
