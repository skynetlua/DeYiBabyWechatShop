package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"fmt"
	"github.com/kataras/iris/v12"
)


//=>/gm/team/list true get {} 
func Gm_team_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	dbMyTeams := mysqld.GetAllDBMyTeams()
	var teamList []map[string]interface{}
	playerId2Team := map[int]*map[string]interface{}{}
	var playerIds [] int
	for _,myTeam := range *dbMyTeams {
		team := map[string]interface{}{
			"teamId" 	:myTeam.PlayerId,
			"playerId"	:myTeam.PlayerId,
			"status"	:myTeam.Status,
			"nickname"	:myTeam.Nickname,
			"mobile"	:myTeam.Mobile,

			"timeStamp" :myTeam.UpdatedAt.Unix(),
			"dateAdd" 	:myTeam.GetUpdateDateStr(),
		}
		teamList = append(teamList, team)
		playerIds = append(playerIds, myTeam.PlayerId)
		playerId2Team[myTeam.PlayerId] = &team
	}
	if len(playerIds) > 0 {
		dbPlayerInfos := mysqld.GetDBPlayerInfosByPlayerIdsFromDB(&playerIds)
		for _,dbPlayerInfo := range *dbPlayerInfos {
			team := playerId2Team[dbPlayerInfo.PlayerId]
			if team != nil {
				(*team)["avatarUrl"] = dbPlayerInfo.AvatarUrl
			}else{
				fmt.Println("竟然是空 playerId =", dbPlayerInfo.PlayerId)
			}
		}
	}
	data := map[string]interface{}{
		"teamList": teamList,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/do/team true get {teamId, status} 
func Gm_do_team(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	teamId := common.AtoI(ctx.FormValue("teamId"))
	status := common.AtoI(ctx.FormValue("status"))
	if teamId == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数错误"})
		return
	}

	myTeam := mysqld.GetDBMyTeamFromDB(teamId, false)
	if myTeam == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数错误"})
		return
	}
	switch mysqld.TeamStatus(status) {
		case mysqld.TeamStatusNo:
		case mysqld.TeamStatusCheck:
		case mysqld.TeamStatusRefuse:
		case mysqld.TeamStatusSeller:
		case mysqld.TeamStatusCancel:
		default:
			ctx.JSON(iris.Map{"code": -1, "msg":"参数错误"})
			return
	}
	myTeam.BeginWrite()
	myTeam.Status = status
	myTeam.EndWrite()
	myTeam.Save()

	// dbPlayer := mysqld.GetDBPlayerByPlayerId(myTeam.PlayerId)
	// if dbPlayer != nil {
	// 	dbPlayer.ClearMyTeam()
	// }
	ctx.JSON(iris.Map{"code": 0})
}

