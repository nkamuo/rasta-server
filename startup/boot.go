package startup

import (
	"github.com/nkamuo/rasta-server/initializers"
	"github.com/nkamuo/rasta-server/service"
)

func Boot() (err error) {
	_, err = initializers.LoadConfig()
	if err != nil {
		return err
	}
	respondentAccessProductBalanceService := service.GetRespondentAccessProductBalanceService()
	respondentAccessProductSubscriptionService := service.GetRespondentAccessProductSubscriptionService()
	if err = respondentAccessProductBalanceService.SetupForAllRespondents(); err != nil {
		return err
	}
	if err = respondentAccessProductSubscriptionService.SetupForAllRespondents(); err != nil {
		return err
	}
	return err
}
