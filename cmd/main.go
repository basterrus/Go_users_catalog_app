package main

import (
	"fmt"
	"github.com/basterrus/Go_users_catalog_app/internal/app/catalogBL"
	"github.com/basterrus/Go_users_catalog_app/internal/app/repository/groupUserBL"
	"github.com/basterrus/Go_users_catalog_app/internal/app/repository/usersBL"
	"github.com/basterrus/Go_users_catalog_app/internal/db/memDB/memGroupDB"
	"github.com/basterrus/Go_users_catalog_app/internal/db/memDB/memUserDB"
)

func main() {
	memUserDB := memUserDB.NewMemUserDB()
	memGroupDB := memGroupDB.NewMemGroup()

	userBL := usersBL.NewUsersStore(memUserDB)
	groupBL := groupUserBL.NewGroupUsersStore(memGroupDB)

	catalog := catalogBL.NewCatalog(userBL, groupBL)

	fmt.Println(catalog)

}
