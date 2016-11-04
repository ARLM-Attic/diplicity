package diptest

import "github.com/zond/diplicity/game"

func testOrders(gameDesc string, envs []*Env) {
	g := envs[0].GetRoute(game.IndexRoute).Success().
		Follow("my-started-games", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Desc"}, gameDesc)

	okParts := []string{"", "Hold"}
	badParts := []string{"", "Hold"}

	nation := g.
		Find([]string{"Properties", "Members"}, []string{"User", "Id"}, envs[0].GetUID()).GetValue("Nation")

	switch nation {
	case "Austria":
		okParts[0] = "vie"
		badParts[0] = "ber"
	case "Germany":
		okParts[0] = "ber"
		badParts[0] = "ank"
	case "Turkey":
		okParts[0] = "ank"
		badParts[0] = "rom"
	case "Italy":
		okParts[0] = "rom"
		badParts[0] = "bre"
	case "France":
		okParts[0] = "bre"
		badParts[0] = "mos"
	case "Russia":
		okParts[0] = "mos"
		badParts[0] = "lon"
	case "England":
		okParts[0] = "lon"
		badParts[0] = "vie"
	}

	phase := g.
		Follow("phases", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Season"}, "Spring")

	phase.Follow("orders", "Links").Success().
		AssertEmpty("Properties")

	otherPlayerPhase := envs[1].GetRoute(game.IndexRoute).Success().
		Follow("my-started-games", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Desc"}, gameDesc).
		Follow("phases", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Season"}, "Spring")

	otherPlayerPhase.Follow("orders", "Links").Success().
		AssertEmpty("Properties")

	phase.Follow("create-order", "Links").Body(map[string]interface{}{
		"Parts":    okParts,
		"Province": okParts[0],
	}).Success()

	phase.Follow("orders", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Nation"}, nation)

	otherPlayerPhase.Follow("orders", "Links").Success().
		AssertEmpty("Properties")

	phase.Follow("orders", "Links").Success().
		Find([]string{"Properties"}, []string{"Properties", "Nation"}, nation).
		Follow("delete", "Links").Success()

	phase.Follow("orders", "Links").Success().
		AssertEmpty("Properties")

	phase.Follow("create-order", "Links").Body(map[string]interface{}{
		"Parts":    badParts,
		"Province": badParts[0],
	}).Failure()

	phase.Follow("orders", "Links").Success().
		AssertEmpty("Properties")

}
